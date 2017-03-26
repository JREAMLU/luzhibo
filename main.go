//go:generate goversioninfo -icon=icon.ico

package main

import (
	"fmt"

)

const ver  =2017032600
const p  = "录直播"

func main() {
	s := ":12216"
	fmt.Printf("正在\"%s\"处监听WebUI...\n", s)
	go startServer(s)
	openWebUI()
	cmd()
}
