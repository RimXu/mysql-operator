package controllers

import (
	"context"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/util/intstr"
	mysqlv1 "mysql-operator/api/v1"
	"mysql-operator/pkg/constants"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// 创建Mysql方法,返回appsv1.deployment类型
func (r *MysqlReconciler) CreateMysql(m *mysqlv1.Mysql, role string, instance string, ctx context.Context) (*appsv1.Deployment, error) {
	var replicas int32 = 1
	mysql_name := m.Name + role
	cpu := constants.InstanceReflect[instance]["CPU"]
	memory := constants.InstanceReflect[instance]["Memory"]
	logrus.Infof("MySQL instance { cpu:%s, memory:%s }", cpu, memory)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysql_name,
			Namespace: m.Namespace,
			Labels: map[string]string{
				"app":            mysql_name,
				"name":           mysql_name,
				"system/appName": mysql_name,
				"system/svcName": mysql_name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": mysql_name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":            mysql_name,
						"name":           mysql_name,
						"system/appName": mysql_name,
						"system/svcName": mysql_name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            mysql_name,
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
											"sh /etc/mysql/init.sh",
										},
									},
								},
							},
							// Type 0表示整型
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
									Value: constants.Exporter_User + ":" + constants.Exporter_Password + "@(" + constants.Exporter_Reg + ":3306)/",
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
									"cpu":    resource.MustParse(constants.Exporter_CPU_lim),
									"memory": resource.MustParse(constants.Exporter_Mem_lim),
								},
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse(constants.Exporter_CPU_req),
									"memory": resource.MustParse(constants.Exporter_Mem_req),
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
									ClaimName: mysql_name + "-data",
								},
							},
						},
						{
							Name: "mysql-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: mysql_name,
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "my.cnf",
											Path: "my.cnf",
										},
										{
											Key:  "init.sh",
											Path: "init.sh",
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
			Strategy: appsv1.DeploymentStrategy{
				Type: "Recreate",
			},
		},
	}
	// 设置deployment的上级控制器
	err := controllerutil.SetControllerReference(m, deployment, r.Scheme)
	if err != nil {
		logrus.Error(err, "MySQL created failed", " Name:", deployment.Name, " Namespace:", deployment.Namespace)
	}

	err = r.Create(context.TODO(), deployment)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// 更新Mysqloperator状态
	m.Status.Status = "DeploymentCompleted"
	if err = r.Status().Update(ctx, m); err != nil {
		logrus.Error(err, "Operator status update error")
		return nil, err
	}

	logrus.Infof("MySQL created successful { name:%s, namespace:%s }", deployment.Name, deployment.Namespace)
	return deployment, nil
}

// 创建proxysql返回deployments类型
func (r *MysqlReconciler) CreateProxy(m *mysqlv1.Mysql, ctx context.Context) (*appsv1.Deployment, error) {
	prxoy_name := m.Name + "-proxy"
	var replicas int32 = 2
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      prxoy_name,
			Namespace: m.Namespace,
			Labels: map[string]string{
				"app":            prxoy_name,
				"name":           prxoy_name,
				"system/appName": prxoy_name,
				"system/svcName": prxoy_name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": prxoy_name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":            prxoy_name,
						"name":           prxoy_name,
						"system/appName": prxoy_name,
						"system/svcName": prxoy_name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            prxoy_name,
							Image:           constants.Registry_Addr + constants.ProxySQL_Image,
							ImagePullPolicy: "IfNotPresent",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 6033,
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"cpu":    resource.MustParse(constants.Proxy_CPU_req),
									"memory": resource.MustParse(constants.Proxy_Mem_req),
								},
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse(constants.Proxy_CPU_req),
									"memory": resource.MustParse(constants.Proxy_Mem_req),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "proxy-localtime",
									MountPath: "/etc/localtime",
								},
								{
									Name:      "proxy-config",
									MountPath: "/etc/proxysql.cnf",
									SubPath:   "proxysql.cnf",
								},
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{
											Type:   0,
											IntVal: 6033,
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
					},
					Volumes: []corev1.Volume{
						{
							Name: "proxy-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: m.Name + "-proxy",
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "proxysql.cnf",
											Path: "proxysql.cnf",
										},
									},
								},
							},
						},
						{
							Name: "proxy-localtime",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/etc/localtime",
								},
							},
						},
					},
				},
			},
			// Type 1表示字符串
			Strategy: appsv1.DeploymentStrategy{
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   1,
						StrVal: "25%",
					},
					MaxSurge: &intstr.IntOrString{
						Type:   1,
						StrVal: "25%",
					},
				},
				Type: "RollingUpdate",
			},
		},
	}
	// 设置deployment的上级控制器
	err := controllerutil.SetControllerReference(m, deployment, r.Scheme)
	if err != nil {
		logrus.Error(err, "ProxySQL created failed", " Name:", deployment.Name, " Namespace:", deployment.Namespace)
	}

	err = r.Create(context.TODO(), deployment)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// 更新Mysqloperator状态
	m.Status.Status = "ProxyCompleted"
	if err = r.Status().Update(ctx, m); err != nil {
		logrus.Error(err, "Operator status update error")
		return nil, err
	}

	logrus.Infof("ProxySQL created successful { name:%s, namespace:%s }", deployment.Name, deployment.Namespace)
	return deployment, err
}
