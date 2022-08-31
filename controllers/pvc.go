package controllers

import (
	"context"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

// 查询pvc是否存在,如果不存在返回 PersistentVolume "XXX" not found 错误类型
func (r *MysqlReconciler) QueryPVC(ns string, sc string, name string, ctx context.Context) error {
	//
	foundPVC := &corev1.PersistentVolumeClaim{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, foundPVC)
	if err != nil {
		logrus.Warn(err)
		return err
	}
	logrus.Errorf(" PersistentVolumeClaim exists: { namespace: %s, name : %s }",ns,name)
	return nil
}

// 创建PVC
func (r *MysqlReconciler) CreatePVC(ns string, sc string, name string, size string,ctx context.Context) error {
	logrus.Infof("PersistentVolumeClaim creating: { namespace:'%s', name:'%s' }",ns,name)
	optionPVC := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Namespace: ns,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{Requests: corev1.ResourceList{corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(size),
			},
			},
			StorageClassName: &sc,
		},
	}
	err := r.Create(context.TODO(),optionPVC)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("PersistentVolumeClaim created successful { Namespace : %s, name : %s, size : %s }",ns,name,size)
	return nil
}
