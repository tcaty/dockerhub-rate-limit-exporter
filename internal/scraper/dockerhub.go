package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tcaty/dockerhub-rate-limit-exporter/cmd"
	"github.com/tcaty/dockerhub-rate-limit-exporter/pkg/utils"
)

type DockerHub struct {
	repository string
	username   string
	password   string
}

func NewDockerHub(flags cmd.Flags) *DockerHub {
	return &DockerHub{
		repository: flags.Repository,
		username:   flags.DockerHubUsername,
		password:   flags.DockerHubPassword,
	}
}

func (dh *DockerHub) FetchMetaData() (*MetaData, error) {
	headers, err := dh.fetchHeaders(false)

	if err != nil {
		return nil, err
	}

	host, err := utils.ParseHeader(headers, "Docker-Ratelimit-Source")

	if err != nil {
		return nil, err
	}

	var username string

	if dh.IsAuthenticatedMode() {
		username = dh.username
	} else {
		username = ""
	}

	metaData := &MetaData{
		Host:     host,
		Username: username,
	}

	return metaData, nil
}

func (dh *DockerHub) FetchRateLimitData() (*RateLimitData, error) {
	headers, err := dh.fetchHeaders(dh.IsAuthenticatedMode())

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

	rateLimitData := &RateLimitData{
		Total:     limit,
		Remaining: remaining,
	}

	return rateLimitData, nil
}

func parseRateLimitHeader(header http.Header, name string) (float64, error) {
	h, err := utils.ParseHeader(header, fmt.Sprintf("Ratelimit-%s", name))
	if err != nil {
		return 0, err
	}

	v, err := strconv.ParseFloat(strings.Split(h, ";")[0], 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func (dh *DockerHub) IsAuthenticatedMode() bool {
	return !(dh.username == "" && dh.password == "")
}

func (dh *DockerHub) fetchHeaders(IsAuthenticatedMode bool) (http.Header, error) {
	token, err := dh.fetchToken(IsAuthenticatedMode)

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

func (dh *DockerHub) fetchToken(IsAuthenticatedMode bool) (string, error) {
	url := fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull", dh.repository)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", err
	}

	if IsAuthenticatedMode {
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

type TokenResponseBody struct {
	Token       string    `json:"token"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	IssuedAt    time.Time `json:"issued_at"`
}
