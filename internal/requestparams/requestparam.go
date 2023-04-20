package requestparams

import (
	"bytes"
	"fmt"
	"github.com/masinger/courl/internal/util"
	"io"
	"net/url"
	"strings"
)

type RequestParams []RequestParam

type RequestParam interface {
	String() string
}

type rawRequestParam struct {
	value string
}

func (r rawRequestParam) String() string {
	return r.value
}

func FromLiteral(literal string) RequestParam {
	return &rawRequestParam{
		value: strings.ReplaceAll(
			strings.ReplaceAll(literal, "\r", ""),
			"\n",
			"",
		),
	}
}

func FromLiteralExpression(value string) (RequestParam, error) {
	data := bytes.Buffer{}

	_, err := util.CopyRawOrStream(value, &data)
	if err != nil {
		return nil, err
	}

	return FromLiteral(data.String()), nil
}

func FromUnencodedExpression(expression string) (RequestParam, error) {
	separatorIndex := strings.Index(expression, "=")
	var key string
	var value string
	if separatorIndex == -1 {
		key = ""
		value = expression
	} else {
		key = expression[:separatorIndex]
		value = expression[separatorIndex+1:]
	}

	fileInterpolationIndex := strings.Index(expression, "@")
	if fileInterpolationIndex == 0 {
		key = ""
		value = expression
	} else if fileInterpolationIndex > 0 {
		key = expression[:fileInterpolationIndex]
		value = expression[fileInterpolationIndex:]
	}

	data := &bytes.Buffer{}
	_, err := util.CopyRawOrStream(value, data)
	if err != nil {
		return nil, err
	}
	var result RequestParam
	if len(key) == 0 {
		result = &rawRequestParam{
			value: url.QueryEscape(data.String()),
		}
	} else {
		result = &rawRequestParam{
			value: fmt.Sprintf(
				"%s=%s",
				key,
				url.QueryEscape(data.String()),
			),
		}
	}
	return result, nil
}

func (requestParams RequestParams) Reader() io.Reader {
	var encodedParams []string
	for _, requestParam := range requestParams {
		encodedParams = append(encodedParams, requestParam.String())
	}

	encoded := strings.Join(encodedParams, "&")
	return bytes.NewBufferString(encoded)
}
