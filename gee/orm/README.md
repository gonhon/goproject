<!--
 * @Author: gaoh
 * @Date: 2023-04-14 00:00:06
 * @LastEditTime: 2023-04-14 00:01:11
-->

# Orm

## sqllite安装

https://www.runoob.com/sqlite/sqlite-installation.html

## wsl

```shell
apt-get install sqlite3

cd /mnt/d/project/go_project/goproject


# cygwin 添加 MinGW-64 gcc 安装修改 go env CC
go env -w 
```

```shell
# 进入gee库
 sqlite3 gee.db


 CREATE TABLE User(Name text, Age integer);

 INSERT INTO User(Name, Age) VALUES ('Tom', 18), ('Jack', 25);

 SELECT * FROM User

# 查看所有表
 .table

# 查看表结构
 .schema User
```