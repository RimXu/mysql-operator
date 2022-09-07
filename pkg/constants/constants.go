package constants

const (
	MYSQL_ROOT_PASSWORD = "root1234"
	// 镜像仓库使用"/"结尾
	Registry_Addr  = "10.198.140.35/kce/"
	Mysql_Image    = "mysql:5.7.35"
	Exporter_Image = "mysqld-exporter:v0.12.1"
	ProxySQL_Image = "proxysql:2.2.2"
)

// 定义MySQL规格和单位
var InstanceReflect = map[string]map[string]string{
	"small":  {"CPU": "100m", "Memory": "400Mi", "Disk": "1Gi"},
	"medium": {"CPU": "200m", "Memory": "600Mi", "Disk": "2Gi"},
	"large":  {"CPU": "400m", "Memory": "800Mi", "Disk": "3Gi"},
}

// 定义数据库用户配置
const (
	// 用于proxysql进行状态监控
	Monitor_User     = "monitor123"
	Monitor_Password = "monitorMWF"
	Monitor_reg      = "%"
	// 用于mysql主从同步
	Repl_User     = "repl"
	Repl_Password = "replMWF"
	Repl_Reg      = "%"
	// 用于mysql exporter指标查询获取
	Exporter_User     = "exporter"
	Exporter_Password = "exporterMWF"
	Exporter_Reg      = "127.0.0.1"
)

//
