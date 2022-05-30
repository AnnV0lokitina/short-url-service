package service

import (
	"context"
	"fmt"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"golang.org/x/sync/errgroup"
	"log"
)

type Repo interface {
	SetURL(ctx context.Context, userID uint32, url *entity.URL) error
	GetURL(ctx context.Context, shortURL string) (*entity.URL, error)
	GetUserURLList(ctx context.Context, id uint32) ([]*entity.URL, error)
	PingBD(ctx context.Context) bool
	Close(context.Context) error
	AddBatch(ctx context.Context, userID uint32, list []*entity.BatchURLItem) error
	DeleteBatch(ctx context.Context, userID uint32, listShortURL []string) error
	CheckUserBatch(ctx context.Context, userID uint32, listShortURL []string) ([]string, error)
}

type Service struct {
	baseURL     string
	repo        Repo
	jobChDelete chan *JobDelete
}

type JobDelete struct {
	UserID uint32
	URLs   []string
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

func (s *Service) CreateDeleteWorkerPull(ctx context.Context, nOfWorkers int) {
	s.jobChDelete = make(chan *JobDelete)
	g, _ := errgroup.WithContext(ctx)

	for i := 1; i <= nOfWorkers; i++ {
		g.Go(func() error {
			for job := range s.jobChDelete {
				err := s.repo.DeleteBatch(ctx, job.UserID, job.URLs)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
			}
			return nil
		})
	}

	go func() {
		<-ctx.Done()
		close(s.jobChDelete)
	}()

	if err := g.Wait(); err != nil {
		log.Println(err)
	}
}

func (s *Service) DeleteURLList(ctx context.Context, userID uint32, checksums []string) error {
	list := make([]string, 0, len(checksums))
	for _, checksum := range checksums {
		shortURL := entity.CreateShortURL(checksum, s.baseURL)
		list = append(list, shortURL)
	}
	var err error
	list, err = s.repo.CheckUserBatch(ctx, userID, list)
	if err != nil {
		return err
	}
	if len(list) <= 0 {
		return nil
	}
	job := &JobDelete{
		UserID: userID,
		URLs:   list,
	}
	if s.jobChDelete != nil {
		s.jobChDelete <- job
	}
	return nil
}
