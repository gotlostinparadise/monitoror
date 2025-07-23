package repository

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"github.com/monitoror/monitoror/monitorables/http/api"
	"github.com/monitoror/monitoror/monitorables/http/api/models"
	"github.com/monitoror/monitoror/monitorables/http/config"
)

type (
	httpRepository struct {
		verifyClient *http.Client
		skipClient   *http.Client
		config       *config.HTTP
	}
)

func NewHTTPRepository(config *config.HTTP) api.Repository {
	var certificates []tls.Certificate

	if config.Certificate != "" && config.Key != "" {
		cert, error := tls.LoadX509KeyPair(config.Certificate, config.Key)

		if error == nil {
			certificates = append(certificates, cert)
		}
	}

	trVerify := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false, Certificates: certificates}}
	trSkip := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true, Certificates: certificates}}

	verifyClient := &http.Client{Transport: trVerify, Timeout: time.Duration(config.Timeout) * time.Millisecond}
	skipClient := &http.Client{Transport: trSkip, Timeout: time.Duration(config.Timeout) * time.Millisecond}

	return &httpRepository{verifyClient: verifyClient, skipClient: skipClient, config: config}
}

func (r *httpRepository) Get(url string, headers map[string]string, sslVerify *bool) (response *models.Response, err error) {
	useVerify := r.config.SSLVerify
	if sslVerify != nil {
		useVerify = *sslVerify
	}

	client := r.verifyClient
	if !useVerify {
		client = r.skipClient
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	response = &models.Response{
		StatusCode: resp.StatusCode,
		Body:       bytes,
	}

	return
}
