package main

import (
	"fmt"
	"strings"
	"strconv"
)

func FormatBufferpool(m string) (string,error) {
	if find := strings.Contains(m, "M"); find {
		memory := strings.Split(m,"M")
	        fmt.Println(memory[0],memory[1])
        	intmemory, _ := strconv.ParseFloat(memory[0],32)
        	intbuffermem := intmemory * 0.6
        	buffermem, _ := strconv.Atoi(fmt.Sprintf("%1.0f", intbuffermem))
        	strmem := fmt.Sprintf("%dM",buffermem)
		fmt.Println(strmem)
        	return strmem,nil

	} else if find := strings.Contains(m, "G"); find {
		memory := strings.Split(m,"G")
                fmt.Println(memory[0],memory[1])
                intmemory, _ := strconv.ParseFloat(memory[0],32)
                intbuffermem := intmemory * 0.6
                buffermem, _ := strconv.Atoi(fmt.Sprintf("%1.0f", intbuffermem))
                strmem := fmt.Sprintf("%dG",buffermem)
                fmt.Println(strmem)
                return strmem,nil
	}
	return "",nil
}


func main() {
	FormatBufferpool("5Gi")
}
