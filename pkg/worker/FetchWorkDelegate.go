package worker

import "context"

type (
	FetchWorkDelegate[T any] func(
		ctx context.Context,
		fetchedWorkCounter *int,
		workChannel chan T,
	) error
)
