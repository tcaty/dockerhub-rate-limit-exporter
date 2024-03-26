package ratelimit

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/tcaty/dockerhub-rate-limit-exporter/pkg/utils"
)

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
