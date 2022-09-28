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
func (r *MysqlReconciler) QueryMysqlCM(ns string, name string, ctx context.Context) error {
	foundCM := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, foundCM)
	if err == nil {
		msg := fmt.Sprintf("MySQL Configmaps exists: { namespace: %s, name : %s }", ns, name)
		errmsg := errors.New(msg)
		logrus.Error(errmsg)
		return errmsg
	}
	logrus.Warnf("MySQL Configmaps not found: { namespace: %s, name : %s }", ns, name)
	return err
}


// 创建ConfigMaps,如果不存在则返回错误
func (r *MysqlReconciler) CreateRepMysqlCM(m *mysqlv1.Mysql, ns string, name string, role string, instance string, ctx context.Context) error {
	logrus.Infof("MySQL ConfigMaps creating: { namespace:'%s', name:'%s' }", ns, name)
	var server_id string
	if find := strings.Contains(name, "master"); find {
		server_id = "10"
		config_cm, _ := ReadMycnf(constants.MySQLCfg, server_id, FormatBufferpool(constants.InstanceReflect[instance]["Memory"]))
		init_cm := constants.InitCfg
		optionCM := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: ns,
			},
			Data: map[string]string{
				"my.cnf":  config_cm,
				"init.sh": init_cm,
			},
		}
		// 设置Configmaps的上级控制器
		err := controllerutil.SetControllerReference(m, optionCM, r.Scheme)
		if err != nil {
			logrus.Errorf("MySQL Configmaps set controller failed { namespace:'%s', name:'%s' }", ns, name)
		}

		err = r.Create(context.TODO(), optionCM)
		if err != nil {
			logrus.Error(err)
			return err
		}
	} else {
		server_id = "11"
		config_cm, _ := ReadMycnf(constants.MySQLSlaveCfg, server_id, FormatBufferpool(constants.InstanceReflect[instance]["Memory"]))
		init_cm := constants.InitCfg
		optionCM := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: ns,
			},
			Data: map[string]string{
				"my.cnf":  config_cm,
				"init.sh": init_cm,
			},
		}
		// 设置Configmaps的上级控制器
		err := controllerutil.SetControllerReference(m, optionCM, r.Scheme)
		if err != nil {
			logrus.Errorf("MySQL Configmaps set controller failed { namespace:'%s', name:'%s' }", ns, name)
		}

		err = r.Create(context.TODO(), optionCM)
		if err != nil {
			logrus.Error(err)
			return err
		}
	}
	logrus.Infof("MySQL ConfigMaps created successful { Namespace : %s, name : %s }", ns, name)
	return nil
}

// 创建ConfigMaps,如果不存在则返回错误
func (r *MysqlReconciler) CreateSingleMysqlCM(m *mysqlv1.Mysql, ns string, name string, role string, instance string, ctx context.Context) error {
	logrus.Infof("MySQL ConfigMaps creating: { namespace:'%s', name:'%s' }", ns, name)
	server_id := "11"
	config_cm, _ := ReadMycnf(constants.MySQLCfg, server_id, FormatBufferpool(constants.InstanceReflect[instance]["Memory"]))
	init_cm := constants.InitCfg
	mysql_job :=  constants.MySQLSingleJob
	optionCM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Data: map[string]string{
			"my.cnf":  config_cm,
			"init.sh": init_cm,
			"mysqljob.sh":mysql_job,
		},
	}
	// 设置Configmaps的上级控制器
	err := controllerutil.SetControllerReference(m, optionCM, r.Scheme)
	if err != nil {
		logrus.Errorf("MySQL Configmaps set controller failed { namespace:'%s', name:'%s' }", ns, name)
	}

	err = r.Create(context.TODO(), optionCM)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("MySQL ConfigMaps created successful { Namespace : %s, name : %s }", ns, name)
	return nil
}



// 查询ProxySQL ConfigMaps,如果不存在则返回错误
func (r *MysqlReconciler) QueryProxyCM(ns string, name string, ctx context.Context) error {
	foundCM := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, foundCM)
	if err == nil {
		msg := fmt.Sprintf("ProxySQL Configmaps exists: { namespace: %s, name : %s }", ns, name)
		errmsg := errors.New(msg)
		logrus.Error(errmsg)
		return errmsg
	}
	logrus.Warnf("ProxySQL Configmaps not found: { namespace: %s, name : %s }", ns, name)
	return err
}

// 创建Proxy ConfigMaps,如果不存在则返回错误
func (r *MysqlReconciler) CreateProxyCM(m *mysqlv1.Mysql, ns string, name string, databases []map[string]string, ctx context.Context) error {
	logrus.Infof("Proxy ConfigMaps creating: { namespace:'%s', name:'%s' }", ns, name)
	cm, _ := ReadProxycnf(constants.ProxyCfg, ns, name, "mydb", "mydbpass", databases)
	mysql_job := constants.MySQLJob
	optionCM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name + "-proxy",
			Namespace: ns,
		},
		Data: map[string]string{
			"proxysql.cnf": cm,
			"mysqljob.sh":  mysql_job,
		},
	}
	// 设置Configmaps的上级控制器
	err := controllerutil.SetControllerReference(m, optionCM, r.Scheme)
	if err != nil {
		logrus.Errorf("Proxy Configmaps set controller failed { namespace:'%s', name:'%s' }", ns, name)
	}

	err = r.Create(context.TODO(), optionCM)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("Proxy ConfigMaps created successful { Namespace : %s, name : %s }", ns, name)
	return nil
}

// 查询init ConfigMaps,如果不存在则返回错误
func (r *MysqlReconciler) QueryInitCM(ns string, name string, ctx context.Context) error {
	foundCM := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, foundCM)
	if err == nil {
		msg := fmt.Sprintf("Init Configmaps exists: { namespace: %s, name : %s }", ns, name)
		errmsg := errors.New(msg)
		logrus.Error(errmsg)
		return errmsg
	}
	logrus.Warnf("Init Configmaps not found: { namespace: %s, name : %s }", ns, name)
	return err
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

// 格式化Mysql配置文件
func ReadMycnf(s string, id string, buffer string) (string, error) {
	cfg1 := strings.Replace(s, "MYSQL_SERVER_ID", id, -1)
	cfg2 := strings.Replace(cfg1, "MYSQL_BUFFER_POOL_SIZE", buffer, -1)
	return cfg2, nil
}

// 格式化ProxySQL配置文件
func ReadProxycnf(s string, ns string, name string, user string, password string, databases []map[string]string) (string, error) {
	cfg1 := strings.Replace(s, "NAMESPACE", ns, -1)
	cfg2 := strings.Replace(cfg1, "MYSQL-ADDR", name, -1)
	cfg3 := strings.Replace(cfg2, "MONITOR-USER", constants.Monitor_User, -1)
	cfg4 := strings.Replace(cfg3, "MONITOR-PASSWORD", constants.Monitor_Password, -1)
	cfg5 := FormatMysqlusers(databases)
	cfg6 := strings.Replace(cfg4, "MYSQL_USERS", cfg5, -1)
	return cfg6, nil
}

// 格式化proxysql中MYSQL_SERVERS配置
func FormatMysqlusers(databases []map[string]string) string {
	dbs := ""
	for _, db := range databases {
		cfg := `
    {
        username = "MYSQL-USERNAME"
        password = "MYSQL-PASSWORD"
        active = 1
        max_connections = 1000
        default_hostgroup = 1
    },
`
		cfg1 := strings.Replace(cfg, "MYSQL-USERNAME", db["user"], -1)
		cfg2 := strings.Replace(cfg1, "MYSQL-PASSWORD", db["passwd"], -1)
		dbs = dbs + cfg2
	}
	return dbs
}
