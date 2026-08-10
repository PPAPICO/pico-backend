package cron

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron      *cron.Cron
	policyJob *PolicyJob
}

func NewScheduler(policyJob *PolicyJob) *Scheduler {
	loc, _ := time.LoadLocation("Asia/Seoul")

	c := cron.New(
		cron.WithLocation(loc),
	)

	_, err := c.AddFunc("0 9,12,15,18 * * *", policyJob.Run)
	if err != nil {
		panic(err)
	}

	return &Scheduler{
		cron:      c,
		policyJob: policyJob,
	}
}

func (s *Scheduler) Start() {
	go s.policyJob.Run()
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
