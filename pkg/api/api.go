package api

import (
	"context"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	//"k8s.io/client-go/tools/clientcmd"
	//"os"
	//"path/filepath"
	mysqlv1 "mysql-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"mysql-operator/pkg/constants"
	ctrl "sigs.k8s.io/controller-runtime"
)

var gvr = schema.GroupVersionResource{
	constants.Group,
	constants.Version,
	constants.Resource,
}




func InitDynamic() (dynamic.Interface){
	config := ctrl.GetConfigOrDie()
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return client
}

// 通过配置文件获得权限
//func InitDynamic() (dynamic.Interface){
//	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
//	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
//	if err != nil {
//		panic(err)
//	}
//
//	client, err := dynamic.NewForConfig(config)
//	if err != nil {
//		panic(err)
//	}
//	return client
//}



// 使用dynamicClient更新CRD spec
func UpdateMysqlStatus(namespace string, name string, value string) error {
	client := InitDynamic()
	// 获取非结构化数据
	unStructData, err := client.Resource(gvr).Namespace(namespace).Get(context.TODO(),name,metav1.GetOptions{})
	if err != nil {
		logrus.Error(err,"Fetch unstructdata err")
		return err
	}
	logrus.Infof("Fetch unstructdata: { namespace:'%s', name:'%s' }", namespace, name)
	// 将非结构化的数据转换为mysqlv1.Mysql类型
	mysql := &mysqlv1.Mysql{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(
		unStructData.UnstructuredContent(),
		mysql)
	if err != nil {
		logrus.Error(err,"format from unstructdata err")
		return err
	}
	logrus.Infof("Format from unstructdata: { namespace:'%s', name:'%s' }", namespace, name)
	mysql.Spec.Phase = value

	// 将更新完的mysqlv1.Mysql类型转换为unstructdata
	data, err := runtime.DefaultUnstructuredConverter.ToUnstructured(mysql)
	if err != nil {
		logrus.Error(err,"format to unstructdata err")
		return err
	}
	d := unstructured.Unstructured{data}
	logrus.Infof("Format to unstructdata: { namespace:'%s', name:'%s' }", namespace, name)

	// 更新crd
	_, err = client.Resource(gvr).Namespace(namespace).Update(context.TODO(),&d, metav1.UpdateOptions{})
	if err != nil {
		logrus.Errorf("CRD update failed: { namespace:'%s', name:'%s' }", namespace, name)
	}
	return nil

}

