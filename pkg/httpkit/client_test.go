package httpkit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"gs-mono/pkg/logger"
)

var (
	mockCTX = context.Background()
)

type httpSuite struct {
	suite.Suite

	sender Sender
}

func (s *httpSuite) SetupSuite() {
	// prepare logger
	logger.Register(logger.Config{
		Exporter:    logger.ExporterZap,
		Development: true,
	})
}

func (s *httpSuite) TearDownSuite() {
	logger.Flush()
}

func (s *httpSuite) SetupTest() {
	s.sender = NewSender()
}

func (s *httpSuite) TearDownTest() {}

func TestHTTPSuite(t *testing.T) {
	suite.Run(t, new(httpSuite))
}

func (s *httpSuite) TestSend() {
	tests := []struct {
		Desc        string
		Method      string
		Value       interface{}
		URL         string
		ExpCode     int
		ExpErrorStr string
	}{
		{
			Desc:        "json.Marshal failed",
			Method:      http.MethodPost,
			Value:       make(chan int),
			ExpCode:     400,
			ExpErrorStr: "json: unsupported type: chan int",
		},
		{
			Desc:        "url.Parse failed",
			Method:      http.MethodPost,
			URL:         "\t",
			Value:       nil,
			ExpCode:     400,
			ExpErrorStr: "parse \"\\t\": net/url: invalid control character in URL",
		},
		{
			Desc:        "client.Do failed",
			Method:      http.MethodPost,
			URL:         "google.com",
			Value:       nil,
			ExpCode:     400,
			ExpErrorStr: "Post \"google.com\": unsupported protocol scheme \"\"",
		},
		{
			Desc:    "get",
			Method:  http.MethodGet,
			ExpCode: 200,
		},
		{
			Desc:   "post",
			Method: http.MethodPost,
			Value: map[string]interface{}{
				"key":   "value",
				"array": []string{"a", "b"},
				"level1": map[string]interface{}{
					"level2": "string",
				},
			},
			ExpCode: 200,
		},
		{
			Desc:   "Method Not Allowed",
			Method: http.MethodGet,
			Value: map[string]interface{}{
				"key":   "value",
				"array": []string{"a", "b"},
				"level1": map[string]interface{}{
					"level2": "string",
				},
			},
			ExpCode: 405,
		},
	}

	for _, t := range tests {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(t.ExpCode)
		}))

		url := svr.URL
		if t.URL != "" {
			url = t.URL
		}

		code, _, err := s.sender.Send(mockCTX, t.Method, url, nil, t.Value)
		if err != nil {
			s.Require().Equal(t.ExpErrorStr, err.Error(), t.Desc)

			svr.Close()

			continue
		}

		s.Require().Equal(t.ExpCode, code, t.Desc)

		svr.Close()
	}
}
