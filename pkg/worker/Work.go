package worker

import (
	"context"
	"sync"
)

type (
	Work[T any] interface {
		Start(ctx context.Context) WorkResults
	}

	work[T any] struct {
		bufferSize      int
		numberOfWorkers int

		fetchWorkDelegate FetchWorkDelegate[T]
		workDelegate      WorkDelegate[T]

		logWorkerStart  WorkLogger
		logWorkerEnd    WorkLogger
		logWorkFinished WorkLogger
	}
)

func NewWork[T any](
	bufferSize int,
	numberOfWorkers int,
	fetchWorkDelegate FetchWorkDelegate[T],
	workDelegate WorkDelegate[T],
) Work[T] {
	return NewWorkWithLogs(
		bufferSize,
		numberOfWorkers,
		fetchWorkDelegate,
		workDelegate,
		emptyWorkLogger,
		emptyWorkLogger,
		emptyWorkLogger,
	)
}

func NewWorkWithLogs[T any](
	bufferSize int,
	numberOfWorkers int,
	fetchWorkDelegate FetchWorkDelegate[T],
	workDelegate WorkDelegate[T],
	logWorkerStart WorkLogger,
	logWorkerEnd WorkLogger,
	logWorkFinished WorkLogger,
) Work[T] {
	return &work[T]{
		bufferSize:        bufferSize,
		numberOfWorkers:   numberOfWorkers,
		fetchWorkDelegate: fetchWorkDelegate,
		workDelegate:      workDelegate,
		logWorkerStart:    logWorkerStart,
		logWorkerEnd:      logWorkerEnd,
		logWorkFinished:   logWorkFinished,
	}
}

func (w *work[T]) Start(
	ctx context.Context,
) WorkResults {
	workChannel := make(chan T, w.bufferSize)
	var workerWG sync.WaitGroup
	results := newWorkResults()

	// Spawn Workers to do the Job
	for i := 0; i < w.numberOfWorkers; i++ {
		workerWG.Add(1)
		go w.startWorker(
			ctx,
			i,
			&workerWG,
			workChannel,
			results,
		)
	}

	// Fetch items into channel
	err := w.fetchWorkDelegate(
		ctx,
		&results.fetchedWorkCount,
		workChannel,
	)
	if err != nil {
		results.appendError(err)
	}

	// Close the channel so to stop the async functions
	close(workChannel)

	// Wait until all the remaining processes finished by the async functions
	workerWG.Wait()

	return results
}

func (w *work[T]) startWorker(
	ctx context.Context,
	worker int,
	workerWG *sync.WaitGroup,
	workChannel chan T,
	results *workResults,
) {
	defer workerWG.Done()

	workCnt := 0

	w.logWorkerStart(
		ctx,
		worker,
		workCnt,
		results,
	)

	defer func() {
		w.logWorkerEnd(
			ctx,
			worker,
			workCnt,
			results,
		)
	}()

	for {
		select {
		case work, open := <-workChannel:
			if !open {
				return
			}

			workCnt++

			err := w.workDelegate(ctx, work)
			if err != nil {
				results.appendError(err)
			}

			w.logWorkFinished(
				ctx,
				worker,
				workCnt,
				results,
			)

			if err == nil {
				results.incrementSuccessfulWorkCount()
			}
			results.incrementTotalWorkCount()
		case <-ctx.Done():
			return
		}
	}
}
