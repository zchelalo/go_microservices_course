package course

import (
	"context"
	"log"
	"time"

	"github.com/zchelalo/go_microservices_domain/domain"
)

type (
	Filters struct {
		Name string
	}

	Service interface {
		Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Get(ctx context.Context, id string) (*domain.Course, error)
		Update(ctx context.Context, id string, name, startDate, endDate *string) error
		Delete(ctx context.Context, id string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	service struct {
		log        *log.Logger
		repository Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:        log,
		repository: repo,
	}
}

func (srv *service) Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error) {
	srv.log.Println("create course service")

	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		srv.log.Println(err)
		return nil, ErrInvalidStartDate
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		srv.log.Println(err)
		return nil, ErrInvalidEndDate
	}

	course := domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}
	if err := srv.repository.Create(ctx, &course); err != nil {
		return nil, err
	}
	return &course, nil
}

func (srv *service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	srv.log.Println("get all courses service")
	courses, err := srv.repository.GetAll(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (srv *service) Get(ctx context.Context, id string) (*domain.Course, error) {
	srv.log.Println("get course service")
	course, err := srv.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (srv *service) Update(ctx context.Context, id string, name *string, startDate, endDate *string) error {
	srv.log.Println("update course service")

	var startDateParsed *time.Time
	if startDate != nil {
		parsed, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			srv.log.Println(err)
			return ErrInvalidStartDate
		}
		startDateParsed = &parsed
	}

	var endDateParsed *time.Time
	if endDate != nil {
		parsed, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			srv.log.Println(err)
			return ErrInvalidEndDate
		}
		endDateParsed = &parsed
	}

	return srv.repository.Update(ctx, id, name, startDateParsed, endDateParsed)
}

func (srv *service) Delete(ctx context.Context, id string) error {
	srv.log.Println("delete course service")
	return srv.repository.Delete(ctx, id)
}

func (srv *service) Count(ctx context.Context, filters Filters) (int, error) {
	srv.log.Println("count course service")
	return srv.repository.Count(ctx, filters)
}
