package worker

import "context"

type (
	WorkDelegate[T any] func(
		ctx context.Context,
		workItem T,
	) error
)
