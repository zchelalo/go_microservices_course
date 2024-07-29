package course

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/zchelalo/go_microservices_domain/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(ctx context.Context, course *domain.Course) error
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Get(ctx context.Context, id string) (*domain.Course, error)
		Update(ctx context.Context, id string, name *string, startDate, endDate *time.Time) error
		Delete(ctx context.Context, id string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	repository struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepository(log *log.Logger, db *gorm.DB) Repository {
	return &repository{
		log: log,
		db:  db,
	}
}

func (repo *repository) Create(ctx context.Context, course *domain.Course) error {
	if err := repo.db.WithContext(ctx).Create(course).Error; err != nil {
		repo.log.Println(err)
		return err
	}

	repo.log.Println("course created with id: ", course.Id)
	return nil
}

func (repo *repository) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	var courses []domain.Course

	tx := repo.db.WithContext(ctx).Model(&courses)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	if err := tx.Order("created_at desc").Find(&courses).Error; err != nil {
		repo.log.Println(err)
		return nil, err
	}

	return courses, nil
}

func (repo *repository) Get(ctx context.Context, id string) (*domain.Course, error) {
	course := domain.Course{
		Id: id,
	}

	if err := repo.db.WithContext(ctx).Model(&course).First(&course).Error; err != nil {
		repo.log.Println(err)
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound{id}
		}
		return nil, err
	}

	return &course, nil
}

func (repo *repository) Update(ctx context.Context, id string, name *string, startDate, endDate *time.Time) error {
	values := make(map[string]interface{})

	if name != nil {
		values["name"] = *name
	}

	if startDate != nil {
		values["start_date"] = *startDate
	}

	if endDate != nil {
		values["end_date"] = *endDate
	}

	result := repo.db.WithContext(ctx).Model(&domain.Course{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		repo.log.Printf("course with id %s doesn't exists", id)
		return ErrNotFound{id}
	}

	return nil
}

func (repo *repository) Delete(ctx context.Context, id string) error {
	course := domain.Course{
		Id: id,
	}

	result := repo.db.WithContext(ctx).Model(&course).Delete(&course)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		repo.log.Printf("course with id %s doesn't exists", id)
		return ErrNotFound{id}
	}

	return nil
}

func (repo *repository) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(&domain.Course{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		repo.log.Println(err)
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}

	return tx
}
