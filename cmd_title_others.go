//+build !windows

package main

import "fmt"

func setConsoleTitle(){
	fmt.Printf("\033]0;%s\007",title)
}