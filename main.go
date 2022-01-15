package main

import (
	"walletboot/bootcron"

	"github.com/sirupsen/logrus"
)

func main() {
	serve, err := bootcron.New()
	if err != nil {
		logrus.Error(err)
		return
	}
	go serve.Start()
	select {}
}
