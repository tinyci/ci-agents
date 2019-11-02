package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
)

// JSONIO is a copy function that uses JSON as an intermediary.
func JSONIO(from, to interface{}) error {
	// Yes. I know it sucks. I blame go-swagger.
	content, err := json.Marshal(from)
	if err != nil {
		return errors.New(err)
	}

	return errors.New(json.Unmarshal(content, to))
}

// JSONContext provides an easy method to extract json from gin parameters.
func JSONContext(ctx *gin.Context, parameter string, obj interface{}) error {
	content, ok := ctx.Get(parameter)
	if !ok {
		return errors.New(fmt.Sprintf("parameter %q not found", parameter))
	}
	return errors.New(json.Unmarshal(content.([]byte), obj))
}
