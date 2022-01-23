package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Loki struct {
	app     string
	baseUrl string
	client  *http.Client
}

func NewLoki(baseUrl, app string) *Loki {
	return &Loki{
		app:     app,
		baseUrl: baseUrl,
		client:  http.DefaultClient,
	}
}

func makeReq(url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (l Loki) push(raw []byte) error {
	url := l.baseUrl + "/loki/api/v1/push"

	req, err := makeReq(url, raw)
	if err != nil {
		return err
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return fmt.Errorf("send failed :%w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	return nil
}

func (l Loki) Push(objects ...interface{}) error {
	time := strconv.FormatInt(time.Now().UnixNano(), 10)

	values := make([][]string, len(objects))

	for i, o := range objects {
		raw, err := json.Marshal(o)
		if err != nil {
			return fmt.Errorf("marshaling object failed: %w", err)
		}

		values[i] = []string{time, string(raw)}
	}

	log, err := json.Marshal(NewLog(l.app, values))
	if err != nil {
		return fmt.Errorf("marshaling log failed: %w", err)
	}

	return l.push(log)
}
