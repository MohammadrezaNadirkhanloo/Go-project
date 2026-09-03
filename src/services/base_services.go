package services

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/MohammadrezaNadirkhanloo/Go-project/api/dto"
	"github.com/MohammadrezaNadirkhanloo/Go-project/common"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/constans"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/db"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/logging"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/service_errors"
	"gorm.io/gorm"
)

type preload struct {
	string
}
type BaseService[TEntity any, TCreate any, TUpdate any, TResponse any] struct {
	Database *gorm.DB
	logger   logging.Logger
	Preload  []preload
}

func NewBaseService[TEntity any, TCreate any, TUpdate any, TResponse any](cfg *config.Config) *BaseService[TEntity, TCreate, TUpdate, TResponse] {
	return &BaseService[TEntity, TCreate, TUpdate, TResponse]{
		Database: db.GetDB(),
		logger:   logging.NewLogger(cfg),
	}
}

func (s *BaseService[TEntity, TCreate, TUpdate, TResponse]) Create(ctx context.Context, req *TCreate) (*TResponse, error) {
	model, _ := common.TypeConverter[TEntity](req)
	tx := s.Database.WithContext(ctx).Begin()

	err := tx.Create(model).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return nil, err
	}
	tx.Commit()
	return common.TypeConverter[TResponse](model)
}

func (s *BaseService[TEntity, TCreate, TUpdate, TResponse]) Update(ctx context.Context, id int, req *TUpdate) (*TResponse, error) {

	updateMap, _ := common.TypeConverter[map[string]interface{}](req)
	(*updateMap)["modified_by"] = &sql.NullInt64{Int64: int64(ctx.Value(constans.UserIdKey).(float64)), Valid: true}
	(*updateMap)["modified_at"] = sql.NullTime{Valid: true, Time: time.Now().UTC()}

	model := new(TEntity)
	tx := s.Database.WithContext(ctx).Begin()

	err := tx.Model(model).Where("id = ? and deleted_by is null", id).Updates(*updateMap).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return nil, err
	}
	tx.Commit()
	s.GetById(ctx, id)
	return s.GetById(ctx, id)
}

func (s *BaseService[TEntity, TCreate, TUpdate, TResponse]) Delete(ctx context.Context, id int) error {
	tx := s.Database.WithContext(ctx).Begin()
	model := new(TEntity)

	deleteMap := map[string]interface{}{
		"deleted_by": &sql.NullInt64{Int64: int64(ctx.Value(constans.UserIdKey).(float64)), Valid: true},
		"deleted_at": sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}
	deleteMap["modified_by"] = &sql.NullInt64{Int64: int64(ctx.Value(constans.UserIdKey).(float64)), Valid: true}
	deleteMap["modified_at"] = sql.NullTime{Valid: true, Time: time.Now().UTC()}
	if ctx.Value(constans.UserIdKey) == nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	cnt := tx.Model(model).Where("id = ? and deleted_by is null", id).Updates(deleteMap).RowsAffected
	if cnt == 0 {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Update, service_errors.RecordNotFound, nil)
		return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
	}
	tx.Commit()
	return nil
}

func (s *BaseService[TEntity, TCreate, TUpdate, TResponse]) GetById(ctx context.Context, id int) (*TResponse, error) {
	model := new(TEntity)
	err := s.Database.Where("id = ? and deleted_by is null", id).First(model).Error
	if err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return nil, err
	}
	return common.TypeConverter[TResponse](model)
}

func (s *BaseService[T, Tc, Tu, Tr]) GetByFilter(ctx context.Context, req *dto.PaginationInputWithFilter) (*dto.PagedList[Tr], error) {
	return Paginate[T, Tr](req, s.Preload, s.Database)

}

func NewPagedList[T any](items *[]T, count int64, pageNumber int, pageSize int64) *dto.PagedList[T] {
	pl := &dto.PagedList[T]{
		PageNumber: pageNumber,
		TotalRows:  count,
		Items:      items,
	}
	pl.TotalPages = int(math.Ceil(float64(count) / float64(pageSize)))
	pl.HasNextPage = pl.PageNumber < pl.TotalPages
	pl.HasPreviousPage = pl.PageNumber > 1

	return pl
}

// Paginate
func Paginate[T any, Tr any](pagination *dto.PaginationInputWithFilter, preloads []preload, db *gorm.DB) (*dto.PagedList[Tr], error) {
	model := new(T)
	var items *[]T
	var rItems *[]Tr
	db = Preload(db, preloads)
	query := getQuery[T](&pagination.DynamicFilter)
	sort := getSort[T](&pagination.DynamicFilter)

	var totalRows int64 = 0

	db.
		Model(model).
		Where(query).
		Count(&totalRows)

	err := db.
		Where(query).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetPageSize()).
		Order(sort).
		Find(&items).
		Error

	if err != nil {
		return nil, err
	}
	rItems, err = common.TypeConverter[[]Tr](items)
	if err != nil {
		return nil, err
	}
	return NewPagedList(rItems, totalRows, pagination.PageNumber, int64(pagination.PageSize)), err

}

func getQuery[T any](filter *dto.DynamicFilter) string {
	t := new(T)
	typeT := reflect.TypeOf(*t)
	query := make([]string, 0)
	query = append(query, "deleted_by is null")
	if filter.Filter != nil {
		for name, filter := range filter.Filter {
			fld, ok := typeT.FieldByName(name)
			if ok {
				fld.Name = common.ToSnakeCase(fld.Name)
				switch filter.Type {
				case "contains":
					query = append(query, fmt.Sprintf("%s ILike '%%%s%%'", fld.Name, filter.From))
				case "notContains":
					query = append(query, fmt.Sprintf("%s not ILike '%%%s%%'", fld.Name, filter.From))

				case "startsWith":
					query = append(query, fmt.Sprintf("%s ILike '%s%%'", fld.Name, filter.From))

				case "endsWith":
					query = append(query, fmt.Sprintf("%s ILike '%%%s'", fld.Name, filter.From))

				case "equals":
					query = append(query, fmt.Sprintf("%s = '%s'", fld.Name, filter.From))

				case "notEquals":
					query = append(query, fmt.Sprintf("%s != '%s'", fld.Name, filter.From))

				case "lessThan":
					query = append(query, fmt.Sprintf("%s < '%s'", fld.Name, filter.From))

				case "lessThanOrEqual":
					query = append(query, fmt.Sprintf("%s <= '%s'", fld.Name, filter.From))

				case "greaterThan":
					query = append(query, fmt.Sprintf("%s > '%s'", fld.Name, filter.From))

				case "greaterThanOrEqual":
					query = append(query, fmt.Sprintf("%s >= '%s'", fld.Name, filter.From))

				case "inRange":
					if fld.Type.Kind() == reflect.String {
						query = append(query, fmt.Sprintf("%s >= '%s'", fld.Name, filter.From))
						query = append(query, fmt.Sprintf("%s <= '%s'", fld.Name, filter.To))
					} else {
						query = append(query, fmt.Sprintf("%s >= %s", fld.Name, filter.From))
						query = append(query, fmt.Sprintf("%s <= %s", fld.Name, filter.To))
					}
				}
			}
		}
	}
	return strings.Join(query, " AND ")
}

func getSort[T any](filter *dto.DynamicFilter) string {
	t := new(T)
	typeT := reflect.TypeOf(*t)
	sort := make([]string, 0)
	if filter.Sort != nil {
		for _, tp := range *filter.Sort {
			fld, ok := typeT.FieldByName(tp.ColId)
			if ok && (tp.Sort == "asc" || tp.Sort == "desc") {
				fld.Name = common.ToSnakeCase(fld.Name)
				sort = append(sort, fmt.Sprintf("%s %s", fld.Name, tp.Sort))
			}
		}
	}
	return strings.Join(sort, ", ")
}

func Preload(db *gorm.DB, preloads []preload) *gorm.DB {
	for _, item := range preloads {
		db = db.Preload(item.string)
	}
	return db
}
