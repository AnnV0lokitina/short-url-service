package service

import (
	"context"
	"fmt"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"log"
	"sync"
	"time"
)

const batchSize = 2

type Repo interface {
	SetURL(ctx context.Context, userID uint32, url *entity.URL) error
	GetURL(ctx context.Context, shortURL string) (*entity.URL, bool, error)
	GetUserURLList(ctx context.Context, id uint32) ([]*entity.URL, bool, error)
	PingBD(ctx context.Context) bool
	Close(context.Context) error
	AddBatch(ctx context.Context, userID uint32, list []*entity.BatchURLItem) error
	DeleteBatch(ctx context.Context, list []*entity.UserShortURL) error
}

type Service struct {
	baseURL         string
	repo            Repo
	deleteChanInput chan *JobDelete
}

type JobDelete struct {
	UserID   uint32
	Checksum string
	BaseURL  string
}

func NewJobDelete(userID uint32, checksum string, baseURL string) *JobDelete {
	return &JobDelete{
		UserID:   userID,
		Checksum: checksum,
		BaseURL:  baseURL,
	}
}

func NewService(baseURL string, repo Repo) *Service {
	return &Service{
		baseURL: baseURL,
		repo:    repo,
	}
}

func (s *Service) GetBaseURL() string {
	return s.baseURL
}

func (s *Service) SetBaseURL(baseURL string) {
	s.baseURL = baseURL
}

func (s *Service) GetRepo() Repo {
	return s.repo
}

func (s *Service) DeleteURLList(userID uint32, checksums []string) {
	for _, checksum := range checksums {
		fmt.Println(checksum)
		s.deleteChanInput <- NewJobDelete(userID, checksum, s.baseURL)
	}
}

func (s *Service) ProcessDeleteRequests(ctx context.Context, workersCount int) {
	duration := 5 * time.Minute
	ctx, _ = context.WithTimeout(context.Background(), duration)
	s.deleteChanInput = make(chan *JobDelete)
	fmt.Println("deleteChanInput create")
	fanOutChs := fanOutDeleteURL(ctx, s.deleteChanInput, workersCount)

	workerChs := make([]chan *entity.UserShortURL, 0, workersCount)
	for _, fanOutCh := range fanOutChs {
		w := newWorkerDeleteURL(ctx, fanOutCh)
		workerChs = append(workerChs, w)
	}

	chanOut := fanInDeleteURL(workerChs...)

	go func() {
		batch := make([]*entity.UserShortURL, 0, batchSize)
	loop:
		for {
			select {
			case job, ok := <-chanOut:
				if !ok {
					fmt.Println("quit close chanOut")
					return
				}
				batch = append(batch, job)
				if len(batch) < cap(batch) {
					continue
				}
				err := s.repo.DeleteBatch(ctx, batch)
				if err != nil {
					close(s.deleteChanInput)
					log.Fatal(err)
				}
				batch = batch[:0]
			case <-ctx.Done():
				fmt.Println("quit")
				break loop
			}
		}
		fmt.Println("deleteChanInput close")
		close(s.deleteChanInput)
	}()
	return
}

func newWorkerDeleteURL(ctx context.Context, inputCh <-chan *JobDelete) chan *entity.UserShortURL {
	outCh := make(chan *entity.UserShortURL)

	go func() {
	loop:
		for {
			select {
			case job, ok := <-inputCh:
				if !ok {
					fmt.Println("quit close i")
					break loop
				}
				// параллельная обработка входящих данных
				outCh <- entity.NewUserShortURL(job.UserID, job.Checksum, job.BaseURL)
			case <-ctx.Done():
				fmt.Println("quit")
				break loop
			}
		}
		close(outCh)
	}()

	return outCh
}

func fanOutDeleteURL(ctx context.Context, inputCh chan *JobDelete, n int) []chan *JobDelete {
	chs := make([]chan *JobDelete, 0, n)
	for i := 0; i < n; i++ {
		fmt.Println("create i ch")
		ch := make(chan *JobDelete)
		chs = append(chs, ch)
	}

	go func() {
		defer func(chs []chan *JobDelete) {
			for _, ch := range chs {
				fmt.Println("close i ch")
				close(ch)
			}
		}(chs)

		for i := 0; ; i++ {
			if i == len(chs) {
				i = 0
			}
			select {
			case job, ok := <-inputCh:
				if !ok {
					fmt.Println("quit close")
					return
				}
				ch := chs[i]
				fmt.Print("ch <- job ")
				fmt.Println(job)
				ch <- job
			case <-ctx.Done():
				fmt.Println("quit")
				return
			}
		}
	}()

	return chs
}

func fanInDeleteURL(inputChs ...chan *entity.UserShortURL) chan *entity.UserShortURL {
	outCh := make(chan *entity.UserShortURL)

	go func() {
		wg := &sync.WaitGroup{}

		for _, inputCh := range inputChs {
			wg.Add(1)

			go func(inputCh chan *entity.UserShortURL) {
				defer wg.Done()
				for item := range inputCh {
					fmt.Print("outCh <- item ")
					fmt.Println(item)
					outCh <- item
				}
			}(inputCh)
		}

		wg.Wait()
		fmt.Println("close outCh")
		close(outCh)
	}()

	return outCh
}
