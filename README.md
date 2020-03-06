# Grandet
Grandet 程序化A股交易工具

# 环境
已测试支持环境：  
- OS
    - Ubuntu 16.04
    - Ubuntu18.04
    - MacOS
- Python
    - 3.7.6  
- Golang
    - 1.11.5  
- Docker
    - 18.09.2
- Docker-compose
    - 1.23.2
- DB
    - postgres 11.1

# 安装
1. 安装go package, 拷贝配置文件
```
go get -v ./...
mv conf.yaml.example conf.yaml
mv Makefile.example Makefile
```

2. 修改 Makefile中 ${YOUR_TUSHARE_TOKEN} 为实际token值([Tushare Token注册申请](https://tushare.pro/))

# 使用
## 启动数据库
```
docker-compose up -d
```
## 获取股票列表以及历史数据
```
make start
```

## 拉取股票列表
支持保存到DB和Excel两种方式，按需调整conf.yaml中 `storage_db` 和 `storage_excel`的值
```
make stock_list
```

## 拉取股票日线数据
支持保存到DB和Excel两种方式，按需调整conf.yaml中 `storage_db` 和 `storage_excel`的值
```
make daily
```

TODO ...

