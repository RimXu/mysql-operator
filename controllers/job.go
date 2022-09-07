package controllers

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"mysql-operator/pkg/constants"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	mysqlv1 "mysql-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

// 创建MySQL初始化job
func (r *MysqlReconciler) CreateRepJob(m *mysqlv1.Mysql,ns string, name string, databases []map[string]string,ctx context.Context) error{
	logrus.Infof("Job creating: { namespace:'%s', name:'%s' }",ns,name)
	job_name := name + "-mysqlinit"
	var activeDeadlineSeconds int64 = 900
	var backoffLimit int32 = 3
	optionJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: job_name,
			Namespace: ns,
			Labels: map[string]string{
				"name":  job_name,
				"system/appName": job_name,
				"system/svcName": job_name,
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: job_name,
							Image: constants.Registry_Addr + constants.Mysql_Image,
							ImagePullPolicy: "IfNotPresent",
							Args: []string{
								"/bin/sh",
								"-c",
								"sleep 1000",
							},
							Env: []corev1.EnvVar{
								{
									Name: "m_mysql",
									Value: name + "-master" + "." + ns,
								},
								{
									Name: "s_mysql",
									Value: name + "-slave" + "." + ns,
								},
								{
									Name: "MONITOR_USER",
									Value: constants.Monitor_User,
								},
								{
									Name: "MONITOR_PASS",
									Value: constants.Monitor_Password,
								},
								{
									Name: "EXPORTER_USER",
									Value: constants.Exporter_User,
								},
								{
									Name: "EXPORTER_PASS",
									Value: constants.Exporter_Password,
								},
								{
									Name: "REPL_USER",
									Value: constants.Repl_User,
								},
								{
									Name: "REPL_PASS",
									Value: constants.Exporter_Password,
								},
							},

						},
					},
					RestartPolicy: "Never",
					ActiveDeadlineSeconds: &activeDeadlineSeconds,
				},
			},
		},
	}
	// 设置Job的上级控制器
	err := controllerutil.SetControllerReference(m, optionJob, r.Scheme)
	if err != nil {
		logrus.Errorf("Job set controller failed { namespace:'%s', name:'%s' }",ns,name)
	}
	err = r.Create(context.TODO(),optionJob)
	if err != nil {
		logrus.Error(err)
		return err
	}
	fmt.Println(ns,name,databases)
	return nil
}