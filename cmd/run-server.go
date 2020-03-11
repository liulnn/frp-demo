package main

import (
	"frp-demo/frp_server"
	"fmt"
)

func main(){

	server, err := frp_server.NewServer(nil)
	if err != nil {
		fmt.Print(err)
		return
	}
	err = server.Run()
	if err != nil {
		fmt.Println(err)
	}
}
