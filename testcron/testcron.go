package testcron

import (
	"time"
	"walletboot/appcore"
	"walletboot/config"
	"walletboot/httpxfs"

	"github.com/sirupsen/logrus"
)

type Cron struct {
	httpClient *httpxfs.Client
	appcore    *appcore.AppCore
	spec       string
}

func (job *Cron) CronBatchRunRand() {
	for i := 0; i < config.NewAccountNumber; i++ {
		job.appcore.RunRand()
	}
}

func (job *Cron) CronBatchRunSendTx() {
	for i := 0; i < config.SendTxNumber; i++ {
		job.appcore.RunSendTx()
	}
}

func (job *Cron) Run() {
	job.appcore.Run()
	job.CronBatchRunRand()
	job.CronBatchRunSendTx()
}

func New() (*Cron, error) {
	appcore, err := appcore.New()
	if err != nil {
		return nil, err
	}
	return &Cron{
		appcore: appcore,
		spec:    config.CronSpec,
	}, nil
}

func (c *Cron) Start() {

	timeDur, err := time.ParseDuration(config.CronSpec)
	if err != nil {
		logrus.Error(err)
		return
	}
	for {
		select {
		case <-time.After(timeDur):
			c.Run()
		}
	}
}