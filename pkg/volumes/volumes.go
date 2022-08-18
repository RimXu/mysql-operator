package volumes

import (
	"fmt"

)

func CreateVolumes(size string,sc string) error {
	if sc == "" {
		sc = "local"
	}
	fmt.Println("CreateVolumes",size,sc)
	return nil
}
