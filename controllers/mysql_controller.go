/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	mysqlv1 "mysql-operator/api/v1"
	"mysql-operator/pkg/constants"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	apiv1 "k8s.io/api/core/v1"
)

// MysqlReconciler reconciles a Mysql object
type MysqlReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=database.operator.io,resources=mysqls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=database.operator.io,resources=mysqls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=database.operator.io,resources=mysqls/finalizers,verbs=update
func (r *MysqlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)

	logrus.Info("MySQL-Operator reconciler start")
	mysqloperator := &mysqlv1.Mysql{}


	// 查询Namespace下是否存在mysqloperator,如果不存在则满足errors.IsNotFound(err),函数返回
	err := r.Get(context.TODO(), req.NamespacedName, mysqloperator)

	// uuid 判空说明是删除namespace下指定的operator
	uuid := mysqloperator.ObjectMeta.UID
	if uuid == "" {
		logrus.Info("MySQL-Operator reconciler delete")
	}
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}
	mysqldep := &appsv1.Deployment{}
	combo := mysqloperator.Spec.Combo
	sc := mysqloperator.Spec.StorageClass
	size := constants.ComboReflect[combo]["Disk"]

	//判断是否需要创建MySQL主从主从
	if mysqloperator.Spec.Replication == true {
		// 查询{Name:my.Name,Namespace:my.Namespace} 是否存在deployment,如果不存在则满足errors.IsNotFound(err)
		err = r.Get(context.TODO(), types.NamespacedName{Name: mysqloperator.Name + "-master", Namespace: mysqloperator.Namespace}, mysqldep)
		if err != nil {
			logrus.Info("MySQL replication configuration")
			// TEST BEGIN
			//pvc := &apiv1.PersistentVolume{}
			//err := r.Get(context.TODO(), types.NamespacedName{Namespace: mysqloperator.Namespace}, pvc)
			//fmt.Println(pvc)
			// TEST END
			if errors.IsNotFound(err) {
				volumeerr := r.CreateVolumes(mysqloperator.Namespace,sc,mysqloperator.Name + "-master-data",size,ctx)
				if volumeerr != nil {
					fmt.Println(volumeerr)
				}
				master := r.CreateMysql(mysqloperator, "-master", combo)
				if err = r.Create(context.TODO(), master); err != nil {
					return ctrl.Result{}, err
				}
				if err := r.Status().Update(ctx, mysqloperator); err != nil {
					logrus.Error(err, "MySQL master status update error")
				}
				return ctrl.Result{Requeue: true}, nil
			} else {
				return ctrl.Result{}, err
			}
		}
		err = r.Get(context.TODO(), types.NamespacedName{Name: mysqloperator.Name + "-slave", Namespace: mysqloperator.Namespace}, mysqldep)
		if err != nil {
			if errors.IsNotFound(err) {
				slave := r.CreateMysql(mysqloperator, "-slave", combo)
				if err = r.Create(context.TODO(), slave); err != nil {
					return ctrl.Result{}, err
				}
				if err := r.Status().Update(ctx, mysqloperator); err != nil {
					logrus.Error(err, "MySQL slave status update error")
				}
				return ctrl.Result{Requeue: true}, nil
			} else {
				return ctrl.Result{}, err
			}
		}
	} else {
		logrus.Info("MySQL Standalone configuration")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MysqlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// 定义operator事件过滤
	e := MysqlPrepare{}
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1.Mysql{}).
		// 使用WithEventFilter过滤资源事件
		WithEventFilter(e).
		// 如果关注operator创建的deployment事件,可以使用Owns方法,其他资源可以使用Watch方法
		Owns(&appsv1.Deployment{}).
		Complete(r)
}


// FuncVolumes
func (r *MysqlReconciler) CreateVolumes(ns string, sc string, name string, size string,ctx context.Context) error {
	//
	foundConfigMap := &apiv1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: "kube-root-ca.crt", Namespace: "default"}, foundConfigMap)
	if err != nil {
		// If a configMap name is provided, then it must exist
		// You will likely want to create an Event for the user to understand why their reconcile is failing.
		return err
	}
	fmt.Println(foundConfigMap)
	return nil
}