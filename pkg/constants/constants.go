package constants


const (
	MYSQL_ROOT_PASSWORD = "root1234"
	// 镜像仓库使用"/"结尾
	Registry_Addr = "bitnami/"
	Image = "mysql:5.7.39"
)

// 定义MySQL规格,单位分别为core,GB,GB
const (
	Small_CPU = 1
	Small_Memory = 1
	Small_Disk = 10
	Medium_CPU = 2
	Medium_Memory = 2
	Medium_Disk = 50
	Large_CPU = 4
	Large_Memory = 8
	Large_Disk = 100
)
