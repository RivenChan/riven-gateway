package main

import (
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"riven-gateway/config"
	"riven-gateway/proxy"
	"riven-gateway/router"
)

func main() {

	c, err := os.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	var bc config.Bootstrap
	err = yaml.Unmarshal(c, &bc)
	if err != nil {
		panic(err)
	}
	server := proxy.NewProxy(nil, &bc)
	http.HandleFunc("/", router.RequestAllocate)
	err = server.Start(nil)
	if err != nil {
		panic(err)
	}
}
