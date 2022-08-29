kubectl delete -f database_v1_mysql.yaml
kubectl delete pvc -n mysql mysql-sample-master-data
kubectl delete pvc -n mysql mysql-sample-slave-data
kubectl delete configmaps -n mysql mysql-sample-slave
kubectl delete configmaps -n mysql mysql-sample-master
