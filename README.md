# MySQL Operator for Kubernetes

## 实现方式

使用mysql主从作为数据库，StorageClass作为后端存储，proxysql作为应用透明代理，job完成初始化配置任务，
支持mysql5.7和mysql8.0

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
示例database_v2_mysql.yaml
```yaml
apiVersion: database.operator.io/v1
kind: Mysql
metadata:
  # MySQL operator名称和命名空间
  name: mysql-v8
  namespace: mysql
spec:
  # 定义MySQL是否需要主从，必选
  replication: true
  # 定义MySQL实例规格,可以在constants.InstanceReflect中定义
  # 也可以通过环境变量配置Base_MEM/Base_CPU定义基数,small/medium/large分别是1倍、2倍、4倍基数
  # Base_MEM单位是Mi
  # Base_CPU单位是m
  instance: medium
  # 定义MySQL存储后端Storageclass
  storageclass: "nfs-client-storageclass"
  # 数据库版本定义：与镜像的名称相同
  # 镜像仓库地址可以通过Registry_Addr设置
  version: mysql:8.0.33
  # 数据库列表,支持多数据库多用户
  databases:
    - name: mydb1
      user: mydb1
      passwd: mydbPWD1
    - name: mydb2
      user: mydb2
      passwd: mydbPWD2


```
```sh
kubectl apply -f database_v1_mysql.yaml
```

### 使用kubectl查询状态
```sh
kubectl get mysqls -n mysql
NAME           REPLICATION   INSTANCE   STORAGECLASS              PHASE
mysql-sample   false         small      nfs-client-storageclass   JobRunning

# 当mysql初始化任务完成之后PHASE会变成JobSuccess(或者JobFailed)
NAME           REPLICATION   INSTANCE   STORAGECLASS              PHASE
mysql-sample   false         small      nfs-client-storageclass   JobSuccess

```
