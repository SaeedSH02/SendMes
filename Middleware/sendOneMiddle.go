package middle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	models "sendMes/Models"
	log "sendMes/logs"
	"strings"

	"go.uber.org/zap"
)

func SendOneSMS(receptor string) (result *models.Results,err error) {
	bodyMsg := &models.BodyMessage{}
	bodyMsg.Message = models.OrginalMsg
	bodyMsg.Sender = models.OwnSender
	bodyMsg.Receptor = receptor

	ApiUrl := os.Getenv("ApiUrl")

	formdata := url.Values{}
	formdata.Set("message", bodyMsg.Message)
	formdata.Set("receptor", bodyMsg.Receptor)
	formdata.Set("sender", bodyMsg.Sender)


	r, err := http.NewRequest("POST",ApiUrl, strings.NewReader(formdata.Encode()))
	if err != nil {
		log.Gl.Error("Can't Creating Request", zap.Error(err))
		return nil, err
	}

	header := os.Getenv("ApiKey")
	r.Header.Set("apikey", header)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")


	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Gl.Error("Can't send Req to Api", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()



	response := &models.Results{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		log.Gl.Error("Can't decode response", zap.Error(err))
		return nil, err
	}
	fmt.Printf("Send message to %s is %v", bodyMsg.Receptor, response)

	return response, nil
}


