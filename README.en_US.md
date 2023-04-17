# <img src='docs/icon.png' width='45' align='center' alt='icon'> AWPark

![](https://img.shields.io/badge/license-MIT-green)  ![](https://img.shields.io/badge/platform-MacOS-purple) ![](https://img.shields.io/badge/language-Go-blue)

Dedicated to creating an Alfred Workflow Store that integrates commonly used tools for developers.

## Alfred Workflow Store
The Alfred Workflow Store allows users to search for target workflows using keywords such as name, author, or tags, and install them with a single click. Currently, there are over 1,000 workflows available for download, all of which are included in the [Workflows.json](/static/workflows.json) file.

Default keyword: `wf`

![workflows store](docs/alfred-workflow-store.png)

* Pressing [Enter]: initiates a background download, and upon completion, an Alfred prompt will appear to import the workflow. Click "Import" to complete the installation.
* Pressing [Cmd + Enter]: opens the workflow's homepage.
* Pressing [Cmd + C]: copies the workflow download link.

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
