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
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/service_errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUsecase struct {
	logger       logging.Logger
	cfg          *config.Config
	otpUsecase   *OtpUsecase
	database     *gorm.DB
	tokenUsecase *TokenUsecase
	// repository   repository.UserRepository
}

func NewUserUsecase(cfg *config.Config) *UserUsecase {
	logger := logging.NewLogger(cfg)
	database := db.GetDB()
	return &UserUsecase{
		cfg: cfg,
		// repository:   repository,
		logger:       logger,
		database:     database,
		otpUsecase:   NewOtpUsecase(cfg),
		tokenUsecase: NewTokenUsecase(cfg),
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

// Register by username
func (u *UserUsecase) RegisterByUsername(req *dto.RegisterUserByUsernameRequest) error {
	user := models.User{Username: req.Username, Firstname: req.FirstName, Lastname: req.LastName, Email: req.Email, Password: req.Password}

	exists, err := u.existsByEmail(req.Email)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.EmailExists}
	}
	exists, err = u.existsByUsername(req.Username)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.UsernameExists}
	}

	bp := []byte(req.Password)
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return err
	}
	user.Password = string(hp)
	roleId, err := u.getDefaultRole()
	if err != nil {
		u.logger.Error(logging.General, logging.DefaultRoleNotFound, err.Error(), nil)
		return err
	}
	tx := u.database.Begin()
	err = tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		u.logger.Error(logging.General, logging.Rollback, err.Error(), nil)
		return err
	}
	err = tx.Create(&models.UserRole{RoleId: roleId, UserId: user.Id}).Error
	if err != nil {
		tx.Rollback()
		u.logger.Error(logging.General, logging.Rollback, err.Error(), nil)
		return err
	}
	tx.Commit()
	return nil
}

// Register/login by mobile number
func (u *UserUsecase) RegisterAndLoginByMobileNumber(ctx context.Context, req *dto.RegisterLoginByMobileRequest) (*dto.TokenDetail, error) {
	err := u.otpUsecase.ValidateOtp(ctx, req.MobileNumber, req.Otp)
	if err != nil {
		return nil, err
	}
	exists, err := u.existsByMobileNumber(req.MobileNumber)
	if err != nil {
		return nil, err
	}

	user := models.User{MobileNumber: req.MobileNumber, Username: req.MobileNumber}

	if exists {
		var user models.User
		err = u.database.Model(&models.User{}).Where("username = ?", user.Username).Preload("UserRole", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).Find(&user).Error

		if err != nil {
			return nil, err
		}
		tdto := tokenDto{UserId: user.Id, FirstName: user.Firstname, LastName: user.Lastname, Email: user.Email, MobileNumber: user.MobileNumber}

		if len(*user.UserRole) > 0 {
			for _, role := range *user.UserRole {
				tdto.Roles = append(tdto.Roles, role.Role.Name)
			}
		}

		token, err := u.tokenUsecase.GenerateToken(&tdto)
		if err != nil {
			return nil, err
		}
		return token, nil
	}

	bp := []byte(common.GeneratePassword())
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return nil, err
	}
	user.Password = string(hp)
	roleId, err := u.getDefaultRole()
	if err != nil {
		u.logger.Error(logging.General, logging.DefaultRoleNotFound, err.Error(), nil)
		return nil, err
	}
	tx := u.database.Begin()
	err = tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		u.logger.Error(logging.General, logging.Rollback, err.Error(), nil)
		return nil, err
	}
	err = tx.Create(&models.UserRole{RoleId: roleId, UserId: user.Id}).Error
	if err != nil {
		tx.Rollback()
		u.logger.Error(logging.General, logging.Rollback, err.Error(), nil)
		return nil, err
	}
	tx.Commit()

	err = u.database.Model(&models.User{}).Where("username = ?", user.Username).Preload("UserRole", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Role")
	}).Find(&user).Error

	if err != nil {
		return nil, err
	}
	tdto := tokenDto{UserId: user.Id, FirstName: user.Firstname, LastName: user.Lastname, Email: user.Email, MobileNumber: user.MobileNumber}

	if len(*user.UserRole) > 0 {
		for _, role := range *user.UserRole {
			tdto.Roles = append(tdto.Roles, role.Role.Name)
		}
	}

	token, err := u.tokenUsecase.GenerateToken(&tdto)
	if err != nil {
		return nil, err
	}
	return token, nil

}

func (u *UserUsecase) LoginByUsername(ctx context.Context, req *dto.LoginByUsernameRequest) (*dto.TokenDetail, error) {
	var user models.User
	err := u.database.Model(&models.User{}).Where("username = ?", user.Username).Preload("UserRole", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Role")
	}).Find(&user).Error
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(req.Password))
	tdto := tokenDto{UserId: user.Id, FirstName: user.Firstname, LastName: user.Lastname, Email: user.Email, MobileNumber: user.MobileNumber}

	if len(*user.UserRole) > 0 {
		for _, role := range *user.UserRole {
			tdto.Roles = append(tdto.Roles, role.Role.Name)
		}
	}

	token, err := u.tokenUsecase.GenerateToken(&tdto)
	if err != nil {
		return nil, err
	}
	return token, nil
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
