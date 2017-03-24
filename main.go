package main

import (
	"fmt"

)

const ver  =2017032401

func main() {
	s := ":12216"
	fmt.Printf("正在\"%s\"处监听WebUI...\n", s)
	go startServer(s)
	openWebUI()
	cmd()
}
