package main


import (
	"fmt"
	"io/ioutil"
	"strings"
)




func ReadMycnf(MyKeys ...string) error {
	var k = make([]byte, 128)
	k, _ = ioutil.ReadFile("/code/go/operator/mysql-operator/config/samples/my.cnf")
	cfg1 := strings.Replace(string(k),"MYSQL_SERVER_ID","20",-1)
	cfg2 := strings.Replace(cfg1,"MYSQL_BUFFER_POOL_SIZE","200M",-1)
	fmt.Println(cfg2)
	return nil

}


func main() {
	ReadMycnf()

}
