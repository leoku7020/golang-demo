package httpkit

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"demo/pkg/logger"

	"github.com/gorilla/schema"
)

func NewSender() Sender {
	return &sender{}
}

type sender struct{}

func (s *sender) Send(
	ctx context.Context, method, path string, header http.Header, value interface{},
) (int, []byte, error) {

	u, err := url.Parse(path)
	if err != nil {
		logger.Ctx(ctx).Error("failed to parse url", logger.WithError(err))
		return 0, nil, err
	}

	var req *http.Request

	switch header.Get("Content-Type") {
	case "application/x-www-form-urlencoded":
		req, err = reqWithURLForm(ctx, method, u.String(), value)
	default:
		req, err = reqWithJSON(ctx, method, u.String(), value)
	}

	if err != nil {
		logger.Ctx(ctx).Error("failed to create new request", logger.WithField("error", err))
		return 0, nil, err
	}

	for k, vs := range header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Ctx(ctx).Error("failed to send request", logger.WithField("error", err))
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Ctx(ctx).Error("failed to read response", logger.WithField("error", err))
		return 0, nil, err
	}

	return resp.StatusCode, body, nil
}

func (s *sender) SendReturnCSV(
	ctx context.Context, method, path string, header http.Header, value interface{},
) (int, [][]string, error) {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", path, nil)
	resp, err := client.Do(req)
	if err != nil {
		logger.Ctx(ctx).Error("failed to send request", logger.WithField("error", err))
		return 0, nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.FieldsPerRecord = -1 // 不檢查每一行的字段數量

	body, err := reader.ReadAll()
	if err != nil {
		logger.Ctx(ctx).Error("failed to read response", logger.WithField("error", err))
		return 0, nil, err
	}

	return resp.StatusCode, body, nil
}

func reqWithURLForm(ctx context.Context, method, uri string, v interface{}) (*http.Request, error) {
	// http will wrap the body as URLEncoded if the content-type is application/x-www-form-urlencoded
	//
	// --ATTENTION--
	// The values in the struct should add `schema` tag not `json` tag
	// e.g.
	// type Example struct {
	// 	Name string `schema:"name"`
	// 	Age  int    `schema:"age"`
	// }
	// As above example, the post body will be wrapped as name=xxx&age=xxx
	form := url.Values{}

	encoder := schema.NewEncoder()
	if err := encoder.Encode(v, form); err != nil {
		return nil, err
	}

	return http.NewRequestWithContext(ctx, method, uri, strings.NewReader(form.Encode()))
}

func reqWithJSON(ctx context.Context, method, uri string, v interface{}) (*http.Request, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		logger.Ctx(ctx).Error("failed to marshal", logger.WithError(err))
		return nil, err
	}
	logger.Ctx(ctx).Debug("Marshal json format", logger.WithFields(map[string]interface{}{
		"url":    uri,
		"method": method,
		"body":   string(bs),
	}))

	var body io.Reader
	if bytes.NewBuffer(bs).String() != "null" {
		body = bytes.NewBuffer(bs)
	}

	return http.NewRequestWithContext(ctx, method, uri, body)
}
