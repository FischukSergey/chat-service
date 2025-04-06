package validator

import (
	"github.com/go-playground/validator/v10"
	optsGenValidator "github.com/kazhuravlev/options-gen/pkg/validator"
)

// Validator представляет валидатор.
var Validator = validator.New()

// init инициализирует валидатор.
func init() {
	optsGenValidator.Set(Validator)
}
