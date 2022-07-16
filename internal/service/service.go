package service

import (
	"context"
	"log"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/AnnV0lokitina/short-url-service/internal/entity"
)

// Repo declare repository interface
type Repo interface {
	SetURL(ctx context.Context, userID uint32, url *entity.URL) error
	GetURL(ctx context.Context, shortURL string) (*entity.URL, error)
	GetUserURLList(ctx context.Context, id uint32) ([]*entity.URL, error)
	PingBD(ctx context.Context) bool
	Close(context.Context) error
	AddBatch(ctx context.Context, userID uint32, list []*entity.BatchURLItem) error
	DeleteBatch(ctx context.Context, userID uint32, listShortURL []string) error
	CheckUserBatch(ctx context.Context, userID uint32, listShortURL []string) ([]string, error)
	GetStats(ctx context.Context) (urls int, users int, err error)
}

// Service keep information to execute application tasks.
type Service struct {
	mu            sync.Mutex
	baseURL       string
	repo          Repo
	jobChDelete   chan *JobDelete
	trustedSubnet string
}

// JobDelete keep user and urls list to delete
type JobDelete struct {
	UserID uint32   // user identifier
	URLs   []string // urls list
}

// NewService Create new Service struct.
func NewService(baseURL string, repo Repo, trustedSubnet string) *Service {
	return &Service{
		baseURL:       baseURL,
		repo:          repo,
		trustedSubnet: trustedSubnet,
	}
}

// GetBaseURL Return the application base url.
func (s *Service) GetBaseURL() string {
	return s.baseURL
}

// SetBaseURL Set the application base url.
func (s *Service) SetBaseURL(baseURL string) {
	s.baseURL = baseURL
}

// GetRepo Return repository struct.
func (s *Service) GetRepo() Repo {
	return s.repo
}

// CreateDeleteWorkerPull Initialize pull of workers to delete urls list.
func (s *Service) CreateDeleteWorkerPull(ctx context.Context, nOfWorkers int) {
	s.mu.Lock()
	s.jobChDelete = make(chan *JobDelete)
	s.mu.Unlock()
	g, _ := errgroup.WithContext(ctx)

	for i := 1; i <= nOfWorkers; i++ {
		g.Go(func() error {
			for job := range s.jobChDelete {
				err := s.repo.DeleteBatch(ctx, job.UserID, job.URLs)
				if err != nil {
					log.Println(err.Error())
					continue
				}
			}
			return nil
		})
	}

	go func() {
		<-ctx.Done()
		s.mu.Lock()
		close(s.jobChDelete)
		s.mu.Unlock()
	}()

	if err := g.Wait(); err != nil {
		log.Println(err)
	}
}

// DeleteURLList Delete urls list, sent by user.
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
	s.mu.Lock()
	if s.jobChDelete != nil {
		s.jobChDelete <- job
	}
	s.mu.Unlock()
	return nil
}
