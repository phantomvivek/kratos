package main

import (
	"github.com/kratos/config"
	"github.com/kratos/service"
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
