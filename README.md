# <img src='docs/icon.png' width='45' align='center' alt='icon'> AWPark
![](https://img.shields.io/badge/license-MIT-green)  ![](https://img.shields.io/badge/platform-MacOS-purple) ![](https://img.shields.io/badge/language-Go-blue)

[English](README.en_US.md)

致力打造一个 Alfred Workflow 商城，以及集成开发者常用工具。

## Alfred Workflow Store
可以通过关键字（名称/作者/标签）检索目标 Workflow，并一键安装。目前有 1k+ 的 Workflows 可以下载，都收录在 [Workflows.json](/static/workflows.json) 文件中。

触发关键字: `wf`

![workflows store](docs/alfred-workflow-store.png)

搜索结果列表操作：
* 输入 [Enter] : 会开始后台下载，下载完成后会弹出 Alfred 导入 Workflow 弹窗，点击导入即可完成安装。
* 输入 [Cmd + Enter] : 打开 Workflow 主页。
* 输入 [Cmd + C] : 拷贝 Workflow 下载地址。

## Develop Kit
触发关键字: `kit`

![develop kit](docs/alfred-workflow-kit.png)

Features:

* Base64 Encode/Decode.
* URL Encode/Decode.
* MD5 Lower/Upper.
* Show IP Address.
  * 本地IP。
  * 公网IP。
* Timestamp Getter/Formatter.
  * Timestamp Getter: 时间戳秒、时间戳毫秒。
  * Timestamp Formatter：`s:1642700033` 秒格式化成 `2022-01-21T01:33:53+08:00`。
* UUID Generator.
* JWT Decode.
