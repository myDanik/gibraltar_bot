package services

import (
	"bufio"
	"net/http"
	"time"
)

type ConfigService struct {
	APIUrl string
}

func NewConfigService(apiUrl string) *ConfigService {
	return &ConfigService{
		APIUrl: apiUrl,
	}
}

func (s *ConfigService) GetConfigs() (string, error) {
	responce, err := http.Get(s.APIUrl)
	if err != nil {
		return "", err
	}
	defer responce.Body.Close()
	if responce.StatusCode != http.StatusOK {
		return "", err
	}
	r := responce.Body
	scanner := bufio.NewScanner(r)
	var result string = ""
	for scanner.Scan() {
		result += scanner.Text() + "\n"

	}

	return result, nil
}

func (s *ConfigService) UpdateConfigs() error {
	client := http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodPatch, s.APIUrl, nil)
	if err != nil {
		return err
	}
	responce, err := client.Do(request)
	if err != nil {
		return err
	}
	defer responce.Body.Close()
	return nil

}
