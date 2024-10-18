# gsv

查找 github owner 并展示 stars repository

![demo](./screen.gif)

# 特性
* 支持搜索 github 账号
* 展示 stars 列表
* 支持查看某个库详情（description、stars、forks）
* 支持多个账号切换
* 支持使用浏览器打开指定库

# Install

```shell
go install github.com/xiaoxuan6/gsv@latest
```

# Shell

## install

```shell
bash <(curl -sSL https://github-mirror.us.kg/https:/github.com/xiaoxuan6/gsv/blob/main/gsv.sh) install
```

## remove

```shell
bash <(curl -sSL https://github-mirror.us.kg/https:/github.com/xiaoxuan6/gsv/blob/main/gsv.sh) remove
```

# Releaser

[releases](https://github.com/xiaoxuan6/gsv/releases)

> [!WARNING]  
> [termui](github.com/gizak/termui) Not support win `gitbash`
> 
> `cmder` 无法展示内容
