package main

import (
	"testtx/testcron"

	"github.com/sirupsen/logrus"
)

func main() {
	serve, err := testcron.New()
	if err != nil {
		logrus.Error(err)
		return
	}
	go serve.Start()
	select {}
}
