package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"k8s.io/apimachinery/pkg/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	mysqlv1 "mysql-operator/api/v1"
)

func (r *MysqlReconciler) QuerySVC(ns string, name string, ctx context.Context) error {
	foundSVC := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, foundSVC)
	if err == nil {
		msg := fmt.Sprintf("Service exists: { namespace: %s, name : %s }",ns,name)
		errmsg := errors.New(msg)
		logrus.Error(errmsg)
		return errmsg
	}
	logrus.Warnf("Service not found: { namespace: %s, name : %s }",ns,name)
	return err
}


// 创建SVC
func (r *MysqlReconciler) CreateSVC(m *mysqlv1.Mysql,ns string, name string, ctx context.Context) error {
	logrus.Infof("PersistentVolumeClaim creating: { namespace:'%s', name:'%s' }",ns,name)
	optionPVC := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Namespace: ns,
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
