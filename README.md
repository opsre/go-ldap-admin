<div align="center">
<h1>Go Ldap Admin</h1>

[![Auth](https://img.shields.io/badge/Auth-eryajf-ff69b4)](https://github.com/eryajf)
[![Go Version](https://img.shields.io/github/go-mod/go-version/eryajf-world/go-ldap-admin)](https://github.com/eryajf/go-ldap-admin)
[![Gin Version](https://img.shields.io/badge/Gin-1.6.3-brightgreen)](https://github.com/eryajf/go-ldap-admin)
[![Gorm Version](https://img.shields.io/badge/Gorm-1.24.5-brightgreen)](https://github.com/eryajf/go-ldap-admin)
[![GitHub Pull Requests](https://img.shields.io/github/stars/eryajf/go-ldap-admin)](https://github.com/eryajf/go-ldap-admin/stargazers)
[![HitCount](https://views.whatilearened.today/views/github/eryajf/go-ldap-admin.svg)](https://github.com/eryajf/go-ldap-admin)
[![GitHub license](https://img.shields.io/github/license/eryajf/go-ldap-admin)](https://github.com/eryajf/go-ldap-admin/blob/main/LICENSE)
[![Commits](https://img.shields.io/github/commit-activity/m/eryajf/go-ldap-admin?color=ffff00)](https://github.com/eryajf/go-ldap-admin/commits/main)

<p> 🌉 基于Go+Vue实现的openLDAP后台管理项目 🌉</p>

<img src="https://cdn.jsdelivr.net/gh/eryajf/tu@main/img/image_20240420_214408.gif" width="800"  height="3">
</div><br>

<p align="center">
  <a href="" rel="noopener">
 <img src="https://cdn.jsdelivr.net/gh/eryajf/tu/img/image_20220614_131521.jpg" alt="Project logo"></a>
</p>

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**目录**

- [ℹ️ 项目简介](#-%E9%A1%B9%E7%9B%AE%E7%AE%80%E4%BB%8B)
- [❤️ 赞助商](#-%E8%B5%9E%E5%8A%A9%E5%95%86)
- [🏊 在线体验](#-%E5%9C%A8%E7%BA%BF%E4%BD%93%E9%AA%8C)
- [👨‍💻 项目地址](#-%E9%A1%B9%E7%9B%AE%E5%9C%B0%E5%9D%80)
- [🔗 文档快链](#-%E6%96%87%E6%A1%A3%E5%BF%AB%E9%93%BE)
- [🥰 感谢](#-%E6%84%9F%E8%B0%A2)
- [🤗 另外](#-%E5%8F%A6%E5%A4%96)
- [🤑 捐赠](#-%E6%8D%90%E8%B5%A0)
- [📝 使用登记](#-%E4%BD%BF%E7%94%A8%E7%99%BB%E8%AE%B0)
- [💎 优秀软件推荐](#-%E4%BC%98%E7%A7%80%E8%BD%AF%E4%BB%B6%E6%8E%A8%E8%8D%90)
- [🤝 贡献者](#-%E8%B4%A1%E7%8C%AE%E8%80%85)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## ℹ️ 项目简介

`go-ldap-admin`旨在为`OpenLDAP`服务端提供一个简单易用，清晰美观的现代化管理后台。

> 在完成针对`OpenLDAP`的管理能力之下，支持对`钉钉`，`企业微信`，`飞书`的集成，用户可以选择手动或者自动同步组织架构以及员工信息到平台中，让`go-ldap-admin`项目成为打通企业 IM 与企业内网应用之间的桥梁。

## ❤️ 赞助商

[![](https://wsrv.nl/?url=https://raw.githubusercontent.com/eryajf/tu/refs/heads/main/img/image_20260422_094140.png)](https://aiserve.top/)


- 是 API 中转站: 适合开发者和有接口接入需求的团队。
- 也是 AI 镜像站: 适合个人、运营、内容团队，打开网页直接使用。

[**猿人AI**](https://aiserve.top/): Gemini | ChatGPT | Claude 一站式支持平台，提供先进的AI镜像服务，助你在国内优雅的使用 ChatGPT 和 Claude 等模型。

## 🏊 在线体验

提供在线体验地址如下：

- 地址：[https://demo-go-ldap-admin.eryajf.net/](https://demo-go-ldap-admin.eryajf.net/)
- 登陆信息：admin/123456

> 在线环境可能不稳，如果遇到访问异常，或者数据错乱，请联系我进行修复。请勿填写个人敏感信息。


**页面功能概览：**

|    ![登录页](https://cdn.jsdelivr.net/gh/eryajf/tu/img/image_20220724_165411.png)    | ![首页](https://cdn.jsdelivr.net/gh/eryajf/tu/img/image_20220724_165545.png)      |
| :----------------------------------------------------------------------------------: | --------------------------------------------------------------------------------- |
|   ![用户管理](https://cdn.jsdelivr.net/gh/eryajf/tu/img/image_20220724_165623.png)   | ![分组管理](https://cdn.jsdelivr.net/gh/eryajf/tu/img/image_20220724_165701.png)  |
| ![字段关系管理](https://cdn.jsdelivr.net/gh/eryajf/tu/img/image_20220724_165853.png) | ![菜单管理](https://cdn.jsdelivr.net/gh/eryajf/tu/img/image_20220724_165954.png)  |
|   ![接口管理](https://cdn.jsdelivr.net/gh/eryajf/tu/img/image_20220724_170015.png)   | ![操作日志](https://cdn.jsdelivr.net/gh/eryajf/tu/img/image_20220724_170035.png)  |
|  ![swag](https://cdn.jsdelivr.net/gh/eryajf/tu@main/img/image_20240521_213841.png)   | ![swag](https://cdn.jsdelivr.net/gh/eryajf/tu@main/img/image_20240521_214025.png) |

## 👨‍💻 项目地址

| 分类 |                     GitHub                     |   CNB |                     Gitee                        |
| :--: | :--------------------------------------------: | :-------------------------------------------------:| :-------------------------------------------------: |
| 后端 |  [go-ldap-admin](https://github.com/opsre/go-ldap-admin.git)   | [go-ldap-admin](https://cnb.cool/opsre/go-ldap-admin.git) | [go-ldap-admin](https://gitee.com/eryajf-world/go-ldap-admin.git)   |
| 前端 | [go-ldap-admin-ui](https://github.com/opsre/go-ldap-admin-ui.git) | [go-ldap-admin-ui](https://cnb.cool/opsre/go-ldap-admin-ui.git) | [go-ldap-admin-ui](https://gitee.com/eryajf-world/go-ldap-admin-ui.git) |

## 🔗 文档快链

项目相关介绍，使用，最佳实践等相关内容，都会在官方文档呈现，如有疑问，请先阅读官方文档，以下列举以下常用快链。

- [官网地址](http://ldapdoc.eryajf.net)
- [项目背景](http://ldapdoc.eryajf.net/pages/101948/)
- [快速开始](http://ldapdoc.eryajf.net/pages/706e78/)
- [功能概览](http://ldapdoc.eryajf.net/pages/7a40de/)
- [本地开发](http://ldapdoc.eryajf.net/pages/cb7497/)

> **说明：**
>
> - 本项目的部署与使用需要你对 OpenLDAP 有一定的掌握，如果想要配置 IM 同步，可能还需要一定的 go 基础来调试(如有异常时)。
> - 文档已足够详尽，所有文档已讲过的，将不再提供免费的服务。如果你在安装部署时遇到问题，可通过[付费服务](http://ldapdoc.eryajf.net/pages/7eab1c/)寻求支持。


## 🥰 感谢

感谢如下优秀的项目，没有这些项目，不可能会有 go-ldap-admin：

- 后端技术栈
  - [Gin-v1.6.3](https://github.com/gin-gonic/gin)
  - [Gorm-v1.24.5](https://github.com/go-gorm/gorm)
  - [Sqlite-v1.7.0](https://github.com/glebarez/sqlite)
  - [Go-ldap-v3.4.2](https://github.com/go-ldap/ldap)
  - [Casbin-v2.22.0](https://github.com/casbin/casbin)
- 前端技术栈

  - [axios](https://github.com/axios/axios)
  - [element-ui](https://github.com/ElemeFE/element)

- 另外感谢
  - [go-web-mini](https://github.com/gnimli/go-web-mini)：项目基于该项目重构而成，感谢作者的付出。
  - 感谢 [nangongchengfeng](https://github.com/nangongchengfeng) 提交的 [swagger](https://github.com/eryajf/go-ldap-admin/pull/345) 功能。

## 🤗 另外

- 如果觉得项目不错，麻烦动动小手点个 ⭐️star⭐️!
- 如果你还有其他想法或者需求，欢迎在 issue 中交流！

## 🤑 捐赠

如果你觉得这个项目对你有帮助，你可以请作者喝杯咖啡 ☕️ [点我](http://ldapdoc.eryajf.net/pages/2b6725/)

## 📝 使用登记

如果你所在公司使用了该项目，烦请在这里留下脚印，感谢支持 🥳 [点我](https://github.com/eryajf/go-ldap-admin/issues/18)

## 💎 优秀软件推荐

- [🦄 TenSunS：高效易用的 Consul Web 运维平台](https://github.com/starsliao/TenSunS)
- [ Next Terminal：一个简单好用安全的开源交互审计堡垒机系统](https://github.com/dushixiang/next-terminal)

## 🤝 贡献者

<!-- readme: collaborators,contributors -start -->
<table>
<tr>
    <td align="center">
        <a href="https://github.com/eryajf">
            <img src="https://avatars.githubusercontent.com/u/33259379?v=4" width="100;" alt="eryajf"/>
            <br />
            <sub><b>二丫讲梵</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/xinyuandd">
            <img src="https://avatars.githubusercontent.com/u/3397848?v=4" width="100;" alt="xinyuandd"/>
            <br />
            <sub><b>Xinyuandd</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/RoninZc">
            <img src="https://avatars.githubusercontent.com/u/48718694?v=4" width="100;" alt="RoninZc"/>
            <br />
            <sub><b>Ronin_Zc</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/wang-xiaowu">
            <img src="https://avatars.githubusercontent.com/u/44340137?v=4" width="100;" alt="wang-xiaowu"/>
            <br />
            <sub><b>Xiaowu</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/nangongchengfeng">
            <img src="https://avatars.githubusercontent.com/u/46562911?v=4" width="100;" alt="nangongchengfeng"/>
            <br />
            <sub><b>南宫乘风</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/huxiangquan">
            <img src="https://avatars.githubusercontent.com/u/52623921?v=4" width="100;" alt="huxiangquan"/>
            <br />
            <sub><b>Null</b></sub>
        </a>
    </td></tr>
<tr>
    <td align="center">
        <a href="https://github.com/0x0034">
            <img src="https://avatars.githubusercontent.com/u/39284250?v=4" width="100;" alt="0x0034"/>
            <br />
            <sub><b>0x0034</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/Pepperpotato">
            <img src="https://avatars.githubusercontent.com/u/49708116?v=4" width="100;" alt="Pepperpotato"/>
            <br />
            <sub><b>Null</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/yuliang2">
            <img src="https://avatars.githubusercontent.com/u/63152460?v=4" width="100;" alt="yuliang2"/>
            <br />
            <sub><b>Xu Jiang</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/Foustdg">
            <img src="https://avatars.githubusercontent.com/u/20092889?v=4" width="100;" alt="Foustdg"/>
            <br />
            <sub><b>YD-SUN</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/ckyoung123421">
            <img src="https://avatars.githubusercontent.com/u/16368382?v=4" width="100;" alt="ckyoung123421"/>
            <br />
            <sub><b>Ckyoung123421</b></sub>
        </a>
    </td></tr>
</table>
<!-- readme: collaborators,contributors -end -->
