package main

import (
	"fmt"

)

func main() {
	s := ":12216"
	fmt.Printf("正在\"%s\"处监听WebUI...\n", s)
	go startServer(s)
	openWebUI()
	cmd()
}
