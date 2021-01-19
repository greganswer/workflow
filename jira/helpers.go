package jira

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"
)

// Config keys.
const (
	UsernameConfigKey = "jira.username"
	TokenConfigKey    = "jira.token"
	APIConfigKey      = "jira.api_url"
	WebConfigKey      = "jira.api_url"
)

// URLs.
const (
	APIInstructionsURL = "https://confluence.atlassian.com/cloud/api-tokens-938839638.html"
	APIIssuePath       = "/rest/api/3/issue"
	APIUserPath        = "/rest/api/3/user"
	WebIssuePath       = "/browse"
)

var httpClient *http.Client

const (
	maxIdleConnections int = 20
	requestTimeout     int = 5
)

func init() {
	// Enable line numbers in logging
	log.SetFlags(log.Ltime | log.Llongfile)

	httpClient = createHTTPClient()
}

// Config contains Jira configuration values.
type Config struct {
	Username  string
	AccountID string
	Token     string
	APIURL    string
	WebURL    string
}

// errorResponse is the data structure for an error response from Jira's JSON API.
type errorResponse struct {
	Messages []string `json:"errorMessages"`
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Timeout: time.Duration(requestTimeout) * time.Second,
	}

	return client
}

func makeRequest(method, u string, reqBody []byte, c *Config) (*http.Response, error) {
	req, err := http.NewRequest(method, u, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.Username, c.Token)
	return httpClient.Do(req)
}

func statusSuccess(res *http.Response) bool {
	return res.StatusCode >= 200 && res.StatusCode < 400
}

// Read the body of the response to force the connection to close.
// Ref: https://stackoverflow.com/a/53589787
func readBody(readCloser io.ReadCloser) ([]byte, error) {
	defer readCloser.Close()
	body, err := ioutil.ReadAll(readCloser)
	if err != nil {
		log.Fatalln(err)
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
