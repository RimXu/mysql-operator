package controllers

import (
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type MysqlPrepare struct {
}


func (d MysqlPrepare) Create(evt event.CreateEvent) bool {
	fmt.Println("Create")
	return false
}

func (d MysqlPrepare) Delete(evt event.DeleteEvent) bool {
	fmt.Println("Delete")
	return false
}

func (d MysqlPrepare) Update(evt event.UpdateEvent) bool {
	fmt.Println("Update")
	return false
}

func (d MysqlPrepare) Generic(evt event.GenericEvent) bool {
	fmt.Println("Generic")
	return false
}