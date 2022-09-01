package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	mysqlv1 "mysql-operator/api/v1"
	"mysql-operator/pkg/constants"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strconv"
	"strings"
)

// 查询ConfigMaps,如果不存在则返回错误
func (r *MysqlReconciler) QueryCM(ns string, name string, ctx context.Context) error {
	foundCM := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, foundCM)
	if err == nil {
		msg := fmt.Sprintf("Configmaps exists: { namespace: %s, name : %s }", ns, name)
		errmsg := errors.New(msg)
		logrus.Error(errmsg)
		return errmsg
	}
	logrus.Warnf("Configmaps not found: { namespace: %s, name : %s }", ns, name)
	return err
}

// 创建ConfigMaps,如果不存在则返回错误
func (r *MysqlReconciler) CreateCM(m *mysqlv1.Mysql, ns string, name string, role string, combo string, ctx context.Context) error {
	logrus.Infof("ConfigMaps creating: { namespace:'%s', name:'%s' }", ns, name)
	var server_id string
	if find := strings.Contains(name, "master"); find {
		server_id = "10"
	} else {
		server_id = "11"
	}
	cm, _ := ReadMycnf(constants.MySQLCfg, server_id, FormatBufferpool(constants.ComboReflect[combo]["Memory"]))
	optionCM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Data: map[string]string{
			"my.cnf": cm,
		},
	}
	// 设置Configmaps的上级控制器
	err := controllerutil.SetControllerReference(m, optionCM, r.Scheme)
	if err != nil {
		logrus.Errorf("Configmaps set controller failed { namespace:'%s', name:'%s' }", ns, name)
	}

	err = r.Create(context.TODO(), optionCM)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("ConfigMaps created successful { Namespace : %s, name : %s }", ns, name)
	return nil
}

// 格式化bufferpool配置
func FormatBufferpool(m string) string {
	if find := strings.Contains(m, "M"); find {
		memory := strings.Split(m, "M")
		intmemory, _ := strconv.ParseFloat(memory[0], 32)
		intbuffermem := intmemory * 0.6
		buffermem, _ := strconv.Atoi(fmt.Sprintf("%1.0f", intbuffermem))
		strmem := fmt.Sprintf("%dM", buffermem)
		return strmem

	} else if find := strings.Contains(m, "G"); find {
		memory := strings.Split(m, "G")
		intmemory, _ := strconv.ParseFloat(memory[0], 32)
		intbuffermem := intmemory * 0.6
		buffermem, _ := strconv.Atoi(fmt.Sprintf("%1.0f", intbuffermem))
		strmem := fmt.Sprintf("%dG", buffermem)
		return strmem
	}
	return "1G"
}

// 格式化配置文件
func ReadMycnf(s string, id string, buffer string) (string, error) {
	cfg1 := strings.Replace(s, "MYSQL_SERVER_ID", id, -1)
	cfg2 := strings.Replace(cfg1, "MYSQL_BUFFER_POOL_SIZE", buffer, -1)
	return cfg2, nil
}
