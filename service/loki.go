package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var ErrNotFound = errors.New("not found")

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

func (l *Loki) get(params url.Values) ([]byte, error) {
	url := l.baseUrl + "/loki/api/v1/query_range?" + params.Encode()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	res, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func (l *Loki) Read(start, end time.Time, query string, limit int) ([][]string, error) {
	query = `{app="` + l.app + `"} ` + query

	params := url.Values{}
	params.Add("start", strconv.FormatInt(start.UnixNano(), 10))
	params.Add("end", strconv.FormatInt(end.UnixNano(), 10))
	params.Add("limit", strconv.Itoa(limit))
	params.Add("query", query)

	raw, err := l.get(params)
	if err != nil {
		return nil, err
	}

	q := &Query{}

	if err := json.Unmarshal(raw, q); err != nil {
		return nil, err
	}

	if len(q.Data.Result) != 1 {
		return nil, ErrNotFound
	}

	if len(q.Data.Result[0].Values) == 0 {
		return nil, errors.New("log miss")
	}

	return q.Data.Result[0].Values, nil
}
