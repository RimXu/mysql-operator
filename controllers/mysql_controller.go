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
	//"sigs.k8s.io/controller-runtime/pkg/source"
	"strconv"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	//"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	mysqlv1 "mysql-operator/api/v1"
	"mysql-operator/pkg/constants"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	batchv1 "k8s.io/api/batch/v1"
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
	//logrus.Info("MySQL-Operator reconciler start ",ctx)
	mysqloperator := &mysqlv1.Mysql{}

	// 查询Namespace下是否存在mysqloperator,如果不存在则满足errors.IsNotFound(err),函数返回
	err := r.Get(context.TODO(), req.NamespacedName, mysqloperator)

	// uuid 判空说明是删除namespace下指定的operator
	//uuid := mysqloperator.ObjectMeta.UID
	//if uuid == "" {
	//	logrus.Info("MySQL-Operator reconciler delete")
	//}
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}
	mysqldep := &appsv1.Deployment{}
	instance := mysqloperator.Spec.Instance
	sc := mysqloperator.Spec.StorageClass
	replication := mysqloperator.Spec.Replication
	databases := mysqloperator.Spec.Databases
	size := constants.InstanceReflect[instance]["Disk"]

	//判断是否需要创建MySQL主从主从
	if replication == true {
		// 查询{Name:my.Name,Namespace:my.Namespace} 是否存在deployment,如果不存在则满足errors.IsNotFound(err)
		err = r.Get(context.TODO(), types.NamespacedName{Name: mysqloperator.Name + "-master", Namespace: mysqloperator.Namespace}, mysqldep)
		if err != nil {
			// 如果mysqloperator不存在
			if errors.IsNotFound(err) {
				// 如果CM/SVC/PVC均不存在，则继续CreateCM/SVC/PVC
				mcmerr := r.QueryMysqlCM(mysqloperator.Namespace, mysqloperator.Name+"-master", ctx)
				msvcerr := r.QueryMysqlSVC(mysqloperator.Namespace, mysqloperator.Name+"-master", ctx)
				mpvcerr := r.QueryMysqlPVC(mysqloperator.Namespace, sc, mysqloperator.Name+"-master-data", ctx)
				scmerr := r.QueryMysqlCM(mysqloperator.Namespace, mysqloperator.Name+"-slave", ctx)
				ssvcerr := r.QueryMysqlSVC(mysqloperator.Namespace, mysqloperator.Name+"-slave", ctx)
				spvcerr := r.QueryMysqlPVC(mysqloperator.Namespace, sc, mysqloperator.Name+"-slave-data", ctx)
				pcmerr := r.QueryMysqlCM(mysqloperator.Namespace, mysqloperator.Name+"-proxy", ctx)
				psvcerr := r.QueryProxySVC(mysqloperator.Namespace, mysqloperator.Name+"-proxy", ctx)
				if errors.IsNotFound(mcmerr) && errors.IsNotFound(msvcerr) && errors.IsNotFound(mpvcerr) &&
					errors.IsNotFound(scmerr) && errors.IsNotFound(ssvcerr) && errors.IsNotFound(spvcerr) &&
					errors.IsNotFound(pcmerr) && errors.IsNotFound(psvcerr) {
					// Create Mysql CM
					err = r.CreateRepMysqlCM(mysqloperator, mysqloperator.Namespace, mysqloperator.Name+"-master", "master", instance, ctx)
					if err != nil {
						logrus.Error("CreateMysqlCM error", err)
					}
					err = r.CreateRepMysqlCM(mysqloperator, mysqloperator.Namespace, mysqloperator.Name+"-slave", "slave", instance, ctx)
					if err != nil {
						logrus.Error("CreateMysqlCM error", err)
					}
					// Create Proxysql CM
					err = r.CreateProxyCM(mysqloperator, mysqloperator.Namespace, mysqloperator.Name, databases, ctx)
					if err != nil {
						logrus.Error("CreateProxyCM error", err)
					}

					// Create MySQL SVC
					err = r.CreateRepMysqlSVC(mysqloperator, mysqloperator.Namespace, mysqloperator.Name+"-master", ctx)
					if err != nil {
						logrus.Error("CreateMysqlSVC error", err)
					}
					err = r.CreateRepMysqlSVC(mysqloperator, mysqloperator.Namespace, mysqloperator.Name+"-slave", ctx)
					if err != nil {
						logrus.Error("CreateMysqlSVC error", err)
					}
					// Create Proxy SVC
					err = r.CreateProxySVC(mysqloperator, mysqloperator.Namespace, mysqloperator.Name+"-proxy", ctx)
					if err != nil {
						logrus.Error("CreateProxySVC err", err)
					}

					// CreatePVC
					err = r.CreateMysqlPVC(mysqloperator.Namespace, sc, mysqloperator.Name+"-master-data", size, ctx)
					if err != nil {
						logrus.Error("CreatePVC error", err)
					}
					err = r.CreateMysqlPVC(mysqloperator.Namespace, sc, mysqloperator.Name+"-slave-data", size, ctx)
					if err != nil {
						logrus.Error("CreatePVC error", err)
					}
				} else {
					return ctrl.Result{Requeue: false}, nil
				}

				// 创建MySQL master
				_, err = r.CreateMysql(mysqloperator, "-master", instance, ctx)
				if err != nil {
					return ctrl.Result{}, err
				}

				// 创建MySQL slave
				_, err = r.CreateMysql(mysqloperator, "-slave", instance, ctx)
				if err != nil {
					return ctrl.Result{}, err
				}

				// 创建ProxySQL
				_, err = r.CreateProxy(mysqloperator, ctx)
				if err != nil {
					logrus.Error(err, "Proxy create error")
					return ctrl.Result{}, err
				}

				// 创建初始化主从Job
				for id, db := range databases {
					args := fmt.Sprintf("%s %s %s;", db["name"], db["user"], db["passwd"])
					err = r.CreateRepJob(mysqloperator, mysqloperator.Namespace, mysqloperator.Name, args, strconv.Itoa(id), ctx)
					if err != nil {
						return ctrl.Result{}, nil
					}
				}
				return ctrl.Result{Requeue: true}, nil
			}
		}
	} else if replication == false {
		err = r.Get(context.TODO(), types.NamespacedName{Name: mysqloperator.Name, Namespace: mysqloperator.Namespace}, mysqldep)
		if err != nil {
			// 如果mysqloperator不存在
			if errors.IsNotFound(err) {
				cmerr := r.QueryMysqlCM(mysqloperator.Namespace, mysqloperator.Name, ctx)
				svcerr := r.QueryMysqlSVC(mysqloperator.Namespace, mysqloperator.Name, ctx)
				pvcerr := r.QueryMysqlPVC(mysqloperator.Namespace, sc, mysqloperator.Name, ctx)
				if errors.IsNotFound(cmerr) && errors.IsNotFound(svcerr) && errors.IsNotFound(pvcerr)  {
					err = r.CreateSingleMysqlCM(mysqloperator, mysqloperator.Namespace, mysqloperator.Name, "", instance, ctx)
					if err != nil {
						logrus.Error("CreateMysqlCM error", err)
					}
					// Create MySQL SVC
					err = r.CreateSingleSVC(mysqloperator, mysqloperator.Namespace, mysqloperator.Name, ctx)
					if err != nil {
						logrus.Error("CreateMysqlSVC error", err)
					}
					// CreatePVC
					err = r.CreateMysqlPVC(mysqloperator.Namespace, sc, mysqloperator.Name+"-data", size, ctx)
					if err != nil {
						logrus.Error("CreatePVC error", err)
					}
					// 创建MySQL master
					_, err = r.CreateMysql(mysqloperator, "", instance, ctx)
					if err != nil {
						return ctrl.Result{}, err
					}

					// 创建初始化主从Job
					for id, db := range databases {
						args := fmt.Sprintf("%s %s %s;", db["name"], db["user"], db["passwd"])
						err = r.CreateSingleJob(mysqloperator, mysqloperator.Namespace, mysqloperator.Name, args, strconv.Itoa(id), ctx)
						if err != nil {
							return ctrl.Result{}, nil
						}
					}
					return ctrl.Result{}, nil
				}
			}
		}
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
		//Watches(&source.Kind{Type: &batchv1.Job{}}, &MysqlJob{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}


