package response

import (
	"github.com/go-playground/validator/v10"
)

func ExtractValidationErrors(err error) map[string]string {
	out := map[string]string{}
	if err == nil {
		return out
	}
	if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range verrs {
			out[fe.Field()] = fe.Tag()
		}
		return out
	}
	out["error"] = err.Error()
	return out
}
