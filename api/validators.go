package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/meomeocoj/simplebank/utils"
)

var validateCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if curr, ok := fieldLevel.Field().Interface().(string); ok {
		return utils.IsCurrency(curr)
	}
	return false
}
