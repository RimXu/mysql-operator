package controllers

import (
	"context"
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	mysqlv1 "mysql-operator/api/v1"
	"mysql-operator/pkg/constants"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// 创建主从MySQL初始化job
func (r *MysqlReconciler) CreateRepJob(m *mysqlv1.Mysql, ns string, name string, database string, id string, ctx context.Context) error {
	job_name := name + "-job-" + id
	proxy_name := name + "-proxy"
	args := []string{
		"/bin/sh",
		"-c",
		"/etc/mysql/mysqljob.sh " + database,
	}
	logrus.Infof("Job creating: { namespace:'%s', name:'%s' }", ns, job_name)
	//args = ["/bin/sh", "-c"]
	//for _,db := range databases {
	//	sh := fmt.Sprintf("/etc/mysql/mysqljob.sh %s %s %s;",db["name"],db["user"],db["passwd"])
	//	db_argss = append(db_argss,sh)
	//}
	var activeDeadlineSeconds int64 = 900
	var backoffLimit int32 = 3
	var defaultMode int32 = 0755
	optionJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      job_name,
			Namespace: ns,
			Labels: map[string]string{
				"name":           job_name,
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
							Name:            job_name,
							Image:           constants.Registry_Addr + constants.Mysql_Image,
							ImagePullPolicy: "IfNotPresent",
							Args:            args,
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: constants.MYSQL_ROOT_PASSWORD,
								},
								{
									Name:  "m_mysql",
									Value: name + "-master" + "." + ns,
								},
								{
									Name:  "s_mysql",
									Value: name + "-slave" + "." + ns,
								},
								{
									Name:  "MONITOR_USER",
									Value: constants.Monitor_User,
								},
								{
									Name:  "MONITOR_PASS",
									Value: constants.Monitor_Password,
								},
								{
									Name:  "EXPORTER_USER",
									Value: constants.Exporter_User,
								},
								{
									Name:  "EXPORTER_PASS",
									Value: constants.Exporter_Password,
								},
								{
									Name:  "REPL_USER",
									Value: constants.Repl_User,
								},
								{
									Name:  "REPL_PASS",
									Value: constants.Exporter_Password,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "mysql-config",
									MountPath: "/etc/mysql",
								},
								{
									Name:      "etc-localtime",
									MountPath: "/etc/localtime",
								},
							},
						},
					},
					RestartPolicy:         "Never",
					ActiveDeadlineSeconds: &activeDeadlineSeconds,
					Volumes: []corev1.Volume{
						{
							Name: "mysql-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: proxy_name,
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "mysqljob.sh",
											Path: "mysqljob.sh",
										},
									},
									DefaultMode: &defaultMode,
								},
							},
						},
						{
							Name: "etc-localtime",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/etc/localtime",
								},
							},
						},
					},
				},
			},
		},
	}
	// 设置Job的上级控制器
	err := controllerutil.SetControllerReference(m, optionJob, r.Scheme)
	if err != nil {
		logrus.Errorf("Job set controller failed { namespace:'%s', name:'%s' }", ns, name)
	}
	err = r.Create(context.TODO(), optionJob)
	if err != nil {
		logrus.Error(err)
		return err
	}

	m.Spec.Phase = "JobCreated"
	if err := r.Client.Update(ctx, m); err != nil {
		logrus.Error(err, "Operator status update error")
		return err
	}

	logrus.Infof("Job created successful { name:%s, namespace:%s }", job_name, ns)
	return nil
}

// 创建Single MySQL初始化job
func (r *MysqlReconciler) CreateSingleJob(m *mysqlv1.Mysql, ns string, name string, database string, id string, ctx context.Context) error {
	job_name := name + "-job-" + id
	args := []string{
		"/bin/sh",
		"-c",
		"/etc/mysql/mysqljob.sh " + database,
	}
	logrus.Infof("Job creating: { namespace:'%s', name:'%s' }", ns, job_name)
	//args = ["/bin/sh", "-c"]
	//for _,db := range databases {
	//	sh := fmt.Sprintf("/etc/mysql/mysqljob.sh %s %s %s;",db["name"],db["user"],db["passwd"])
	//	db_argss = append(db_argss,sh)
	//}
	var activeDeadlineSeconds int64 = 900
	var backoffLimit int32 = 3
	var defaultMode int32 = 0755
	optionJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      job_name,
			Namespace: ns,
			Labels: map[string]string{
				"name":           job_name,
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
							Name:            job_name,
							Image:           constants.Registry_Addr + constants.Mysql_Image,
							ImagePullPolicy: "IfNotPresent",
							Args:            args,
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: constants.MYSQL_ROOT_PASSWORD,
								},
								{
									Name:  "m_mysql",
									Value: name + "." + ns,
								},
								{
									Name:  "MONITOR_USER",
									Value: constants.Monitor_User,
								},
								{
									Name:  "MONITOR_PASS",
									Value: constants.Monitor_Password,
								},
								{
									Name:  "EXPORTER_USER",
									Value: constants.Exporter_User,
								},
								{
									Name:  "EXPORTER_PASS",
									Value: constants.Exporter_Password,
								},
								{
									Name:  "REPL_USER",
									Value: constants.Repl_User,
								},
								{
									Name:  "REPL_PASS",
									Value: constants.Exporter_Password,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "mysql-config",
									MountPath: "/etc/mysql",
								},
								{
									Name:      "etc-localtime",
									MountPath: "/etc/localtime",
								},
							},
						},
					},
					RestartPolicy:         "Never",
					ActiveDeadlineSeconds: &activeDeadlineSeconds,
					Volumes: []corev1.Volume{
						{
							Name: "mysql-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: name,
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "mysqljob.sh",
											Path: "mysqljob.sh",
										},
									},
									DefaultMode: &defaultMode,
								},
							},
						},
						{
							Name: "etc-localtime",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/etc/localtime",
								},
							},
						},
					},
				},
			},
		},
	}
	// 设置Job的上级控制器
	err := controllerutil.SetControllerReference(m, optionJob, r.Scheme)
	if err != nil {
		logrus.Errorf("Job set controller failed { namespace:'%s', name:'%s' }", ns, name)
	}
	err = r.Create(context.TODO(), optionJob)
	if err != nil {
		logrus.Error(err)
		return err
	}

	m.Spec.Phase = "JobCreated"
	if err := r.Client.Update(ctx, m); err != nil {
		logrus.Error(err, "Operator status update error")
		return err
	}

	logrus.Infof("Job created successful { name:%s, namespace:%s }", job_name, ns)
	return nil
}

//func (r MysqlReconciler) UpdateSpecStatus(ctx context.Context, m *mysqlv1.Mysql) error {
//	if err := r.Client.Update(ctx, m); err != nil {
//		logrus.Error(err, "Operator status update error")
//		return err
//	}
//	return nil
//}
