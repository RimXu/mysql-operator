package constants

const (
	MYSQL_ROOT_PASSWORD = "root1234"
	// 镜像仓库使用"/"结尾
	Registry_Addr = "bitnami/"
	Image         = "mysql:5.7.39"
)

// 定义MySQL规格和单位
var ComboReflect = map[string] map[string] string {
	"small" :  {"CPU": "100m" , "Memory": "400Mi" ,"Disk": "10Gi"},
	"medium" : {"CPU": "200m" , "Memory": "400Mi", "Disk": "20Gi"},
	"large" :  {"CPU": "400m" , "Memory": "800Mi", "Disk": "30Gi"},
}
