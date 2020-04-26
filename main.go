package main

import (
	"github.com/phantomvivek/kratos/config"
	"github.com/phantomvivek/kratos/service"
)

func init() {

	//Set config in struct for access
	config.SetConfig()
}

func main() {

	service.TestRunner.Initialize()

	service.TestRunner.Start()

	//Wait for the tests to finish
	<-service.TestRunner.TestDoneChan
}
