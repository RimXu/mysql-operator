package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"mysql-operator/pkg/constants"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strconv"
	"strings"
	mysqlv1 "mysql-operator/api/v1"
)




// 查询ConfigMaps,如果不存在则返回错误
func (r *MysqlReconciler) QueryMysqlCM(ns string, name string, ctx context.Context) error {
	foundCM := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, foundCM)
	if err == nil {
		msg := fmt.Sprintf("MySQL Configmaps exists: { namespace: %s, name : %s }",ns,name)
		errmsg := errors.New(msg)
		logrus.Error(errmsg)
		return errmsg
	}
	logrus.Warnf("MySQL Configmaps not found: { namespace: %s, name : %s }",ns,name)
	return err
}

// 创建ConfigMaps,如果不存在则返回错误
func (r *MysqlReconciler) CreateMysqlCM(m *mysqlv1.Mysql,ns string, name string, role string, combo string, ctx context.Context) error {
	logrus.Infof("MySQL ConfigMaps creating: { namespace:'%s', name:'%s' }",ns,name)
	var server_id string
	if find := strings.Contains(name, "master"); find {
		server_id = "10"
	} else {
		server_id = "11"
	}
	config_cm,_:= ReadMycnf(constants.MySQLCfg,server_id,FormatBufferpool(constants.ComboReflect[combo]["Memory"]))
	init_cm := constants.InitCfg
	optionCM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Namespace: ns,
		},
		Data: map[string]string{
			"my.cnf": config_cm,
			"init.sh": init_cm,
		},

	}
	// 设置Configmaps的上级控制器
	err := controllerutil.SetControllerReference(m, optionCM, r.Scheme)
	if err != nil {
		logrus.Errorf("MySQL Configmaps set controller failed { namespace:'%s', name:'%s' }",ns,name)
	}

	err = r.Create(context.TODO(),optionCM)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("MySQL ConfigMaps created successful { Namespace : %s, name : %s }",ns,name)
	return nil
}


// 查询ProxySQL ConfigMaps,如果不存在则返回错误
func (r *MysqlReconciler) QueryProxyCM(ns string, name string, ctx context.Context) error {
	foundCM := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, foundCM)
	if err == nil {
		msg := fmt.Sprintf("ProxySQL Configmaps exists: { namespace: %s, name : %s }",ns,name)
		errmsg := errors.New(msg)
		logrus.Error(errmsg)
		return errmsg
	}
	logrus.Warnf("ProxySQL Configmaps not found: { namespace: %s, name : %s }",ns,name)
	return err
}



// 创建Proxy ConfigMaps,如果不存在则返回错误
func (r *MysqlReconciler) CreateProxyCM(m *mysqlv1.Mysql,ns string, name string, ctx context.Context) error {
	logrus.Infof("Proxy ConfigMaps creating: { namespace:'%s', name:'%s' }",ns,name)
	cm,_:= ReadProxycnf(constants.ProxyCfg,ns, name,"mydb","mydbpass")
	optionCM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name + "-proxy",
			Namespace: ns,
		},
		Data: map[string]string{
			"proxysql.cnf": cm,
		},

	}
	// 设置Configmaps的上级控制器
	err := controllerutil.SetControllerReference(m, optionCM, r.Scheme)
	if err != nil {
		logrus.Errorf("Proxy Configmaps set controller failed { namespace:'%s', name:'%s' }",ns,name)
	}

	err = r.Create(context.TODO(),optionCM)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("Proxy ConfigMaps created successful { Namespace : %s, name : %s }",ns,name)
	return nil
}

// 查询init ConfigMaps,如果不存在则返回错误
func (r *MysqlReconciler) QueryInitCM(ns string, name string, ctx context.Context) error {
	foundCM := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, foundCM)
	if err == nil {
		msg := fmt.Sprintf("Init Configmaps exists: { namespace: %s, name : %s }",ns,name)
		errmsg := errors.New(msg)
		logrus.Error(errmsg)
		return errmsg
	}
	logrus.Warnf("Init Configmaps not found: { namespace: %s, name : %s }",ns,name)
	return err
}


// 格式化bufferpool配置
func FormatBufferpool(m string) string {
	if find := strings.Contains(m, "M"); find {
		memory := strings.Split(m,"M")
		intmemory, _ := strconv.ParseFloat(memory[0],32)
		intbuffermem := intmemory * 0.6
		buffermem, _ := strconv.Atoi(fmt.Sprintf("%1.0f", intbuffermem))
		strmem := fmt.Sprintf("%dM",buffermem)
		return strmem

	} else if find := strings.Contains(m, "G"); find {
		memory := strings.Split(m,"G")
		intmemory, _ := strconv.ParseFloat(memory[0],32)
		intbuffermem := intmemory * 0.6
		buffermem, _ := strconv.Atoi(fmt.Sprintf("%1.0f", intbuffermem))
		strmem := fmt.Sprintf("%dG",buffermem)
		return strmem
	}
	return "1G"
}



// 格式化Mysql配置文件
func ReadMycnf(s string, id string, buffer string) (string,error) {
	cfg1 := strings.Replace(s,"MYSQL_SERVER_ID",id,-1)
	cfg2 := strings.Replace(cfg1,"MYSQL_BUFFER_POOL_SIZE",buffer,-1)
	return cfg2,nil
}

// 格式化ProxySQL配置文件
func ReadProxycnf(s string,ns string, name string, user string, password string) (string,error){
	cfg1 := strings.Replace(s, "NAMESPACE", ns, -1)
	cfg2 := strings.Replace(cfg1, "MYSQL-ADDR", name, -1)
	cfg3 := strings.Replace(cfg2, "MYSQL-USERNAME", user, -1)
	cfg4 := strings.Replace(cfg3, "MYSQL-PASSWORD", password,-1)
	return cfg4,nil
}

