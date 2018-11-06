## 学习Hyperledger Fabric 实战联盟链

### 课程材料
所有的材料都放在imocc目录

#### application
实战项目

#### chaincode
课程中编写的链码

* assetsExchange: 资产交易平台链码
* badexample: 一个错误的示范例子

#### deploy
fabric部署脚本目录

#### docs
课外教辅材料

#### images
一些架构图

### 项目结构

#### bccsp
密码学：加密、签名以及证书等等

#### bddtests
行为驱动开发

需求 概要设计 详细设计 开发
需求 开发

#### common
公共库

* 错误处理
* 日志处理
* 账本存储
* 各种工具
* ...

#### core
核心库，组件核心逻辑

#### devenv
开发环境、Vagrant

#### docs
文档相关

#### events
事件监听机制

#### examples
一些例子

#### gossip
最终一致性共识算法，用于组织内部区块同步

#### images

docker镜像打包

#### msp

成员服务管理 member service provider

#### orderer
排序节点入口

#### peer

peer节点入口

#### proposals
新功能提案

#### protos

grpc: protobuffer + rpc
jsonrpc: json + rpc

**Note:** This is a **read-only mirror** of the formal [Gerrit](https://gerrit.hyperledger.org/r/#/admin/projects/fabric) repository,
where active development is ongoing. Issue tracking is handled in [Jira](https://jira.hyperledger.org/secure/RapidBoard.jspa?projectKey=FAB&rapidView=5&view=planning)

## Status

This project is an _Active_ Hyperledger project. For more information on the history of this project see the [Fabric wiki page](https://wiki.hyperledger.org/projects/Fabric). Information on what _Active_ entails can be found in
the [Hyperledger Project Lifecycle document](https://wiki.hyperledger.org/community/project-lifecycle).

[![Build Status](https://jenkins.hyperledger.org/buildStatus/icon?job=fabric-merge-x86_64)](https://jenkins.hyperledger.org/view/fabric/job/fabric-merge-x86_64/)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/955/badge)](https://bestpractices.coreinfrastructure.org/projects/955)
[![Go Report Card](https://goreportcard.com/badge/github.com/hyperledger/fabric)](https://goreportcard.com/report/github.com/hyperledger/fabric)
[![GoDoc](https://godoc.org/github.com/hyperledger/fabric?status.svg)](https://godoc.org/github.com/hyperledger/fabric)
[![Documentation Status](https://readthedocs.org/projects/hyperledger-fabric/badge/?version=release)](http://hyperledger-fabric.readthedocs.io/en/release/?badge=latest)

## Hyperledger Fabric

Hyperledger Fabric is a platform for distributed ledger solutions, underpinned
by a modular architecture delivering high degrees of confidentiality,
resiliency, flexibility and scalability. It is designed to support pluggable
implementations of different components, and accommodate the complexity and
intricacies that exist across the economic ecosystem.

Hyperledger Fabric delivers a uniquely elastic and extensible architecture,
distinguishing it from alternative blockchain solutions. Planning for the
future of enterprise blockchain requires building on top of a fully-vetted,
open source architecture; Hyperledger Fabric is your starting point.

## Documentation, Getting Started and Developer Guides

Please visit our
[online documentation](http://hyperledger-fabric.readthedocs.io/en/release/) for
information on getting started using and developing with the fabric, SDK and chaincode.

It's recommended for first-time users to begin by going through the
[Getting Started](http://hyperledger-fabric.readthedocs.io/en/release/getting_started.html)
section of the documentation in order to gain familiarity with the Hyperledger
Fabric components and the basic transaction flow.

## Contributing

We welcome contributions to the Hyperledger Fabric Project in many forms.
There’s always plenty to do! Check [the documentation on how to contribute to this project](http://hyperledger-fabric.readthedocs.io/en/latest/CONTRIBUTING.html)
for the full details.

## Community

[Hyperledger Community](https://www.hyperledger.org/community)

[Hyperledger mailing lists and archives](http://lists.hyperledger.org/)

[Hyperledger Chat](http://chat.hyperledger.org/channel/fabric)

[Hyperledger Fabric Issue Tracking](https://jira.hyperledger.org/secure/Dashboard.jspa?selectPageId=10104)

[Hyperledger Wiki](https://wiki.hyperledger.org/)

[Hyperledger Code of Conduct](https://wiki.hyperledger.org/community/hyperledger-project-code-of-conduct)

[Community Calendar](https://wiki.hyperledger.org/community/calendar-public-meetings)

## License <a name="license"></a>

Hyperledger Project source code files are made available under the Apache License, Version 2.0 (Apache-2.0), located in the [LICENSE](LICENSE) file. Hyperledger Project documentation files are made available under the Creative Commons Attribution 4.0 International License (CC-BY-4.0), available at http://creativecommons.org/licenses/by/4.0/.
