# MySQL Operator for Kubernetes

## 实现方式

使用mysql主从作为数据库，StorageClass作为后端存储，proxysql作为应用透明代理，job完成初始化配置任务

## 代码获取

https://github.com/RimXu/mysql-operator.git

## MySQL Operator 安装

### 使用kubectl安装

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

### 使用kubectl部署
示例
```yaml
apiVersion: database.operator.io/v1
kind: Mysql
metadata:
  name: mysql-sample
  namespace: mysql
spec:
  # TODO(user): Add fields here
  replication: true
  instance: medium
  storageclass: "nfs-client-storageclass"
```


