package main

import (
	"flag"
	"log"

	"github.com/docker/go-plugins-helpers/authorization"
)

const (
	pluginSocket = "/run/docker/plugins/nomonkey.sock"
)

func main() {
	flag.Parse()

	nomonkey, err := newPlugin()
	if err != nil {
		log.Fatal(err)
	}

	h := authorization.NewHandler(nomonkey)
	log.Println("nomonkey plugin ready")
	if err := h.ServeUnix("root", pluginSocket); err != nil {
		log.Fatal(err)
	}
}
