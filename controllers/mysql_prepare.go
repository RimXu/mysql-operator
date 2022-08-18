package controllers

import (
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/event"
	//appsv1 "k8s.io/api/apps/v1"
	//mysqlv1 "mysql-operator/api/v1"

)

type MysqlPrepare struct {
}

func (d MysqlPrepare) Create(evt event.CreateEvent) bool {
	fmt.Println("Create")
	//fmt.Println(evt.Object.(*appsv1.Deployment).ObjectMeta.Namespace)
	//fmt.Println(evt.Object.GetName())
	//evt.Object.(*nginxv1.Nginx).ObjectMeta
	// 只有return true 才能将数据加入到workqueue中
	return true
}

func (d MysqlPrepare) Delete(evt event.DeleteEvent) bool {
	fmt.Println("Delete")
	return true
}

func (d MysqlPrepare) Update(evt event.UpdateEvent) bool {
	fmt.Println("Update")
	return true
}

func (d MysqlPrepare) Generic(evt event.GenericEvent) bool {
	fmt.Println("Generic")
	return false
}
