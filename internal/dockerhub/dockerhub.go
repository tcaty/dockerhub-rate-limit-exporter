package dockerhub

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/exporter"
	"github.com/tcaty/dockerhub-rate-limit-exporter/pkg/utils"
)

type DockerHub struct {
	repository string
	username   string
	password   string
}

func New(repository string, username string, password string) *DockerHub {
	return &DockerHub{
		repository: repository,
		username:   username,
		password:   password,
	}
}

func (dh *DockerHub) FetchMetaData() (*exporter.MetaData, error) {
	authenticatedMode := false
	headers, err := dh.fetchHeaders(authenticatedMode)

	if err != nil {
		return nil, err
	}

	host, err := utils.ParseHeader(headers, "Docker-Ratelimit-Source")

	if err != nil {
		return nil, err
	}

	metaData := &exporter.MetaData{
		Host:     host,
		Username: dh.username,
	}

	return metaData, nil
}

func (dh *DockerHub) FetchRateLimitData() (*exporter.RateLimitData, error) {
	authenticatedMode := !(dh.username == AnonymousUsername && dh.password == AnonymousPassword)
	headers, err := dh.fetchHeaders(authenticatedMode)

	if err != nil {
		return nil, err
	}

	limit, err := parseRateLimitHeader(headers, "Limit")

	if err != nil {
		return nil, err
	}

	remaining, err := parseRateLimitHeader(headers, "Remaining")

	if err != nil {
		return nil, err
	}

	rateLimitData := &exporter.RateLimitData{
		Total:     limit,
		Remaining: remaining,
	}

	return rateLimitData, nil
}

func (dh *DockerHub) fetchHeaders(authenticatedMode bool) (http.Header, error) {
	token, err := dh.fetchToken(authenticatedMode)

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://registry-1.docker.io/v2/%s/manifests/latest", dh.repository)
	req, err := http.NewRequest(http.MethodHead, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could fetch dockerhub api, received status code: %d", res.StatusCode)
	}

	return res.Header, err
}

func (dh *DockerHub) fetchToken(authenticatedMode bool) (string, error) {
	url := fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull", dh.repository)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", err
	}

	if authenticatedMode {
		req.SetBasicAuth(dh.username, dh.password)
	}
	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could fetch dockerhub token, received status code: %d", res.StatusCode)
	}

	defer res.Body.Close()

	var tokenResponseBody TokenResponseBody
	d := json.NewDecoder(res.Body)
	if err := d.Decode(&tokenResponseBody); err != nil {
		return "", err
	}

	return tokenResponseBody.Token, nil
}
