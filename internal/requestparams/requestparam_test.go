package requestparams

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testExpressions = map[string]string{
	"value":                "value",
	"=value":               "value",
	"key=value":            "key=value",
	"key=value with space": "key=value+with+space",
	"value with space":     "value+with+space",
	"@testfile":            "test+foo%0D%0Abar",
	"name@testfile":        "name=test+foo%0D%0Abar",
}

func TestFromUnencodedExpression(t *testing.T) {
	var result RequestParam
	var err error

	for expression, expected := range testExpressions {
		result, err = FromUnencodedExpression(expression)
		assert.NoError(t, err)
		assert.Equal(t, expected, result.String())
	}
}
