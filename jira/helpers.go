package jira

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	return client
}

func makeRequest(method, u string, reqBody []byte, c *Config) (*http.Response, error) {
	req, err := http.NewRequest(method, u, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.Username, c.Token)
	return httpClient.Do(req)
}

func statusSuccess(res *http.Response) bool {
	return (res.StatusCode >= 200) || (res.StatusCode < 400)
}

// Read the body of the response to force the connection to close.
// Ref: https://stackoverflow.com/a/53589787
func readBody(readCloser io.ReadCloser) ([]byte, error) {
	defer readCloser.Close()
	body, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func joinURLPath(base string, elem ...string) string {
	u, err := url.Parse(base)
	if err != nil {
		return ""
	}
	u.Path = path.Join(append([]string{u.Path}, elem...)...)
	return u.String()
}
