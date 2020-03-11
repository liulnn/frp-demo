package main

import (
	"frp-demo/frp_client"
	"fmt"
)

func main(){

	cli , err := frp_client.NewClient(nil)
	if err != nil {
		fmt.Print(err)
		return
	}
	err = cli.Run()
	if err != nil {
		fmt.Println(err)
	}
}
