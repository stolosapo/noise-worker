package example

import (
	"context"
	"fmt"

	"github.com/stolosapo/noise-worker/pkg/worker"
)

type (
	UseCase interface {
		DoTheJob(ctx context.Context) error
	}

	useCase struct {
		repository Repository
		service    DomainService

		workerBufferSize int
		workConcurrency  int
	}
)

func NewUseCase(
	repository Repository,
	service DomainService,
	workerBufferSize int,
	workConcurrency int,
) *useCase {
	return &useCase{
		repository:       repository,
		service:          service,
		workerBufferSize: workerBufferSize,
		workConcurrency:  workConcurrency,
	}
}

func (u useCase) DoTheJob(ctx context.Context) error {
	w := worker.NewWorkWithLogs(
		u.workerBufferSize,
		u.workConcurrency,
		u.repository.FetchData,
		u.service.AVeryImportantJob,
		logWorkerStart,
		logWorkerEnd,
		logWorkFinished,
	)

	results := w.Start(ctx)
	if results.HasError() {
		return results.AllErrors()[0]
	}

	fmt.Printf(`
*************************
*   The Job is finished
* -----------------------
* FetchedWorkCount:    %d
* TotalWorkCount:      %d
* SuccessfulWorkCount: %d
*************************
	`,
		results.FetchedWorkCount(),
		results.TotalWorkCount(),
		results.SuccessfulWorkCount(),
	)

	return nil
}
