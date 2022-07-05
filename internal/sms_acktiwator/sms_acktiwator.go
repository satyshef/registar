package sms_acktiwator

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/satyshef/registar/internal/sms_acktiwator/config"
)

type Service struct {
	APIKey  string
	ID      string
	Country string
}

//const API_KEY = "8nHsK3IhRHCHfJQoVVyyRkHZppAveD"

func New(configFile string) *Service {
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	return &Service{
		APIKey:  conf.APIKey,
		ID:      conf.ID,
		Country: conf.Country,
	}
}

func (s *Service) URLGetNumber() string {
	//urlGetNumberStatus := fmt.Sprintf("https://365sms.ru/stubs/handler_api.php?api_key=%s&action=getNumbersStatus&country=0&operator=tg_0", API_KEY)
	//urkGetBalance := fmt.Sprintf("https://365sms.ru/stubs/handler_api.php?api_key=%s&action=getBalance", API_KEY)

	return fmt.Sprintf("https://sms-acktiwator.ru/api/getnumber/%s?id=%s&code=%s", s.APIKey, s.ID, s.Country)
	//return fmt.Sprintf("https://365sms.ru/stubs/handler_api.php?api_key=%s&action=getNumber&service=%s&country=%s", s.APIKey, s.ID, s.Country)

}

// послать запрос на смс сервер для отмены активации
func (s *Service) URLCancelOrder(id string) string {
	//return fmt.Sprintf("https://365sms.ru/stubs/handler_api.php?api_key=%s&action=setStatus&status=8&id=%s", s.APIKey, id)

	return fmt.Sprintf("https://sms-acktiwator.ru/api/setstatus/%s?id=%s&status=1", s.APIKey, id)
}

func (s *Service) URLGetStatus(id string) string {
	//return fmt.Sprintf("https://365sms.ru/stubs/handler_api.php?api_key=%s&action=getStatus&id=%s", s.APIKey, id)

	return fmt.Sprintf("https://sms-acktiwator.ru/api/getstatus/%s?id=%s", s.APIKey, id)
}

// Return  Number, ID, Error
func (s *Service) GetNumber(response string) (string, string, error) {
	if response == "" {
		return "", "", fmt.Errorf("%s", "Empty response")
	}
	if response == "Expected status code 200 but got 500" {
		return "", "", fmt.Errorf("%s", "Bad country")
	}
	var resp map[string]interface{}
	json.Unmarshal([]byte(response), &resp)
	if resp["number"] != nil {
		number := resp["number"].(string)
		number = strings.Trim(number, "+")
		id := fmt.Sprintf("%d", int64(resp["id"].(float64)))
		return number, id, nil
	}
	if resp["name"] == "error" {
		code := int(resp["code"].(float64))
		switch code {
		case 101:
			return "", "", fmt.Errorf("%s", "WRONG_SERVICE")
		case 102:
			return "", "", fmt.Errorf("%s", "NO_NUMBERS")
		case 103:
			return "", "", fmt.Errorf("%s", "NO_NUMBERS")
		}
	}
	return "", "", fmt.Errorf("WRONG RESPONSE : %s", response)
}

// статус заявки
func (s *Service) GetStatus(response string) (string, string) {

	if response == "" {
		return "UNKNOWN", ""
	}

	if response == "null" {
		return "WAIT_CODE", ""
	}

	var resp map[string]interface{}
	json.Unmarshal([]byte(response), &resp)
	if resp["number"] != nil {
		code := fmt.Sprintf("%s", resp["small"])
		return "OK", code
	}

	//fmt.Printf("STATUS %#v\n\n", response)

	return response, ""
}
