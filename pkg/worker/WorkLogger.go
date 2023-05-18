package worker

import (
	"context"
)

type (
	WorkLogger func(
		ctx context.Context,
		worker int,
		workCounter int,
		results WorkResults,
	)
)

func emptyWorkLogger(
	ctx context.Context,
	worker int,
	workCounter int,
	results WorkResults,
) {
	// Nothing to do because is the empty logger
}
