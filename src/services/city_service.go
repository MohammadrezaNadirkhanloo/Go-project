package services

import (
	"context"

	"github.com/MohammadrezaNadirkhanloo/Go-project/api/dto"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/db"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/models"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/logging"
)

type CityService struct {
	base *BaseService[models.City, dto.CreateAndUpdateCityRequest, dto.CreateAndUpdateCityRequest, dto.CityResponse]
}

func NewCityService(cfg *config.Config) *CityService {
	return &CityService{
		base: &BaseService[models.City, dto.CreateAndUpdateCityRequest, dto.CreateAndUpdateCityRequest, dto.CityResponse]{
			Database: db.GetDB(),
			logger:   logging.NewLogger(cfg),
			Preload: []preload{
				{string: "Country"},
			},
		},
	}
}

// Create
func (s *CityService) Create(ctx context.Context, req *dto.CreateAndUpdateCityRequest) (*dto.CityResponse, error) {
	return s.base.Create(ctx, req)
}

// Update
func (s *CityService) Update(ctx context.Context, id int, req *dto.CreateAndUpdateCityRequest) (*dto.CityResponse, error) {
	return s.base.Update(ctx, id, req)
}

// Delete
func (s *CityService) Delete(ctx context.Context, id int) error {
	return s.base.Delete(ctx, id)
}

// Get By Id
func (s *CityService) GetById(ctx context.Context, id int) (*dto.CityResponse, error) {
	return s.base.GetById(ctx, id)
}

// Get By Filter
func (s *CityService) GetByFilter(ctx context.Context, req *dto.PaginationInputWithFilter) (*dto.PagedList[dto.CityResponse], error) {
	return s.base.GetByFilter(ctx, req)
}
