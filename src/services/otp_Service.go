package services

import (
	"context"
	"fmt"
	"time"

	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/constans"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/cache"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/logging"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/service_errors"
	"github.com/redis/go-redis/v9"
)

type OtpUsecase struct {
	logger      logging.Logger
	cfg         *config.Config
	redisClient *redis.Client
}

type otpDto struct {
	Value string
	Used  bool
}

func NewOtpUsecase(cfg *config.Config) *OtpUsecase {
	logger := logging.NewLogger(cfg)
	redis := cache.GetRedis()
	return &OtpUsecase{logger: logger, cfg: cfg, redisClient: redis}
}

// func (u *OtpUsecase) SendOtp(mobileNumber string) error {
// 	otp := common.GenerateOtp()
// 	err := u.SetOtp(mobileNumber, otp)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (u *OtpUsecase) SetOtp(ctx context.Context, mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constans.RedisOtpDefaultKey, mobileNumber)
	val := &otpDto{
		Value: otp,
		Used:  false,
	}

	res, err := cache.Get[otpDto](ctx, u.redisClient, key)
	if err == nil && !res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OptExists}
	} else if err == nil && res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	}

	err = cache.Set(ctx, u.redisClient, key, val, u.cfg.Otp.ExpireTime*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (u *OtpUsecase) ValidateOtp(ctx context.Context, mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constans.RedisOtpDefaultKey, mobileNumber)

	res, err := cache.Get[otpDto](ctx, u.redisClient, key)
	if err != nil {
		return err
	} else if res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	} else if !res.Used && res.Value != otp {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpNotValid}
	} else if !res.Used && res.Value == otp {
		res.Used = true
		err = cache.Set(ctx, u.redisClient, key, res, u.cfg.Otp.ExpireTime*time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}
