package sms3t

import (
	"fmt"
	"log"
	"strings"

	"github.com/satyshef/registar/internal/sms3t/config"
)

type Service struct {
	APIKey  string
	ID      string
	Country string
}

//const API_KEY = "0f0e1da8d76bd061eb4ac70a966d6ac0"

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

	return fmt.Sprintf("https://sms3t.com/stubs/handler_api.php?api_key=%s&action=getNumber&service=%s&country=%s", s.APIKey, s.ID, s.Country)

}

// послать запрос на смс сервер для отмены активации
func (s *Service) URLCancelOrder(id string) string {
	return fmt.Sprintf("https://sms3t.com/stubs/handler_api.php?api_key=%s&action=setStatus&status=8&id=%s", s.APIKey, id)
}

func (s *Service) URLGetStatus(id string) string {
	return fmt.Sprintf("https://sms3t.com/stubs/handler_api.php?api_key=%s&action=getStatus&id=%s", s.APIKey, id)
}

func (s *Service) GetNumber(response string) (string, string, error) {

	//response := "ACCESS_NUMBER:234242:79993456789"
	switch response {
	case "":
		return "", "", fmt.Errorf("%s", "Empty response")
	case "NO_BALANCE":
		return "", "", fmt.Errorf("%s", "NO_BALANCE")
	case "NO_NUMBERS":
		return "", "", fmt.Errorf("%s", "NO_NUMBERS")
	case "WRONG_SERVICE":
		return "", "", fmt.Errorf("%s", "WRONG_SERVICE")
	default:
		if strings.Contains(response, "ACCESS_NUMBER") {
			s := strings.Split(response, ":")
			if len(s) == 3 {
				return s[2], s[1], nil
			}
		}
	}

	return "", "", fmt.Errorf("WRONG RESPONSE : %s", response)
}

// статус заявки
func (s *Service) GetStatus(response string) (string, string) {
	switch response {
	case "":
		return "UNKNOWN", ""

	case "NO_ACTIVATION":
		return "NO_ACTIVATION", ""

	case "STATUS_WAIT_CODE":
		return "WAIT_CODE", ""

	default:
		if strings.Contains(response, "STATUS_WAIT_RETRY") {
			s := strings.Split(response, ":")
			//time.Sleep(time.Second * 10)
			return "OK", s[1]
		}
	}

	return response, ""
}
