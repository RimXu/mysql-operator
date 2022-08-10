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
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	mysqlv1 "mysql-operator/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	err := r.Get(context.TODO(),req.NamespacedName, mysqloperator)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}
	mysqldep := &appsv1.Deployment{}
	combo := mysqloperator.Spec.Combo
	//判断是否需要创建MySQL主从主从
	if mysqloperator.Spec.Replication == true {
		// 查询{Name:my.Name,Namespace:my.Namespace} 是否存在deployment,如果不存在则满足errors.IsNotFound(err)
		err = r.Get(context.TODO(),types.NamespacedName{Name:mysqloperator.Name+"-master",Namespace:mysqloperator.Namespace},mysqldep)
		if err != nil {
			logrus.Info("MySQL replication configuration")
			if errors.IsNotFound(err) {
				master := r.CreateMysql(mysqloperator,"-master",combo)
				if err = r.Create(context.TODO(), master); err != nil {
					return ctrl.Result{}, err
				}
				//mysqloperator.Status.Replicas = fmt.Sprintf("%d",mysqldep.Status.Replicas)
				if err := r.Status().Update(ctx, mysqloperator); err != nil {
					logrus.Error(err,"MySQL master status update error")
				}
				return ctrl.Result{Requeue: true}, nil
			} else {
				return ctrl.Result{},err
			}
		}
		err = r.Get(context.TODO(),types.NamespacedName{Name:mysqloperator.Name+"-slave",Namespace:mysqloperator.Namespace},mysqldep)
		if err != nil {
			if errors.IsNotFound(err) {
				slave := r.CreateMysql(mysqloperator,"-slave",combo)
				if err = r.Create(context.TODO(), slave); err != nil {
					return ctrl.Result{}, err
				}
				//mysqloperator.Status.Replicas = fmt.Sprintf("%d",mysqldep.Status.Replicas)
				if err := r.Status().Update(ctx, mysqloperator); err != nil {
					logrus.Error(err,"MySQL slave status update error")
				}
				return ctrl.Result{Requeue: true}, nil
			} else {
				return ctrl.Result{},err
			}
		}
	} else {
		logrus.Info("MySQL Standalone configuration")
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MysqlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1.Mysql{}).
		Complete(r)
}
