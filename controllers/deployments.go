package controllers

import (
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/intstr"
	mysqlv1 "mysql-operator/api/v1"
	"mysql-operator/pkg/constants"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// 创建Mysql方法,返回appsv1.deployment类型
func (r *MysqlReconciler) CreateMysql(m *mysqlv1.Mysql, role string, combo string) *appsv1.Deployment {
	labels := LabelsForApp(m.Name)
	var replicas int32 = 1
	cpu := constants.ComboReflect[combo]["CPU"]
	memory := constants.ComboReflect[combo]["Memory"]
	logrus.Infof("MySQL combo { cpu:%s, memory:%s }", cpu, memory)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + role,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            m.Name + role,
							Image:           constants.Registry_Addr + constants.Mysql_Image,
							ImagePullPolicy: "IfNotPresent",
							Args: []string{
								"mysqld",
								"--defaults-file=/etc/mysql/my.cnf",
								"--user=mysql",
							},
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: constants.MYSQL_ROOT_PASSWORD,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 3306,
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"cpu":    resource.MustParse(cpu),
									"memory": resource.MustParse(memory),
								},
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse(cpu),
									"memory": resource.MustParse(memory),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "mysql-data",
									MountPath: "/data",
								},
								{
									Name:      "mysql-config",
									MountPath: "/etc/mysql",
								},
								{
									Name:      "etc-localtime",
									MountPath: "/etc/localtime",
								},
							},
							Lifecycle: &corev1.Lifecycle{
								PostStart: &corev1.LifecycleHandler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"/bin/sh",
											"-c",
											"sh /root/init.sh",
										},
									},
								},
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{
											Type:   0,
											IntVal: 3306,
										},
									},
								},
								TimeoutSeconds:      5,
								SuccessThreshold:    1,
								FailureThreshold:    3,
								InitialDelaySeconds: 15,
								PeriodSeconds:       30,
							},
						},
						// exporter container
						{
							Name:            "mysql-exporter",
							Image:           constants.Registry_Addr + constants.Exporter_Image,
							ImagePullPolicy: "IfNotPresent",
							Args: []string{
								"--collect.binlog_size",
								"--collect.engine_innodb_status",
							},
							Env: []corev1.EnvVar{
								{
									Name:  "DATA_SOURCE_NAME",
									Value: "exporter:exporterMWF@(localhost:3306)/",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 9104,
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"cpu":    resource.MustParse("500m"),
									"memory": resource.MustParse("500Mi"),
								},
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse("100m"),
									"memory": resource.MustParse("200Mi"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "etc-localtime",
									MountPath: "/etc/localtime",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "mysql-data",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: m.Name + role + "-data",
								},
							},
						},
						{
							Name: "mysql-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: m.Name + role,
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "my.cnf",
											Path: "my.cnf",
										},
									},
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
	// 设置deployment的上级控制器
	err := controllerutil.SetControllerReference(m, deployment, r.Scheme)
	if err != nil {
		logrus.Error(err, "MySQL created failed", " Name:", deployment.Name, " Namespace:", deployment.Namespace)
	}
	logrus.Infof("MySQL created successful { name:%s, namespace:%s }", deployment.Name, deployment.Namespace)
	return deployment
}

func LabelsForApp(name string) map[string]string {
	return map[string]string{"app": name}
}
