package main

import (
	"silvernote/factory"
	"silvernote/service"
)

func main() {

	var fac factory.Factory

	fac.Initialize()

	restServer := service.HttpServer{Fac: &fac}

	restServer.StartService()
}
