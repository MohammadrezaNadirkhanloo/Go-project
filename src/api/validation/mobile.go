package validation

import (
	"regexp"

	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/logging"
	"github.com/go-playground/validator/v10"
)

var logger = logging.NewLogger(config.GetConfig())

func IranMobileNumberValidator(fld validator.FieldLevel) bool {
	val, ok := fld.Field().Interface().(string) // گرفتن مقدار و بررسی دیتا بودن
	if !ok {
		return false
	}
	res, err := regexp.MatchString(`^09(1[0-9]|2[0-2]|3[0-9]|9[0-9])[0-9]{7}$`, val)
	if err != nil {
		logger.Error(logging.Validation, logging.MobileValidation, err.Error(), nil)
	}
	return res
}
