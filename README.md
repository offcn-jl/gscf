# CSCF ( Chaos SCF ) 
CSCF ( Chaos SCF ) , Is a framework fit [GIN](https://github.com/gin-gonic/gin) to [SCF](https://cloud.tencent.com/document/product/583) . 
CSCF ( Chaos SCF ) 是一个帮助传统 WEB 服务适配 [SCF](https://cloud.tencent.com/document/product/583) 的修改版 [GIN](https://github.com/gin-gonic/gin) 框架。

[![Status](https://img.shields.io/badge/Status-Beta-yellow)](#当前版本) [![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/github.com/offcn-jl/cscf)](https://goreportcard.com/report/github.com/offcn-jl/cscf) [![Master Build](https://github.com/offcn-jl/cscf/workflows/Master%20Build/badge.svg)](https://github.com/offcn-jl/cscf/actions?query=workflow%3A%22Master+Build%22) [![codecov](https://codecov.io/gh/offcn-jl/cscf/branch/master/graph/badge.svg)](https://codecov.io/gh/offcn-jl/cscf) [![Build](https://github.com/offcn-jl/cscf/workflows/Build/badge.svg)](https://github.com/offcn-jl/cscf/actions?query=workflow%3ABuild) [![codecov](https://codecov.io/gh/offcn-jl/cscf/branch/new-feature/graph/badge.svg)](https://codecov.io/gh/offcn-jl/cscf/branch/new-feature) 

基于 [Gin](https://github.com/gin-gonic/gin) 框架进行重写，舍弃 Router 组件 ( 使用 [API 网关](#APIGateway) 替代 ) 、多线程管理相关逻辑、各种加锁解锁逻辑，并对 [SCF](#SCF) 进行适配。预期实现对最终效果是可以直接将基于 Gin 框架编写的 Handler 迁移到 SCF 运行，实现将 Gin 的 Handler 作为 SCF "微服务" 运行的效果。

## 当前版本
当前版本 : 测试版

暂定版本发布流程 : Alpha -> Beta -> RC -> GA

> Alpha : 内部测试版, 一般不向外部发布  
> Beta : 也是测试版, 这个阶段的版本会一直加入新的功能  
> RC : 发行候选版本, 基本不再加入新的功能, 主要进行缺陷修复  
> GA : 正式发布的版本, 采用 Release X.Y.Z 作为发布版本号  
> 参考 : [Alpha、Beta、RC、GA版本的区别](http://www.blogjava.net/RomulusW/archive/2008/05/04/197985.html) 、 [软件版本GA、RC、beta等含义](https://blog.csdn.net/gnail_oug/article/details/79998154)

###### APIGateway
腾讯云 [API 网关](https://cloud.tencent.com/product/apigateway) ( API Gateway )是 [腾讯云](https://cloud.tencent.com) 推出的一种 API 托管服务，能提供 API 的完整生命周期管理，包括创建、维护、发布、运行、下线等。

###### SCF
[云函数](https://cloud.tencent.com/document/product/583) ( Serverless Cloud Function，SCF ) 是 [腾讯云](https://cloud.tencent.com) 为企业和开发者们提供的无服务器执行环境。

###### Chaos
> 卡俄斯（又译卡欧斯、卡奥斯；英语/拉丁语：Chaos；希腊语：Χάος），原初混沌，世界的开始。  
  根据古希腊诗人赫西俄德《神谱》（公元前8世纪）描述：宇宙之初，卡俄斯最先独自诞生，是一条无边无际、充满黑暗的裂缝。随卡俄斯之后，先诞生大地之神盖亚（Gaia）随后地狱深渊神塔耳塔洛斯（Tartarus）和爱神厄洛斯（Eros）相继独立诞生，世界由此开始。  
  卡俄斯原本并非“混沌”之意，这一延续至今的意思源自罗马诗人奥维德《变形记》（1.7-9），其中卡俄斯被描述为“一团乱糟糟，没有秩序的物体”，chaos至此已经失去了神性，演变成了“混沌”之意。  
> 摘自 : [百度百科](https://baike.baidu.com/item/卡俄斯/10724560?fromtitle=chaos&fromid=85022#viewPageContent)
