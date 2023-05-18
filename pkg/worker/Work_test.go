package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_work_Start_WithPointers(t *testing.T) {
	type workItem struct {
		id int
	}
	type fields struct {
		fetchWorkDelegate FetchWorkDelegate[*workItem]
		workDelegate      WorkDelegate[*workItem]
	}
	tests := []struct {
		name    string
		fields  fields
		want    *workResults
		wantErr bool
	}{
		{
			name: "Should fail when fetch failed",
			fields: fields{
				fetchWorkDelegate: func(ctx context.Context, fetchedWorkCounter *int, workChannel chan *workItem) error {
					return errors.New("an error")
				},
				workDelegate: func(ctx context.Context, workItem *workItem) error {
					return nil
				},
			},
			want:    &workResults{},
			wantErr: true,
		},

		{
			name: "Should return correct values",
			fields: fields{
				fetchWorkDelegate: func(ctx context.Context, fetchedWorkCounter *int, workChannel chan *workItem) error {
					for i := 0; i < 10; i++ {
						*fetchedWorkCounter++
						workChannel <- &workItem{id: i}
					}
					return nil
				},
				workDelegate: func(ctx context.Context, workItem *workItem) error {
					fmt.Printf("Item: %d, Time: %v\n", workItem.id, time.Now().UnixMilli())
					return nil
				},
			},
			want: &workResults{
				fetchedWorkCount:    10,
				totalWorkCount:      10,
				successfulWorkCount: 10,
			},
			wantErr: false,
		},

		{
			name: "Should return correct values with some failed",
			fields: fields{
				fetchWorkDelegate: func(ctx context.Context, fetchedWorkCounter *int, workChannel chan *workItem) error {
					for i := 0; i < 10; i++ {
						*fetchedWorkCounter++
						workChannel <- &workItem{id: i}
					}
					return nil
				},
				workDelegate: func(ctx context.Context, workItem *workItem) error {
					if workItem.id == 3 {
						return errors.New("an error")
					}
					fmt.Printf("Item: %d, Time: %v\n", workItem.id, time.Now().UnixMilli())
					return nil
				},
			},
			want: &workResults{
				fetchedWorkCount:    10,
				totalWorkCount:      10,
				successfulWorkCount: 9,
			},
			wantErr: true,
		},

		{
			name: "Should return correct values when no job to do",
			fields: fields{
				fetchWorkDelegate: func(ctx context.Context, fetchedWorkCounter *int, workChannel chan *workItem) error {
					return nil
				},
				workDelegate: func(ctx context.Context, workItem *workItem) error {
					return nil
				},
			},
			want: &workResults{
				fetchedWorkCount:    0,
				totalWorkCount:      0,
				successfulWorkCount: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			w := NewWorkWithLogs(
				1,
				10,
				tt.fields.fetchWorkDelegate,
				tt.fields.workDelegate,
				emptyWorkLogger,
				emptyWorkLogger,
				emptyWorkLogger,
			)
			got := w.Start(ctx)

			if got.HasError() != tt.wantErr {
				t.Errorf("work.Start() error = %v, wantErr %v", got.HasError(), tt.wantErr)
			}

			if got.FetchedWorkCount() != tt.want.fetchedWorkCount {
				t.Errorf("work.Start() fetchedWorkCount = %v, want %v", got.FetchedWorkCount(), tt.want.fetchedWorkCount)
			}

			if got.TotalWorkCount() != tt.want.totalWorkCount {
				t.Errorf("work.Start() totalWorkCount = %v, want %v", got.TotalWorkCount(), tt.want.totalWorkCount)
			}

			if got.SuccessfulWorkCount() != tt.want.successfulWorkCount {
				t.Errorf("work.Start() successfulWorkCount = %v, want %v", got.SuccessfulWorkCount(), tt.want.successfulWorkCount)
			}
		})
	}
}

func Test_work_Start_WithContextDone(t *testing.T) {
	type workItem struct {
		id int
	}
	type fields struct {
		fetchWorkDelegate FetchWorkDelegate[*workItem]
		workDelegate      WorkDelegate[*workItem]
		timeUntilCancel   time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		want    *workResults
		wantErr bool
	}{
		{
			name: "Should process less than fetched",
			fields: fields{
				fetchWorkDelegate: func(ctx context.Context, fetchedWorkCounter *int, workChannel chan *workItem) error {
					for i := 0; i < 10; i++ {
						*fetchedWorkCounter++
						workChannel <- &workItem{id: i}
					}
					return nil
				},
				workDelegate: func(ctx context.Context, workItem *workItem) error {
					time.Sleep(500 * time.Millisecond)
					fmt.Printf("Item: %d, Time: %v\n", workItem.id, time.Now().UnixMilli())
					return nil
				},
				timeUntilCancel: 100 * time.Millisecond,
			},
			want: &workResults{
				fetchedWorkCount: 10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			w := NewWork(
				10,
				1,
				tt.fields.fetchWorkDelegate,
				tt.fields.workDelegate,
			)
			var got WorkResults
			var wg sync.WaitGroup

			wg.Add(1)
			go func(ctx context.Context, wg *sync.WaitGroup) {
				defer wg.Done()
				got = w.Start(ctx)
			}(ctx, &wg)

			time.Sleep(tt.fields.timeUntilCancel)
			cancel()
			wg.Wait()

			if got.HasError() != tt.wantErr {
				t.Errorf("work.Start() error = %v, wantErr %v", got.HasError(), tt.wantErr)
			}

			if got.FetchedWorkCount() != tt.want.fetchedWorkCount {
				t.Errorf("work.Start() fetchedWorkCount = %v, want %v", got.FetchedWorkCount(), tt.want.fetchedWorkCount)
			}

			if got.TotalWorkCount() >= got.FetchedWorkCount() {
				t.Errorf("work.Start() totalWorkCount = %v, want less than %v", got.TotalWorkCount(), got.FetchedWorkCount())
			}
		})
	}
}
