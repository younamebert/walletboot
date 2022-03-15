package bootcron

import (
	"time"
	"walletboot/app"
	"walletboot/config"

	"github.com/sirupsen/logrus"
)

type Cron struct {
	app  *app.App
	spec string
	quit chan struct{}
}

func (job *Cron) createAcount() {
	for i := 0; i < config.NewAccountNumber; i++ {
		if err := job.app.CreateAccount(); err != nil {
			// fmt.Println(err)
			logrus.Warn(err)
			// job.Stop()
			// return
			continue
		}
	}
}

func (job *Cron) transfer() {
	for i := 0; i < config.SendTxNumber; i++ {
		if err := job.app.SendTransaction(); err != nil {
			// logrus.Error(err)
			// job.Stop()
			// return
			logrus.Warn(err)
			// job.Stop()
			continue
		}
	}
}

func (job *Cron) RunCarry() {
	go job.transfer()
	go job.createAcount()
}

func New() (*Cron, error) {
	appcore := app.New()
	// if err != nil {
	// 	return nil, err
	// }
	return &Cron{
		app:  appcore,
		spec: config.CronSpec,
		quit: make(chan struct{}),
	}, nil
}

func (c *Cron) Stop() {
	// close(c.quit)
	// c.quit = make(chan struct{})
}

// func (c *Cron) AppCore() *app.App {
// 	return c.app
// }

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
			c.RunCarry()
			go c.app.UpdateAccountState()
		case <-c.quit:
			logrus.Info("Stop for walletboot")
			c.quit = make(chan struct{})
			break out
		}
	}
}
