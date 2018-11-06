/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tendermint

import (
	"sync"
	"time"

	"github.com/hyperledger/fabric/orderer/multichain"
	cb "github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/op/go-logging"
	"github.com/tendermint/abci/client"
	"github.com/tendermint/abci/types"
)

var logger = logging.MustGetLogger("orderer/tendermint")

type consenter struct{}

type chain struct {
	support  multichain.ConsenterSupport
	sendChan chan *cb.Envelope
	exitChan chan struct{}

	tendermintCli abcicli.Client // 同步通信
	types.BaseApplication
	lastBlock    *types.ResponseInfo
	blocks       map[int64]*block
	deliveredTxs []*cb.Envelope
}

type block struct {
	deliveredTxs []*cb.Envelope
	currentBlock types.RequestBeginBlock
}

var _ types.Application = (*chain)(nil)

// New creates a new consenter for the solo consensus scheme.
// The solo consensus scheme is very simple, and allows only one consenter for a given chain (this process).
// It accepts messages being delivered via Enqueue, orders them, and then uses the blockcutter to form the messages
// into blocks before writing to the given ledger
func New() multichain.Consenter {
	return &consenter{}
}

func (tm *consenter) HandleChain(support multichain.ConsenterSupport, metadata *cb.Metadata) (multichain.Chain, error) {
	return newChain(support), nil
}

func newChain(support multichain.ConsenterSupport) *chain {
	ch := &chain{
		support:      support,
		sendChan:     make(chan *cb.Envelope),
		exitChan:     make(chan struct{}),
		blocks:       make(map[int64]*block),
		deliveredTxs: make([]*cb.Envelope, 0),
	}
	ch.tendermintCli = abcicli.NewLocalClient(&sync.Mutex{}, ch) // 测试
	resp, err := ch.tendermintCli.InfoSync(types.RequestInfo{})
	if err != nil {
		logger.Panic(err)
	}
	ch.lastBlock = resp

	return ch
}

func (ch *chain) Start() {
	go ch.main()
}

func (ch *chain) Halt() {
	select {
	case <-ch.exitChan:
		// Allow multiple halts without panic
	default:
		close(ch.exitChan)
	}
}

// Enqueue accepts a message and returns true on acceptance, or false on shutdown
func (ch *chain) Enqueue(env *cb.Envelope) bool {
	select {
	case ch.sendChan <- env:
		return true
	case <-ch.exitChan:
		return false
	}
}

// Errored only closes on exit
func (ch *chain) Errored() <-chan struct{} {
	return ch.exitChan
}

func (ch *chain) main() {
	ticker := time.NewTicker(ch.support.SharedConfig().BatchTimeout())

	// 新开始一个区块
	_, err := ch.tendermintCli.BeginBlockSync(types.RequestBeginBlock{
		Header: types.Header{
			ChainID: ch.support.ChainID(),
			Height:  int64(ch.support.Height() + 1),
		},
	})
	if err != nil {
		logger.Panicf("Begin block error: %s", err)
	}
	for {
		select {
		case env := <-ch.sendChan:
			envBytes, err := utils.Marshal(env)
			if err != nil {
				logger.Errorf("proto.Marshal return error: %s", err)
				continue
			}

			// 1. 校验交易
			if resp, err := ch.tendermintCli.CheckTxSync(envBytes); err != nil {
				logger.Errorf("env check error: %s", err)
				continue
			} else {
				if resp.Code != uint32(cb.Status_SUCCESS) {
					continue
				}
			}

			// 2. 发送交易
			if resp, err := ch.tendermintCli.DeliverTxSync(envBytes); err != nil {
				logger.Errorf("deliver tx to tendermint error: %s", err)
				continue
			} else {
				if resp.Code != uint32(cb.Status_SUCCESS) {
					continue
				}
			}
		case <-ticker.C:

			// 结束区块
			if _, err := ch.tendermintCli.EndBlockSync(types.RequestEndBlock{Height: int64(ch.support.Height() + 1)}); err != nil {
				logger.Panic(err)
			}

			// 提交区块
			if _, err := ch.tendermintCli.CommitSync(); err != nil {
				logger.Panic(err)
			}

			// 新开始一个区块
			_, err := ch.tendermintCli.BeginBlockSync(types.RequestBeginBlock{
				Header: types.Header{
					ChainID: ch.support.ChainID(),
					Height:  int64(ch.support.Height() + 1),
				},
			})
			if err != nil {
				logger.Panicf("Begin block error: %s", err)
			}
		case <-ch.exitChan:
			logger.Debugf("Exiting")
			return
		}
	}
}

func (ch *chain) Info(req types.RequestInfo) types.ResponseInfo {
	return types.ResponseInfo{
		Data:            "",
		Version:         req.Version,
		LastBlockHeight: int64(ch.support.Height()),
	}
}

func (c *chain) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	c.blocks[req.GetHeader().GetHeight()] = &block{
		deliveredTxs: make([]*cb.Envelope, 0),
		currentBlock: req,
	}
	return types.ResponseBeginBlock{}
}

func (c *chain) CheckTx(tx []byte) types.ResponseCheckTx {
	msg, err := utils.UnmarshalEnvelope(tx)
	if err != nil {
		return types.ResponseCheckTx{Code: uint32(cb.Status_BAD_REQUEST)}
	}

	// 消息体的简单校验
	payload, err := utils.UnmarshalPayload(msg.Payload)
	if err != nil {
		logger.Warningf("Received malformed message, dropping connection: %s", err)
		return types.ResponseCheckTx{Code: uint32(cb.Status_BAD_REQUEST)}
	}

	if payload.Header == nil {
		logger.Warningf("Received malformed message, with missing header, dropping connection")
		return types.ResponseCheckTx{Code: uint32(cb.Status_BAD_REQUEST)}
	}

	_, err = utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		logger.Warningf("Received malformed message (bad channel header), dropping connection: %s", err)
		return types.ResponseCheckTx{Code: uint32(cb.Status_BAD_REQUEST)}
	}

	return types.ResponseCheckTx{Code: uint32(cb.Status_SUCCESS)}
}

// Deliver a tx for full processing
func (ch *chain) DeliverTx(tx []byte) types.ResponseDeliverTx {
	msg, err := utils.UnmarshalEnvelope(tx)
	if err != nil {
		logger.Errorf("Unmarshal tx envelope error: %s", err)
		return types.ResponseDeliverTx{Code: uint32(cb.Status_BAD_REQUEST)}
	}

	select {
	case <-ch.exitChan:
		return types.ResponseDeliverTx{Code: uint32(cb.Status_INTERNAL_SERVER_ERROR)}
	default:
		ch.deliveredTxs = append(ch.deliveredTxs, msg)
		return types.ResponseDeliverTx{Code: uint32(cb.Status_SUCCESS)}
	}
}

func (ch *chain) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	ch.blocks[req.GetHeight()].deliveredTxs = append(ch.blocks[req.GetHeight()].deliveredTxs, ch.deliveredTxs...)
	ch.deliveredTxs = ch.deliveredTxs[:0]
	return types.ResponseEndBlock{}
}

func (ch *chain) Commit() types.ResponseCommit {
	for _, blk := range ch.blocks {
		block := ch.support.CreateNextBlock(blk.deliveredTxs)
		ch.support.WriteBlock(block, nil, nil)
		ch.lastBlock = &types.ResponseInfo{
			LastBlockHeight: int64(block.GetHeader().GetNumber()),
		}
	}
	ch.blocks = make(map[int64]*block)
	ch.deliveredTxs = ch.deliveredTxs[:0]

	return types.ResponseCommit{}
}
