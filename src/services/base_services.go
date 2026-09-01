package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/MohammadrezaNadirkhanloo/Go-project/common"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/constans"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/db"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/logging"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/service_errors"
	"gorm.io/gorm"
)

type BaseService[TEntity any, TCreate any, TUpdate any, TResponse any] struct {
	Database *gorm.DB
	logger   logging.Logger
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

// func (u *BaseUsecase[TEntity, TCreate, TUpdate, TResponse]) GetByFilter(ctx context.Context, req filter.PaginationInputWithFilter) (*filter.PagedList[TResponse], error) {
// 	var response *filter.PagedList[TResponse]
// 	count, entities, err := u.repository.GetByFilter(ctx, req)
// 	if err != nil {
// 		return response, err
// 	}

// 	return filter.Paginate[TEntity, TResponse](count, entities, req.PageNumber, int64(req.PageSize))
// }
