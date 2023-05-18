package main

import (
	"context"

	"github.com/stolosapo/noise-worker/internal/example"
)

const (
	itemsToFetch     int = 1000
	workerBufferSize int = 100
	workConcurrency  int = 100
)

func main() {
	ctx := context.Background()

	repository := example.NewRepository(itemsToFetch)
	service := example.NewDomainService()
	useCase := example.NewUseCase(repository, service, workerBufferSize, workConcurrency)

	err := useCase.DoTheJob(ctx)
	if err != nil {
		panic(err)
	}
}
