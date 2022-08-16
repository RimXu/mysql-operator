package controllers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	mysqlv1 "mysql-operator/api/v1"
	"mysql-operator/pkg/constants"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// 创建Mysql方法,返回appsv1.deployment类型
func (r *MysqlReconciler) CreateMysql(m *mysqlv1.Mysql,role string,combo string) *appsv1.Deployment {
	labels := LabelsForApp(m.Name)
	var replicas int32 = 1
	var cpu,memory,disk int
	switch combo {
	case "small":
		cpu = constants.Small_CPU
		memory = constants.Small_Memory
		disk = constants.Small_Disk
	case "medium":
		cpu = constants.Medium_CPU
		memory = constants.Medium_Memory
		disk = constants.Medium_Disk
	case "large":
		cpu = constants.Large_CPU
		memory = constants.Large_Memory
		disk = constants.Large_Disk
	}
	fmt.Print(cpu,memory,disk)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.Name+role,
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
							Name: "mysql",
							Image: constants.Registry_Addr + constants.Image,
							ImagePullPolicy: "IfNotPresent",
							Env: []apiv1.EnvVar{
								{
									Name: "MYSQL_ROOT_PASSWORD",
									Value: constants.MYSQL_ROOT_PASSWORD,
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									Name: "http",
									Protocol: apiv1.ProtocolTCP,
									ContainerPort: 3306,
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
		logrus.Error(err,"MySQL created failed"," Name:",deployment.Name," Namespace:",deployment.Namespace)
	}
	logrus.Infof("MySQL created successful Name:%s,Namespace:%s",deployment.Name,deployment.Namespace)
	return deployment
}


func LabelsForApp(name string) map[string]string {
	return map[string]string{"app":name}
}
