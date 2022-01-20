# <img src='docs/icon.png' width='45' align='center' alt='icon'> AWPark

![](https://img.shields.io/badge/license-MIT-green)  ![](https://img.shields.io/badge/platform-MacOS-purple) ![](https://img.shields.io/badge/language-Go-blue)

Alfred Workflow for engineer.

## Alfred Workflow Store
Search and install [Workflows.json](/static/workflows.json)

Default keyword: `wf`

![workflows store](docs/alfred-workflow-store.png)

* Input [Enter] : Download & Install; (Open workflow's website if url is empty)
* Input [Cmd + Enter] : Open workflow's website.
* Input [Cmd + C] : Copy target URL to the clipboard.

## Develop Kit
Default keyword: `kit`

![develop kit](docs/alfred-workflow-kit.png)

Features:

* Base64 Encode/Decode.
* URL Encode/Decode.
* MD5 Lower/Upper.
* Show IP Address.
  * Local IP.
  * External IP.
* Timestamp Getter/Formatter.
  * Timestamp Getter: Seconds / Millisecond.
  * Timestamp Formatter: Seconds `s:1642700033` to `2022-01-21T01:33:53+08:00`。
* UUID Generator.
* JWT Decode.
