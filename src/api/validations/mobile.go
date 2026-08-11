package validation
import (
    "log"
    "regexp"
    "github.com/go-playground/validator/v10"
)
func IranMobileNumberValidator(fld validator.FieldLevel) bool {
    val, ok := fld.Field().Interface().(string)// گرفتن مقدار و بررسی دیتا بودن
    if !ok {
        return false
    }
    res, err := regexp.MatchString(`^09(1[0-9]|2[0-2]|3[0-9]|9[0-9])[0-9]{7}$`, val)
    if err != nil {
        log.Print(err.Error())
    }
    return res
}