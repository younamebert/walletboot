package bootcron

import (
	"fmt"
	"time"
	"walletboot/appcore"
	"walletboot/config"

	"github.com/sirupsen/logrus"
)

type Cron struct {
	appcore *appcore.AppCore
	spec    string
	quit    chan struct{}
}

func (job *Cron) CronBatchRunRand() {
	for i := 0; i < config.NewAccountNumber; i++ {
		if err := job.appcore.RunRand(); err != nil {
			fmt.Println(err)
			// job.Stop()
			// return
			continue
		}
	}
}

func (job *Cron) CronBatchRunSendTx() {
	for i := 0; i < config.SendTxNumber; i++ {
		if err := job.appcore.RunSendTx(); err != nil {
			// logrus.Error(err)
			// job.Stop()
			// return
			fmt.Println(err)
			continue
		}
	}
}

func New() (*Cron, error) {
	appcore, err := appcore.New()
	if err != nil {
		return nil, err
	}
	return &Cron{
		appcore: appcore,
		spec:    config.CronSpec,
		quit:    make(chan struct{}),
	}, nil
}

func (c *Cron) Stop() {
	close(c.quit)
	// c.quit = make(chan struct{})
}

func (c *Cron) AppCore() *appcore.AppCore {
	return c.appcore
}

func (c *Cron) Start() {

	timeDur, err := time.ParseDuration(config.CronSpec)
	if err != nil {
		logrus.Error(err)
		return
	}
out:
	for {
		select {
		case <-time.After(timeDur):
			c.CronBatchRunRand()
			c.CronBatchRunSendTx()
		case <-c.quit:
			logrus.Info("Stop for walletboot")
			c.quit = make(chan struct{})
			break out
		}
	}
}
