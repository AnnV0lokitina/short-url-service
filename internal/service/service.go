package service

import (
	"context"
	"fmt"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"golang.org/x/sync/errgroup"
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
	s.deleteChanInput = make(chan *JobDelete, len(checksums))
	go func() {
		fmt.Println("start wait")
		s.processDeleteRequests(nOfWorkers)
		fmt.Println("end wait")
	}()
	fmt.Println("start send")
	for _, checksum := range checksums {
		fmt.Println(checksum)
		s.deleteChanInput <- NewJobDelete(userID, checksum, s.baseURL)
	}
	close(s.deleteChanInput)
	fmt.Println("end send")
}

func (s *Service) processDeleteRequests(workersCount int) {
	fmt.Println("fanOutDeleteURL create")
	fanOutChs := fanOutDeleteURL(s.deleteChanInput, workersCount)

	fmt.Println("workerChs create")
	workerChs := make([]chan *entity.UserShortURL, 0, workersCount)
	for _, fanOutCh := range fanOutChs {
		w := newWorkerDeleteURL(fanOutCh)
		workerChs = append(workerChs, w)
	}
	fmt.Println("workerChs created!")

	chanOut := fanInDeleteURL(workerChs...)
	fmt.Println("fanInDeleteURL created!")

	ctx, cancel := context.WithTimeout(context.Background(), writeToDBDuration)
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.writeToDB(ctx, chanOut)
	})
	g.Wait()
	cancel()

}

func (s *Service) writeToDB(ctx context.Context, chanOut <-chan *entity.UserShortURL) error {
	batch := make([]*entity.UserShortURL, 0, batchSize)
	fmt.Println("start writeToDb")
	for urlInfo := range chanOut {
		fmt.Println("read from out ch: writeToDb i")
		batch = append(batch, urlInfo)
		if len(batch) < cap(batch) {
			continue
		}
		err := s.repo.DeleteBatch(ctx, batch)
		if err != nil {
			fmt.Println("db exception")
			return err
		}
		batch = batch[:0]
	}
	fmt.Print("batch")
	fmt.Println(batch)
	if len(batch) > 0 {
		fmt.Println("writeToDb i after")
		err := s.repo.DeleteBatch(ctx, batch)
		if err != nil {
			return err
		}
	}
	fmt.Println("end!!!!!!!!!!!!!!!!!!!!!!")
	return nil
}

func newWorkerDeleteURL(inputCh <-chan *JobDelete) chan *entity.UserShortURL {
	outCh := make(chan *entity.UserShortURL)
	fmt.Println("make worker")

	go func() {
		for job := range inputCh {
			// параллельная обработка входящих данных
			fmt.Print("process worker ")
			fmt.Println(job.Checksum)
			outCh <- entity.NewUserShortURL(job.UserID, job.Checksum, job.BaseURL)
		}

		close(outCh)
	}()

	return outCh
}

func fanOutDeleteURL(inputCh chan *JobDelete, n int) []chan *JobDelete {
	chs := make([]chan *JobDelete, 0, n)
	for i := 0; i < n; i++ {
		fmt.Println("create i ch fanOut")
		ch := make(chan *JobDelete)
		chs = append(chs, ch)
	}

	go func() {
		defer func(chs []chan *JobDelete) {
			for _, ch := range chs {
				fmt.Println("close i ch fanOut")
				close(ch)
			}
		}(chs)

		for i := 0; ; i++ {
			if i == len(chs) {
				i = 0
			}
			job, ok := <-inputCh
			if !ok {
				fmt.Println("quit close inputCh fanOut")
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
	fmt.Println("made fan in chanel")

	go func() {
		wg := &sync.WaitGroup{}

		for _, inputCh := range inputChs {
			wg.Add(1)
			fmt.Println("start wait fanInDeleteURL 1")

			go func(inputCh chan *entity.UserShortURL) {
				fmt.Println("do fanInDeleteURL")
				defer func() {
					wg.Done()
				}()
				for item := range inputCh {
					fmt.Print("outCh <- item fanInDeleteURL")
					fmt.Println(item)
					outCh <- item
				}
				fmt.Println("end tread")
			}(inputCh)
		}
		fmt.Println("wait fanInDeleteURL 1")
		wg.Wait()
		fmt.Println("close outCh fanInDeleteURL")
		close(outCh)
	}()

	return outCh
}
