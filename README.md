# Doubak

[![GitHub Action CI](https://github.com/its-my-data/doubak/workflows/CI/badge.svg?branch=master)](https://github.com/its-my-data/doubak/actions?query=workflow%3ACI)
[![GitHub Action CI](https://github.com/its-my-data/doubak/workflows/Release/badge.svg?branch=master)](https://github.com/its-my-data/doubak/actions?query=workflow%3ARelease)

豆瓣 (Douban.com) 备份工具命令行版，备份完整页面镜像，Golang实现。

## 为什么会有此项目？
偶然间发现以前玩过的游戏突然不见了，看过的电影也404了，这种情况真的很崩溃。
互联网本身就带有很强的健忘症属性，加上祖国的规则，自己生成（写）的数据绝对是没有保障性的。
曾经是网易核心产品的163博客，现在也全部404。博文图片也都挂掉了。
仔细想想，这些自己写的数据何尝不是一种日记呢？如果日记都可以被随意删除，我觉得这些平台可以死ね！尤其是豆瓣！

我个人是豆瓣的重度用户，所有电影、书、游戏，开坑前和完坑后都会在豆瓣上标记一下，这是我的记忆，应该由我做主。
所以我决定 self-host 豆瓣。我现有的数据必然都需要爬下来，然后我需要能创建自己的条目（包括被祖国封锁的条目），我做了我就要标记！
当然，自己生成的东西都要有上下文才make sense，所以整个页面都会爬下来，版本也会保留。
不再受制于平台是我的终极目标。

类似的项目？有的。譬如 [豆伴](https://github.com/doufen-org/tofu)。
功能很完善了，但是用起来还是不顺手，导出的东西缺失context，而且不利于自己host网页版的豆瓣，简言之就是纯“自嗨”用的。

本项目会可以增量爬豆瓣，而且可以保留版本历史。
所有数据都会存在数据库里，不出意外的话是SQLite，然后也可以导出成json。

还有最重要的一点，你们看电影会看第二遍吗？每看一遍感悟肯定会不一样吧？
豆瓣怎么就没这个功能呢？明明是刚需啊。。。
所以我这个工具保留版本历史的作用之一就是，你后来又看了一遍，就可以保留以前的感悟和现在的感悟。

One more thing, 我就特别需要标记为“不想看”的功能。
原因有很多，比如我看了某个短片介绍，我觉得这个片子不好，标记一下雷区。这个也是刚需吧？
反正自己来实现好了。

目前本项目爬出来/处理过的数据储存在 [mewx.github.io-Generator](https://github.com/MewX/mewx.github.io-Generator/tree/master/data/doubak) 中。
我个人会在另一个项目里使用 [Gatsby.js](https://www.gatsbyjs.org/) 用这些爬到的数据生成 self-hosted 豆瓣，托管在 GitHub Pages 上。
当然啦， [Archive.org](http://archive.org) 会是我的另一个备份。

欢迎大家交流 self-hosting everything 的经验。

## How to compile?

### Using Bazel (recommended)

1. Install [Bazel](https://docs.bazel.build/versions/master/install-ubuntu.html) or [Bazelisk](https://github.com/bazelbuild/bazelisk/releases).

2. Install protoc and golang-goprotobuf-dev. Ubuntu for example:

    ```shell
    $ sudo apt install protobuf-compiler golang-goprotobuf-dev
    ```

3. Compile.

    ```shell
    $ bazel build :doubak

    # Or if you are using Bazelisk.
    $ bazelisk build :doubak
    ```
4. Run. Note that not using `bazel/bazelisk run` to avoid generating files into temp folders.

    ```shell
    $ bazel-bin/doubak_/doubak \
        --user=<user_name> \
        --categories=<categories> \
        --cookies_file=<cookies.txt>
    ```

### Manually

1. If you want to compile from proto directly,
install protoc and golang-goprotobuf-dev. Ubuntu for example:

    ```shell
    $ sudo apt install protobuf-compiler golang-goprotobuf-dev
    ```

2. Compile:

    ```shell
    # If you didn't install protoc, just skip this step.
    $ ./compile_protos.sh

    # Using legacy go build.
    $ go build
    ```

3. Run.

    ```shell
    $ ./doubak \
        --user=<user_name> \
        --categories=<categories> \
        --cookies_file=<cookies.txt>
    ```
