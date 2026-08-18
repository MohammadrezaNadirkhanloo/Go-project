package services

import (
	"context"

	"github.com/MohammadrezaNadirkhanloo/Go-project/api/dto"
	"github.com/MohammadrezaNadirkhanloo/Go-project/common"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/constans"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/db"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/models"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/logging"
	"gorm.io/gorm"
)

type UserUsecase struct {
	logger     logging.Logger
	cfg        *config.Config
	otpUsecase *OtpUsecase
	database   *gorm.DB
	// tokenUsecase *TokenUsecase
	// repository   repository.UserRepository
}

func NewUserUsecase(cfg *config.Config) *UserUsecase {
	logger := logging.NewLogger(cfg)
	database := db.GetDB()
	return &UserUsecase{
		cfg: cfg,
		// repository:   repository,
		logger:     logger,
		database:   database,
		otpUsecase: NewOtpUsecase(cfg),
		// tokenUsecase: NewTokenUsecase(cfg),
	}
}

func (s *UserUsecase) SendOtp(ctx context.Context, req *dto.GetOtpRequest) error {
	otp, _ := common.GenerateOtp()
	err := s.otpUsecase.SetOtp(ctx, req.MobileNumber, otp)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserUsecase) existsByEmail(email string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.User{}).Select("count(*)>0").Where("email=?", email).Find(&exists).Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
	}
	return exists, nil
}
func (s *UserUsecase) existsByUsername(username string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.User{}).Select("count(*)>0").Where("username=?", username).Find(&exists).Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
	}
	return exists, nil
}
func (s *UserUsecase) existsByMobileNumber(mobileNumber string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.User{}).Select("count(*)>0").Where("mobile_number=?", mobileNumber).Find(&exists).Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
	}
	return exists, nil
}
func (s *UserUsecase) getDefaultRole() (roleId int, err error) {
	if err := s.database.Model(&models.Role{}).Select("id").Where("name=?", constans.DefaultRoleName).First(&roleId).Error; err != nil {
		return 0, err
	}
	return roleId, nil
}
