package service

import (
	"context"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"sync"
	"time"
)

const batchSize = 2
const writeToDBDuration = 5 * time.Minute
const nOfWorkers = 2

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
	baseURL string
	repo    Repo
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
	deleteChanInput := make(chan *JobDelete, len(checksums))
	go func() {
		s.processDeleteRequests(deleteChanInput, nOfWorkers)
	}()
	for _, checksum := range checksums {
		deleteChanInput <- NewJobDelete(userID, checksum, s.baseURL)
	}
	close(deleteChanInput)
}

func (s *Service) processDeleteRequests(deleteChanInput chan *JobDelete, workersCount int) {
	fanOutChs := fanOutDeleteURL(deleteChanInput, workersCount)
	workerChs := make([]chan *entity.UserShortURL, 0, workersCount)
	for _, fanOutCh := range fanOutChs {
		w := newWorkerDeleteURL(fanOutCh)
		workerChs = append(workerChs, w)
	}

	chanOut := fanInDeleteURL(workerChs...)

	ctx, cancel := context.WithTimeout(context.Background(), writeToDBDuration)
	s.writeToDB(ctx, chanOut)
	cancel()
}

func (s *Service) writeToDB(ctx context.Context, chanOut <-chan *entity.UserShortURL) error {
	batch := make([]*entity.UserShortURL, 0, batchSize)
	for urlInfo := range chanOut {
		batch = append(batch, urlInfo)
		if len(batch) < cap(batch) {
			continue
		}
		err := s.repo.DeleteBatch(ctx, batch)
		if err != nil {
			return err
		}
		batch = batch[:0]
	}
	if len(batch) > 0 {
		err := s.repo.DeleteBatch(ctx, batch)
		if err != nil {
			return err
		}
	}
	return nil
}

func newWorkerDeleteURL(inputCh <-chan *JobDelete) chan *entity.UserShortURL {
	outCh := make(chan *entity.UserShortURL)

	go func() {
		for job := range inputCh {
			// параллельная обработка входящих данных
			outCh <- entity.NewUserShortURL(job.UserID, job.Checksum, job.BaseURL)
		}

		close(outCh)
	}()

	return outCh
}

func fanOutDeleteURL(inputCh chan *JobDelete, n int) []chan *JobDelete {
	chs := make([]chan *JobDelete, 0, n)
	for i := 0; i < n; i++ {
		ch := make(chan *JobDelete)
		chs = append(chs, ch)
	}

	go func() {
		defer func(chs []chan *JobDelete) {
			for _, ch := range chs {
				close(ch)
			}
		}(chs)

		for i := 0; ; i++ {
			if i == len(chs) {
				i = 0
			}
			job, ok := <-inputCh
			if !ok {
				return
			}

			ch := chs[i]
			ch <- job
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
				defer func() {
					wg.Done()
				}()
				for item := range inputCh {
					outCh <- item
				}
			}(inputCh)
		}
		wg.Wait()
		close(outCh)
	}()
	return outCh
}
