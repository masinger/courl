package request

import (
	"fmt"
	"github.com/masinger/courl/internal/requestparams"
	"github.com/masinger/courl/internal/util"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	Method                   string
	Headers                  []string
	LiteralDataExpressions   []string
	UnencodedDataExpressions []string
	DataBinaryExpressions    []string
}

func (r Request) CreateRequest(url string) (*http.Request, error) {
	body, suggestedRequestMethod, suggestedContentType, err := r.body()
	if err != nil {
		return nil, err
	}

	urlWithProtocol := prependProtocol(url)

	req, err := http.NewRequest(
		util.PresentOrDefault(util.CoalesceStrings(util.NotEmpty(&r.Method, &suggestedRequestMethod)...), http.MethodGet),
		urlWithProtocol,
		body,
	)

	if util.AllStringsPresent(&suggestedContentType) && !r.hasCustomContentType() {
		req.Header.Set("Content-Type", suggestedContentType)
	}

	if err = r.applyHeaders(req); err != nil {
		return req, err
	}

	log.Debug(">>>>>>>>>>>>>> REQUEST >>>>>>>>>>>>>>>>>>>>>>")
	log.Debugf("%s %s", req.Method, req.URL.String())
	for headerName, headerValues := range req.Header {
		for _, headerValue := range headerValues {
			log.Debugf("%s: %s", headerName, headerValue)
		}
	}
	log.Debug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	return req, err
}

func (r Request) body() (data io.Reader, suggestedRequestMethod string, suggestedContentType string, err error) {
	postParamsReader, err := r.form()
	if err != nil {
		return nil, "", "", err
	}

	binaryReader, err := r.binary()
	if err != nil {
		return nil, "", "", err
	}

	if postParamsReader == nil && binaryReader == nil {
		return nil, "", "", nil
	}

	var resultReader io.Reader
	if postParamsReader == nil || binaryReader == nil {
		if postParamsReader != nil {
			resultReader = postParamsReader
		} else {
			resultReader = binaryReader
		}
	} else {
		resultReader = io.MultiReader(postParamsReader, binaryReader)
	}

	return resultReader, http.MethodPost, "application/x-www-form-urlencoded", nil
}

func (r Request) binary() (io.Reader, error) {
	var readers []io.Reader
	for _, binaryExpression := range r.DataBinaryExpressions {
		reader, err := util.GetRawOrStream(binaryExpression)
		if err != nil {
			return nil, err
		}
		readers = append(readers, reader)
	}
	if readers == nil || len(readers) == 0 {
		return nil, nil
	}
	if len(readers) == 1 {
		return readers[0], nil
	}
	return io.MultiReader(readers...), nil
}

func (r Request) form() (io.Reader, error) {
	var result requestparams.RequestParams
	for _, expression := range r.LiteralDataExpressions {
		param, err := requestparams.FromLiteralExpression(expression)
		if err != nil {
			return nil, err
		}
		result = append(result, param)
	}
	for _, expression := range r.UnencodedDataExpressions {
		param, err := requestparams.FromUnencodedExpression(expression)
		if err != nil {
			return nil, err
		}
		result = append(result, param)
	}
	if result == nil || len(result) == 0 {
		return nil, nil
	}
	return result.Reader(), nil
}

func (r Request) hasCustomContentType() bool {
	for _, header := range r.Headers {
		if strings.Index(strings.TrimSpace(strings.ToLower(header)), "content-type:") == 0 {
			return true
		}
	}
	return false
}

func (r Request) applyHeaders(req *http.Request) error {
	for _, header := range r.Headers {
		firstIndexOfDelimiter := strings.Index(header, ":")
		if firstIndexOfDelimiter == -1 {
			return fmt.Errorf("invalid header: %s", header)
		}
		headerName := strings.TrimSpace(header[:firstIndexOfDelimiter])
		headerValue := strings.TrimSpace(header[firstIndexOfDelimiter+1:])

		req.Header.Add(headerName, headerValue)
	}

	return nil
}

func prependProtocol(url string) string {
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		return url
	}
	if strings.HasPrefix(url, "://") {
		return fmt.Sprintf("https%s", url)
	}
	if strings.HasPrefix(url, "//") {
		return fmt.Sprintf("https:%s", url)
	}
	return fmt.Sprintf("https://%s", url)
}
