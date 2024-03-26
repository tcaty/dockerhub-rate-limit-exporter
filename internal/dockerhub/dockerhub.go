package dockerhub

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DockerHub struct {
	repository string
}

func New(repository string) *DockerHub {
	return &DockerHub{
		repository: repository,
	}
}

func (dh *DockerHub) FetchHeaders() (http.Header, error) {
	token, err := dh.fetchToken()
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

	return res.Header, err
}

func (dh *DockerHub) fetchToken() (string, error) {
	url := fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull", dh.repository)
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var tokenResponseBody TokenResponseBody
	d := json.NewDecoder(res.Body)
	if err := d.Decode(&tokenResponseBody); err != nil {
		return "", err
	}

	return tokenResponseBody.Token, nil
}
