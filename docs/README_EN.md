<div align="center">
<h1>Go Ldap Admin</h1>

[简体中文](../README.md) | **English**

[![Auth](https://img.shields.io/badge/Auth-eryajf-ff69b4)](https://github.com/eryajf)
[![Go Version](https://img.shields.io/github/go-mod/go-version/eryajf-world/go-ldap-admin)](https://github.com/eryajf/go-ldap-admin)
[![Gin Version](https://img.shields.io/badge/Gin-1.6.3-brightgreen)](https://github.com/eryajf/go-ldap-admin)
[![Gorm Version](https://img.shields.io/badge/Gorm-1.24.5-brightgreen)](https://github.com/eryajf/go-ldap-admin)
[![GitHub Pull Requests](https://img.shields.io/github/stars/eryajf/go-ldap-admin)](https://github.com/eryajf/go-ldap-admin/stargazers)
[![HitCount](https://views.whatilearened.today/views/github/eryajf/go-ldap-admin.svg)](https://github.com/eryajf/go-ldap-admin)
[![GitHub license](https://img.shields.io/github/license/eryajf/go-ldap-admin)](https://github.com/eryajf/go-ldap-admin/blob/main/LICENSE)
[![Commits](https://img.shields.io/github/commit-activity/m/eryajf/go-ldap-admin?color=ffff00)](https://github.com/eryajf/go-ldap-admin/commits/main)

<p>🌉 An OpenLDAP administration platform built with Go + Vue 🌉</p>

<img src="./docs/img/00.webp" width="800"  height="3">
</div><br>

<p align="center">
  <a href="" rel="noopener">
 <img src="./docs/img/01.webp" alt="Project logo"></a>
</p>

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**目录**

- [ℹ️ Introduction](#-introduction)
- [❤️ Sponsors](#-sponsors)
- [🏊 Online Demo](#-online-demo)
- [👨‍💻 Project Repositories](#-project-repositories)
- [🔗 Documentation Links](#-documentation-links)
- [🥰 Acknowledgements](#-acknowledgements)
- [🤗 One More Thing](#-one-more-thing)
- [🤑 Donation](#-donation)
- [📝 User Registration](#-user-registration)
- [💎 Recommended Software](#-recommended-software)
- [🤝 Contributors](#-contributors)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## ℹ️ Introduction

`go-ldap-admin` aims to provide a simple, easy-to-use, clean, and modern administration console for `OpenLDAP`.

> In addition to OpenLDAP administration, the project supports integration with `DingTalk`, `WeCom`, and `Lark`. Users can manually or automatically synchronize organization structures and employee information into the platform, making `go-ldap-admin` a bridge between enterprise IM systems and internal enterprise applications.

## ❤️ Sponsors

[![](./docs/img/12.webp)](https://aiserve.top/)

- An API relay platform for developers and teams that need API access.
- An AI mirror platform for individuals, operations teams, and content teams that want to use AI directly from the browser.

[**Yuanren AI**](https://aiserve.top/): an all-in-one Gemini, ChatGPT, and Claude support platform that provides advanced AI mirror services, helping users in China access models such as ChatGPT and Claude more smoothly.

## 🏊 Online Demo

The online demo is available at:

- URL: [https://demo-go-ldap-admin.eryajf.net/](https://demo-go-ldap-admin.eryajf.net/)
- Login: admin/123456

> The demo environment may be unstable. If access fails or data becomes inconsistent, please contact the author for recovery. Do not enter personal or sensitive information.

**Feature overview:**

|    ![Login page](./docs/img/02.webp)    | ![Home](./docs/img/03.webp)      |
| :----------------------------------------------------------------------------------: | --------------------------------------------------------------------------------- |
|   ![User management](./docs/img/04.webp)   | ![Group management](./docs/img/05.webp)  |
| ![Field relation management](./docs/img/06.webp) | ![Menu management](./docs/img/07.webp)  |
|   ![API management](./docs/img/08.webp)   | ![Operation logs](./docs/img/09.webp)  |
|  ![Swagger](./docs/img/10.webp)   | ![Swagger](./docs/img/11.webp) |

## 👨‍💻 Project Repositories

| Category |                     GitHub                     |   CNB |                     Gitee                        |
| :--: | :--------------------------------------------: | :-------------------------------------------------:| :-------------------------------------------------: |
| Backend |  [go-ldap-admin](https://github.com/opsre/go-ldap-admin.git)   | [go-ldap-admin](https://cnb.cool/opsre/go-ldap-admin.git) | [go-ldap-admin](https://gitee.com/eryajf-world/go-ldap-admin.git)   |
| Frontend | [go-ldap-admin-ui](https://github.com/opsre/go-ldap-admin-ui.git) | [go-ldap-admin-ui](https://cnb.cool/opsre/go-ldap-admin-ui.git) | [go-ldap-admin-ui](https://gitee.com/eryajf-world/go-ldap-admin-ui.git) |

## 🔗 Documentation Links

Project introductions, usage guides, and best practices are all available in the official documentation. If you have questions, please read the documentation first. Common links are listed below.

- [Official website](http://ldapdoc.eryajf.net)
- [Project background](http://ldapdoc.eryajf.net/pages/101948/)
- [Quick start](http://ldapdoc.eryajf.net/pages/706e78/)
- [Feature overview](http://ldapdoc.eryajf.net/pages/7a40de/)
- [Local development](http://ldapdoc.eryajf.net/pages/cb7497/)

> **Notes:**
>
> - Deploying and using this project requires some knowledge of OpenLDAP. If you want to configure IM synchronization, you may also need basic Go knowledge for debugging when issues occur.
> - The documentation is already detailed. Free support is not provided for topics already covered in the documentation. If you encounter problems during installation or deployment, you can request support through the [paid service](http://ldapdoc.eryajf.net/pages/7eab1c/).

## 🥰 Acknowledgements

Thanks to the following excellent projects. Without them, `go-ldap-admin` would not exist:

- Backend stack
  - [Gin-v1.6.3](https://github.com/gin-gonic/gin)
  - [Gorm-v1.24.5](https://github.com/go-gorm/gorm)
  - [Sqlite-v1.7.0](https://github.com/glebarez/sqlite)
  - [Go-ldap-v3.4.2](https://github.com/go-ldap/ldap)
  - [Casbin-v2.22.0](https://github.com/casbin/casbin)
- Frontend stack
  - [axios](https://github.com/axios/axios)
  - [element-ui](https://github.com/ElemeFE/element)
- Special thanks
  - [go-web-mini](https://github.com/gnimli/go-web-mini): this project was refactored based on `go-web-mini`. Thanks to the author for the original work.
  - Thanks to [nangongchengfeng](https://github.com/nangongchengfeng) for contributing the [Swagger](https://github.com/eryajf/go-ldap-admin/pull/345) feature.

## 🤗 One More Thing

- If you find this project useful, please give it a star.
- If you have ideas or feature requests, feel free to discuss them in issues.

## 🤑 Donation

If this project helps you, you can buy the author a coffee: [click here](http://ldapdoc.eryajf.net/pages/2b6725/)

## 📝 User Registration

If your company uses this project, please leave a footprint here. Thank you for your support: [click here](https://github.com/eryajf/go-ldap-admin/issues/18)

## 💎 Recommended Software

- [🦄 TenSunS: an efficient and easy-to-use Consul Web operations platform](https://github.com/starsliao/TenSunS)
- [Next Terminal: a simple, easy-to-use, and secure open-source interactive audit bastion host system](https://github.com/dushixiang/next-terminal)

## 🤝 Contributors

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
