package response

import (
	"fmt"
	"github.com/masinger/courl/internal/util"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	AllowBinaryToStdout bool
	OutputPath          string
	Verbose             bool
	ErrorOnFail         bool
}

func (r Response) Handle(httpResponse *http.Response, responseError error) error {
	if responseError != nil {
		return responseError
	}

	if err := r.verbose(httpResponse); err != nil {
		return err
	}

	if r.ErrorOnFail {
		if respErr := toResponseError(httpResponse); respErr != nil {
			return respErr
		}
	}

	return r.body(httpResponse)
}

func (r Response) verbose(response *http.Response) error {
	if !r.Verbose {
		return nil
	}
	log.Debug("<<<<<<<<<<<<<< RESPONSE <<<<<<<<<<<<<<<<<<<<<")
	log.Debugf("%s %d", response.Proto, response.StatusCode)
	for headerName, headerValues := range response.Header {
		for _, headerValue := range headerValues {
			log.Debugf("%s: %s", headerName, headerValue)
		}
	}
	log.Debug("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	return nil
}

func (r Response) body(response *http.Response) (err error) {
	body := response.Body
	defer func() {
		if tempErr := body.Close(); tempErr != nil {
			err = tempErr
		}
	}()
	return r.writeOutput(response, body)
}

func (r Response) writeOutput(response *http.Response, reader io.Reader) (err error) {
	var writer io.Writer
	if util.AllStringsPresent(&r.OutputPath) {
		f, err := os.Create(r.OutputPath)
		if err != nil {
			return err
		}
		defer func() {
			if tempErr := f.Close(); tempErr != nil {
				err = tempErr
			}
		}()
		writer = f
	} else {
		if !r.AllowBinaryToStdout {
			if isBinary, indicator := isBinaryData(response); isBinary {
				return fmt.Errorf("will not output binary data to stdout unless --binary-stdout is set, reason: %s", indicator)
			}
		}
		writer = os.Stdout
	}
	_, err = io.Copy(writer, reader)
	return err
}

func toResponseError(resp *http.Response) error {
	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}
	return nil
}

func isBinaryData(resp *http.Response) (bool, string) {
	contentDisposition := resp.Header.Get("Content-Disposition")
	if strings.Index(strings.TrimSpace(contentDisposition), "attachment;") == 0 {
		return true, "response has 'Content-Disposition' header set to 'attachment'"
	}

	return false, ""
}
