package middle

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	models "sendMes/models"
	logger "sendMes/logs"

	"strings"
)

func SendMessage(receptor, message string) (result *models.Results, err error) {
	bodyMsg := &models.BodyMessage{}
	bodyMsg.Message = message
	bodyMsg.Sender = models.OwnSender
	bodyMsg.Receptor = receptor

	apiURL := os.Getenv("ApiUrl")

	formdata := url.Values{}
	formdata.Set("message", bodyMsg.Message)
	formdata.Set("receptor", bodyMsg.Receptor)
	formdata.Set("sender", bodyMsg.Sender)

	r, err := http.NewRequest("POST", apiURL, strings.NewReader(formdata.Encode()))
	if err != nil {
		logger.Gl.Error(
			"failed to create http request",
			"err", err,
			"api_url", apiURL,
			"step", "create_request",
		)
		return nil, err
	}

	header := os.Getenv("ApiKey")
	r.Header.Set("apikey", header)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		logger.Gl.Error(
			"failed to send request to api",
			"err", err,
			"api_url", apiURL,
			"step", "http_do",
		)
		return nil, err
	}
	defer resp.Body.Close()

	response := &models.Results{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		logger.Gl.Error(
			"failed to decode api response",
			"err", err,
			"status_code", resp.StatusCode,
			"step", "decode_response",
		)
		return nil, err
	}
	logger.Gl.Info(
		"message send request completed",
		"status_code", response.MessageIds,
		"receptor_suffix", bodyMsg.Receptor[len(bodyMsg.Receptor)-4:],
	)

	return response, nil
}
