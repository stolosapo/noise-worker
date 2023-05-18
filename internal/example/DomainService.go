package example

import (
	"context"
	"math/rand"
	"time"
)

type (
	DomainService interface {
		AVeryImportantJob(ctx context.Context, item TheModel) error
	}

	domainService struct {
		randomizer *rand.Rand
		max        int
	}
)

func NewDomainService() *domainService {
	return &domainService{
		randomizer: rand.New(rand.NewSource(time.Now().UnixNano())),
		max:        10,
	}
}

func (s domainService) AVeryImportantJob(ctx context.Context, item TheModel) error {
	randomSleep := s.randomizer.Intn(s.max) * 100 * int(time.Millisecond)

	time.Sleep(time.Duration(randomSleep))

	_ = item.WithValue(randomSleep)

	return nil
}
