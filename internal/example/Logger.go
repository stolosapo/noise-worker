package example

import (
	"context"
	"fmt"

	"github.com/stolosapo/noise-worker/pkg/worker"
)

func logWorkerStart(
	ctx context.Context,
	worker int,
	workCounter int,
	results worker.WorkResults,
) {
	fmt.Printf("A new Worker (%d) started the accept jobs\n", worker)
}

func logWorkerEnd(
	ctx context.Context,
	worker int,
	workCounter int,
	results worker.WorkResults,
) {
	fmt.Printf("The Worker (%d) ended %d jobs\n", worker, workCounter)
}

func logWorkFinished(
	ctx context.Context,
	worker int,
	workCounter int,
	results worker.WorkResults,
) {
	fmt.Printf("The job (%d) ended from Worker: %d\n", workCounter, worker)
}
