package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	mysqlv1 "mysql-operator/api/v1"
)

// 查询MySQL SVC
func (r *MysqlReconciler) QueryMysqlSVC(ns string, name string, ctx context.Context) error {
	foundSVC := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, foundSVC)
	if err == nil {
		msg := fmt.Sprintf("MySQL service exists: { namespace: %s, name : %s }",ns,name)
		errmsg := errors.New(msg)
		logrus.Error(errmsg)
		return errmsg
	}
	logrus.Warnf("MySQL service not found: { namespace: %s, name : %s }",ns,name)
	return err
}

// 查询ProxySQL SVC
func (r *MysqlReconciler) QueryProxySVC(ns string, name string, ctx context.Context) error {
	foundSVC := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, foundSVC)
	if err == nil {
		msg := fmt.Sprintf("Proxy service exists: { namespace: %s, name : %s }",ns,name)
		errmsg := errors.New(msg)
		logrus.Error(errmsg)
		return errmsg
	}
	logrus.Warnf("Proxy service not found: { namespace: %s, name : %s }",ns,name)
	return err
}


// 创建Mysql SVC
func (r *MysqlReconciler) CreateMysqlSVC(m *mysqlv1.Mysql,ns string, name string, ctx context.Context) error {
	logrus.Infof("Mysql service creating: { namespace:'%s', name:'%s' }",ns,name)
	optionPVC := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Namespace: ns,
			Labels: map[string]string{
				"name":  name,
				"system/appName": name,
				"system/svcName": name,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "tcp-mysql",
					Protocol: "TCP",
					Port: 3306,
					TargetPort: intstr.IntOrString{
						Type: 0,
						IntVal: 3306,
					},
				},
				{
					Name: "tcp-exporter",
					Protocol: "TCP",
					Port: 9104,
					TargetPort: intstr.IntOrString{
						Type: 0,
						IntVal: 9104,
					},
				},
			},
			Selector: map[string]string {
				"name": name,
			},
			Type: "ClusterIP",
			ClusterIP: "None",
		},
	}
	// 设置Service的上级控制器
	err := controllerutil.SetControllerReference(m, optionPVC, r.Scheme)
	if err != nil {
		logrus.Errorf("Service set controller failed { namespace:'%s', name:'%s' }",ns,name)
	}

	err = r.Create(context.TODO(),optionPVC)
	if err != nil {
		logrus.Error(err)
		return err
	}
	// 更新Mysqloperator状态
	m.Status.Status = "svc"
	if err = r.Status().Update(ctx, m); err != nil {
		logrus.Error(err, "Operator status update error")
		return err
	}
	logrus.Infof("Service created successful { Namespace : %s, name : %s }",ns,name)
	return nil
}

// 创建ProxySQL SVC
func (r *MysqlReconciler) CreateProxySVC(m *mysqlv1.Mysql,ns string, name string, ctx context.Context) error {
	logrus.Infof("Proxy service creating: { namespace:'%s', name:'%s' }",ns,name)
	optionPVC := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Namespace: ns,
			Labels: map[string]string{
				"name":  name,
				"system/appName": name,
				"system/svcName": name,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "tcp-proxysql",
					Protocol: "TCP",
					Port: 6033,
					TargetPort: intstr.IntOrString{
						Type: 0,
						IntVal: 6033,
					},
				},
			},
			Selector: map[string]string {
				"name": name,
			},
			Type: "ClusterIP",
		},
	}
	// 设置Service的上级控制器
	err := controllerutil.SetControllerReference(m, optionPVC, r.Scheme)
	if err != nil {
		logrus.Errorf("Proxy service set controller failed { namespace:'%s', name:'%s' }",ns,name)
	}

	err = r.Create(context.TODO(),optionPVC)
	if err != nil {
		logrus.Error(err)
		return err
	}
	// 更新Mysqloperator状态
	m.Status.Status = "svc"
	if err = r.Status().Update(ctx, m); err != nil {
		logrus.Error(err, "Operator status update error")
		return err
	}
	logrus.Infof("Proxy service created successful { Namespace : %s, name : %s }",ns,name)
	return nil
}
