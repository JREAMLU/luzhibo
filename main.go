package main

import (
	"fmt"
	"github.com/pkg/browser"

)

func main() {
	s := ":12216"
	fmt.Printf("正在\"%s\"处监听WebUI...\n", s)
	go startServer(s)
	browser.OpenURL("http://localhost:12216")
	cmd()
}
