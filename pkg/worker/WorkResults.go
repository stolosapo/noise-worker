package worker

import "sync"

type (
	WorkResults interface {
		FetchedWorkCount() int
		TotalWorkCount() int
		SuccessfulWorkCount() int
		AllErrors() []error
		HasError() bool
	}

	workResults struct {
		resultsLocker       sync.RWMutex
		fetchedWorkCount    int
		totalWorkCount      int
		successfulWorkCount int
		allErrors           []error
	}
)

func newWorkResults() *workResults {
	return &workResults{
		allErrors: []error{},
	}
}

func (r *workResults) FetchedWorkCount() int {
	r.resultsLocker.RLock()
	defer r.resultsLocker.RUnlock()
	return r.fetchedWorkCount
}

func (r *workResults) incrementTotalWorkCount() {
	r.resultsLocker.RLock()
	defer r.resultsLocker.RUnlock()
	r.totalWorkCount++
}

func (r *workResults) TotalWorkCount() int {
	r.resultsLocker.RLock()
	defer r.resultsLocker.RUnlock()
	return r.totalWorkCount
}

func (r *workResults) incrementSuccessfulWorkCount() {
	r.resultsLocker.Lock()
	defer r.resultsLocker.Unlock()
	r.successfulWorkCount++
}

func (r *workResults) SuccessfulWorkCount() int {
	r.resultsLocker.RLock()
	defer r.resultsLocker.RUnlock()
	return r.successfulWorkCount
}

func (r *workResults) appendError(err ...error) {
	r.resultsLocker.Lock()
	defer r.resultsLocker.Unlock()
	r.allErrors = append(r.allErrors, err...)
}

func (r *workResults) AllErrors() []error {
	r.resultsLocker.RLock()
	defer r.resultsLocker.RUnlock()
	return r.allErrors
}

func (r *workResults) HasError() bool {
	return len(r.AllErrors()) > 0
}
