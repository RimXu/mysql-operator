package controllers

import (

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	//disk := constants.ComboReflect[combo]["Disk"]
	logrus.Infof("combo { cpu:%s, memory:%s }", cpu, memory)
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
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            "mysql",
							Image:           constants.Registry_Addr + constants.Image,
							ImagePullPolicy: "IfNotPresent",
							Env: []apiv1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: constants.MYSQL_ROOT_PASSWORD,
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 3306,
								},
							},
							Resources: apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									"cpu" :  resource.MustParse(cpu),
									"memory": resource.MustParse(memory),
								},
								Requests: apiv1.ResourceList{
									"cpu": resource.MustParse(cpu),
									"memory": resource.MustParse(memory),
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name: "data",
									MountPath: "/data",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "data",
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: m.Name + role + "-data",
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

