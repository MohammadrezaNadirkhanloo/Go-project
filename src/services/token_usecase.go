package services

import (
	"time"

	"github.com/MohammadrezaNadirkhanloo/Go-project/api/dto"
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/constans"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/logging"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/service_errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type TokenUsecase struct {
	logger logging.Logger
	cfg    *config.Config
}

type tokenDto struct { //اطلاعاتی که از توکن میگیریم 
	UserId       int
	FirstName    string
	LastName     string
	Username     string
	MobileNumber string
	Email        string
	Roles        []string
}

func NewTokenUsecase(cfg *config.Config) *TokenUsecase {
	logger := logging.NewLogger(cfg)
	return &TokenUsecase{
		cfg:    cfg,
		logger: logger,
	}
}

func (u *TokenUsecase) GenerateToken(token tokenDto) (*dto.TokenDetail, error) {
	td := &dto.TokenDetail{}
	td.AccessTokenExpireTime = time.Now().Add(u.cfg.JWT.AccessTokenExpireDuration * time.Minute).Unix()
	td.RefreshTokenExpireTime = time.Now().Add(u.cfg.JWT.RefreshTokenExpireDuration * time.Minute).Unix()

	atc := jwt.MapClaims{}

	atc[constans.UserIdKey] = token.UserId
	atc[constans.FirstNameKey] = token.FirstName
	atc[constans.LastNameKey] = token.LastName
	atc[constans.UsernameKey] = token.Username
	atc[constans.EmailKey] = token.Email
	atc[constans.MobileNumberKey] = token.MobileNumber
	atc[constans.RolesKey] = token.Roles
	atc[constans.ExpireTimeKey] = td.AccessTokenExpireTime

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)

	var err error
	td.AccessToken, err = at.SignedString([]byte(u.cfg.JWT.Secret))

	if err != nil {
		return nil, err
	}

	rtc := jwt.MapClaims{}

	rtc[constans.UserIdKey] = token.UserId
	rtc[constans.FirstNameKey] = token.FirstName
	rtc[constans.LastNameKey] = token.LastName
	rtc[constans.UsernameKey] = token.Username
	rtc[constans.EmailKey] = token.Email
	rtc[constans.MobileNumberKey] = token.MobileNumber
	rtc[constans.RolesKey] = token.Roles
	rtc[constans.ExpireTimeKey] = td.RefreshTokenExpireTime

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtc)

	td.RefreshToken, err = rt.SignedString([]byte(u.cfg.JWT.RefreshSecret))

	if err != nil {
		return nil, err
	}

	return td, nil
}

func (u *TokenUsecase) VerifyToken(token string) (*jwt.Token, error) {
	at, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
		}
		return []byte(u.cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return at, nil
}

func (u *TokenUsecase) GetClaims(token string) (claimMap map[string]interface{}, err error) { // اطلاعات گرفتن از توکن
	claimMap = map[string]interface{}{}

	verifyToken, err := u.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if ok && verifyToken.Valid {
		for k, v := range claims {
			claimMap[k] = v
		}
		return claimMap, nil
	}
	return nil, &service_errors.ServiceError{EndUserMessage: service_errors.ClaimsNotFound}
}

func (s *TokenUsecase) RefreshToken(c *gin.Context) (*dto.TokenDetail, error) {
	refreshToken, err := c.Cookie(constans.RefreshTokenCookieName)
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidRefreshToken}
	}

	claims, err := s.GetClaims(refreshToken)
	if err != nil {
		return nil, err
	}

	// Convert roles to []string
	rolesInterface, ok := claims[constans.RolesKey].([]interface{})
	if !ok {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidRolesFormat}
	}

	roles := make([]string, len(rolesInterface))
	for i, role := range rolesInterface {
		roles[i], ok = role.(string)
		if !ok {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidRolesFormat}
		}
	}

	tokenDto := tokenDto{
		UserId:       int(claims[constans.UserIdKey].(float64)),
		FirstName:    claims[constans.FirstNameKey].(string),
		LastName:     claims[constans.LastNameKey].(string),
		Username:     claims[constans.UsernameKey].(string),
		MobileNumber: claims[constans.MobileNumberKey].(string),
		Email:        claims[constans.EmailKey].(string),
		Roles:        roles,
	}
	newTokenDetail, err := s.GenerateToken(tokenDto)
	if err != nil {
		return nil, err
	}

	return newTokenDetail, nil
}
