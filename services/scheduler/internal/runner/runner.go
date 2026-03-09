package runner

import (
	"context"
	"log"

	"github.com/iamsainadh/neovault/scheduler/internal/jobs"
	"github.com/robfig/cron/v3"
)

type Runner struct {
	cron *cron.Cron
	job  *jobs.RotationJob
}

func NewRunner(job *jobs.RotationJob) *Runner {
	return &Runner{
		cron: cron.New(),
		job:  job,
	}
}

func (r *Runner) Register(spec string) {
	_, err := r.cron.AddFunc(spec, func() {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		r.job.Run(ctx)
	})

	if err != nil {
		log.Fatalf("failed to register cron job: %v", err)
	}
}

func (r *Runner) Start() {
	r.cron.Start()
}
