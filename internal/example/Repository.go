package example

import (
	"context"
	"fmt"
)

type (
	Repository interface {
		FetchData(
			ctx context.Context,
			fetchedDataCounter *int,
			dataChannel chan TheModel,
		) error
	}

	repository struct {
		howManyItemsToFetch int
	}
)

func NewRepository(
	howManyItemsToFetch int,
) *repository {
	return &repository{
		howManyItemsToFetch: howManyItemsToFetch,
	}
}

func (r repository) FetchData(
	ctx context.Context,
	fetchedDataCounter *int,
	dataChannel chan TheModel,
) error {
	if fetchedDataCounter == nil {
		fetchedDataCounter = new(int)
		*fetchedDataCounter = 0
	}

	for i := 0; i < r.howManyItemsToFetch; i++ {
		item := NewModel(
			i,
			fmt.Sprintf("Item: %d", i),
			fmt.Sprintf("This is the Item %d", i),
		)
		dataChannel <- item
		*fetchedDataCounter++
	}

	return nil
}
