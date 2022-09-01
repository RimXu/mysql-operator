# MySQL Operator for Kubernetes

## Introduction

The MySQL Operator for Kubernetes is an operator for managing MySQL Cluster 

## Pull Requests

https://github.com/RimXu/mysql-operator.git

## MySQL Operator for Kubernetes Installation

### Using Manifest Files with kubectl

1.安装kubebuilder：https://kubebuilder.io/ 

2.从github获取源码
```sh
# git clone https://github.com/RimXu/mysql-operator.git
# cd mysql-operator/
```
3.初始化operator：
```
# mkdir -p /code/go/operator/mysql-operator
# cd /code/go/operator/mysql-operator
# go mod init mysql-operator
# kubebuilder init --domain operator.io
# kubebuilder create api --group database --version v1 --kind Mysql
```

4.创建API:
```sh
# kubebuilder create api --group database --version v1 --kind Mysql
```

5.编译和安装,使用make编译，make install部署CRD:
```sh
# make && make install
```

6.启动controller
```sh
# make run
```

7.制作MySQL-operator镜像:
```sh
# docker build . -t mysql-operator:0.1
```

### Use kubectl deploy to kubernetes

示例



