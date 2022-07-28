package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/satyshef/registar/internal/sim5"
	"github.com/satyshef/registar/internal/sms365"
	"github.com/satyshef/registar/internal/sms3t"
	"github.com/satyshef/registar/internal/sms_acktiwator"
	"github.com/satyshef/tdbot"

	"time"

	tdc "github.com/satyshef/go-tdlib/client"
	"github.com/satyshef/go-tdlib/tdlib"
	"github.com/satyshef/tdbot/config"
	"github.com/satyshef/tdbot/profile"

	"github.com/valyala/fasthttp"
)

type SmsService interface {
	URLGetNumber() string
	URLCancelOrder(id string) string
	URLGetStatus(id string) string
	GetNumber(response string) (string, string, error)
	GetStatus(response string) (string, string)
}

var (
	bot         *tdbot.Bot
	smsServ     SmsService
	accCount    int
	profileDir  string
	configDir   string
	configName  string
	serviceName string
	serviceDir  string
	phoneNumber string
	orderID     string

	runProcess   bool
	runWaiteCode bool
	successCount int
)

func init() {
	flag.StringVar(&profileDir, "p", "./profiles", "Путь к директории с профилями")
	flag.StringVar(&configDir, "d", "./data/config", "Путь к директории конфигураций")
	flag.StringVar(&serviceDir, "r", "./data/service", "Путь к дериктории с сервисами\n\n")
	flag.StringVar(&phoneNumber, "n", "", "Номер телефона. Если указан СМС сервис не используется")
	flag.StringVar(&configName, "c", "bot.toml", "Название конфигурационного файла")
	flag.IntVar(&accCount, "a", 1, "Количество аккаунтов")
	flag.StringVar(&serviceName, "s", "", "Имя СМС сервиса\n  sms365\n  sms3t\n 5sim\n sms-acktiwator\n\n")

}

func main() {

	flag.Parse()
	serviceDir = strings.Trim(serviceDir, "/")
	configDir = strings.Trim(configDir, "/")

	switch serviceName {
	case "365sms":
		smsServ = sms365.New(serviceDir + "/365sms.toml")
	case "sms3t":
		smsServ = sms3t.New(serviceDir + "/sms3t.toml")
	case "5sim":
		smsServ = sim5.New(serviceDir + "/5sim.toml")
	case "sms-acktiwator":
		smsServ = sms_acktiwator.New(serviceDir + "/sms-acktiwator.toml")
	default:
		if phoneNumber == "" {
			fmt.Println("No required parameter. Please set phone number or SMS service")
			return
		}
	}

	if phoneNumber != "" {
		//manual registartion
		fmt.Println("Manual register", phoneNumber)
		manualRegistration()

	} else {
		//auto registrator
		fmt.Printf("SMS Service : %s\n\n", serviceName)
		profile.AddTail(&profileDir)
		autoRegistration()
	}

}

// Manual registartion
func manualRegistration() {
	// создаем профиль
	conf, err := config.Load(configDir + "/" + configName)
	if err != nil {
		fmt.Printf("Load config error : %s\n", err)
		return
	}

	prof, err := profile.New(phoneNumber, profileDir, conf)
	if err != nil {
		fmt.Printf("Load profile error : %s\n", err)
		return
	}

	bot = tdbot.New(prof)
	go func() {
		e := bot.Start()
		if e != nil {
			fmt.Printf("Start bot error : %s\n", err)
			return
		}
	}()

	for {
		time.Sleep(time.Second * 1)
		currentState, _ := bot.Client.Authorize()
		switch currentState.GetAuthorizationStateEnum() {

		case tdlib.AuthorizationStateWaitPhoneNumberType:
			continue
			/*
				fmt.Print("Enter phone: ")
				fmt.Scanln(&phoneNumber)
				_, err := bot.Client.SendPhoneNumber(phoneNumber)
				if err != nil {
					fmt.Printf("Error sending phone number: %v", err)
				}
			*/

		case tdlib.AuthorizationStateWaitCodeType:
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := bot.Client.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %s\n", err.(*tdlib.Error).Message)
			}
			continue
		case tdlib.AuthorizationStateWaitPasswordType:
			fmt.Print("Enter Password: ")
			var password string
			fmt.Scanln(&password)
			_, err := bot.Client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %s\n", err.(*tdlib.Error).Message)
			}
			continue
		case tdlib.AuthorizationStateReadyType:
			fmt.Println("Registration completed")
			goto Exit
		default:
			fmt.Printf("%s\n", currentState)
			continue
		}
	Exit:
		break
	}

}

// Auto registration
func autoRegistration() {
	var err error

	for successCount < accCount {
		runWaiteCode = false
		phoneNumber, orderID, err = getNumber()
		/*
			fmt.Printf("%s - %s ", phoneNumber, orderID)
			os.Exit(1)
		*/
		if err != nil {
			if err.Error() != "NO_NUMBERS" {
				fmt.Println("Get number : ", err)
			}

			if err.Error() == "NO_BALANCE" {
				//exit from program
				return
			}

			time.Sleep(time.Second * 3)
		} else {
			fmt.Println("#", successCount+1)
			err := startReg(phoneNumber)
			if err != nil {
				fmt.Println("Registar Error : ", err)
				time.Sleep(time.Second * 3)
			} else {
				//fmt.Println("OK")
			}
		}
	}
}

func startReg(phone string) error {
	fmt.Println("Start registration", phone)
	bot = nil
	runProcess = true
	//dir := profileDir
	//dir += phone
	//profile.AddTail(&dir)

	//прверяем существование профиля, если err == nil значит профиль существует
	/*
		if profile.IsProfile(dir) {
			cancelOrder(orderID)
			return fmt.Errorf("Профиль существует")
		}
	*/
	/*
		prof, err := profile.Get(dir)
		if err == nil {
			cancelOrder(orderID)
			prof.Close()
			return fmt.Errorf("Профиль существует")
		}

		if err.Error() != fmt.Sprintf("%s does not exist", dir) {
			cancelOrder(orderID)
			return err
		}
	*/

	// создаем профиль
	conf, err := config.Load(configDir + "/" + configName)
	if err != nil {
		cancelOrder(orderID)
		return err
	}

	prof, err := profile.New(phone, profileDir, conf)
	if err != nil {
		cancelOrder(orderID)
		return err
	}
	defer prof.Close()

	//bot.Client.AddEventHandler(eventCatcher)

	go func() {
		for {
			// Init bot
			bot = tdbot.New(prof)
			e := bot.Start()
			if e == nil {
				return
			}
			// реакции на ошибки во время запуска бота
			switch e.Code {
			//Если таймаут делаем еще одну попытку
			case tdc.ErrorCodeTimeout:
				bot.Logger.Errorf("TIMEOUT : %#v\n", e)
				bot.Stop()
				time.Sleep(time.Second * 1)
				continue
			case tdc.ErrorCodeStopped:
				bot.Logger.Errorf("CLIENT ERROR : %#v\n", e)
				time.Sleep(time.Second * 1)
				continue
			case profile.ErrorCodeDirNotExists:
				bot.Logger.Errorln("Dir not exists")
				runProcess = false
				goto Exit
			case profile.ErrorCodeLimitExceeded,
				tdc.ErrorCodePhoneBanned,
				tdc.ErrorCodePhoneInvalid,
				tdc.ErrorCodeAborted,
				tdc.ErrorCodeManyRequests:
				fmt.Println("CANCEL : ", e)
				cancelReg(orderID)
				goto Exit

			default:
				bot.Logger.Errorf("START ERROR: %#v\n", e)
				goto Exit
			}

		}
	Exit:
	}()

	for runProcess {
		time.Sleep(time.Second * 1)
		currentState, _ := bot.Client.Authorize()
		if currentState == nil {
			continue
		}

		switch currentState.GetAuthorizationStateEnum() {
		case tdlib.AuthorizationStateReadyType:
			if bot.Status == tdbot.StatusReady {
				successCount++
				fmt.Println("Completed")
				runProcess = false
			}

		case tdlib.AuthorizationStateWaitCodeType:
			//Если код отправлен в телеграм тогда прекращаем регистрацию
			s := currentState.(*tdlib.AuthorizationStateWaitCode)
			if s.CodeInfo.Type != nil && s.CodeInfo.Type.GetAuthenticationCodeTypeEnum() == tdlib.AuthenticationCodeTypeTelegramMessageType {
				bot.Logger.Infoln("Dont wait telegram code")
				cancelReg(orderID)
			} else {
				startWaiteCode()
			}

		case tdlib.AuthorizationStateWaitPasswordType:
			bot.Logger.Infoln("In account set password")
			cancelReg(orderID)

		default:
			//fmt.Println("Unknown client state : ", currentState.GetAuthorizationStateEnum())
			time.Sleep(time.Second * 3)
		}

	}
	fmt.Println("Success")

	return nil

}

func cancelReg(id string) {
	cancelOrder(id)
	bot.Stop()
	err := bot.Profile.Remove()
	if err != nil {
		fmt.Println("Remove profile error : ", err)
	}
	runProcess = false
}

// запуск процесса получения кода из смс
func startWaiteCode() {

	if !runWaiteCode {
		runWaiteCode = true
		waitTimeout := 60
		sleepTimeout := 2
		go func() {
			for waitTimeout > 0 {
				status, data := getStatus(orderID)
				//fmt.Println(status)
				switch status {
				case "UNKNOWN":
					fmt.Println("Неизвестный статус активации")
					time.Sleep(time.Second * 2)
					continue

				case "NO_ACTIVATION":
					fmt.Println("На СМС сервере нет данной активации")
					err := bot.Profile.Remove()
					if err != nil {
						fmt.Println(err)
					}
					runProcess = false
					bot.Stop()
					return

				case "WAIT_CODE":
					time.Sleep(time.Duration(sleepTimeout) * time.Second)
					waitTimeout -= sleepTimeout
					continue

				case "OK":
					fmt.Println("Code : ", data)
					bot.SendCode(data)
					return

				default:
					fmt.Println("Unknown status :", status)
				}

			}
			fmt.Println("Wait code timeout!!!!")
			cancelReg(orderID)
		}()

	}
}

func getNumber() (string, string, error) {
	urlGetNumber := smsServ.URLGetNumber()
	response := sendGet(urlGetNumber, 10000)
	//response := "ACCESS_NUMBER:234242:79663345925"
	//response := fmt.Sprintf("ACCESS_NUMBER:234242:%d", time.Now().Unix())
	return smsServ.GetNumber(response)
}

// послать запрос на смс сервер для отмены активации
func cancelOrder(id string) {
	urlCancel := smsServ.URLCancelOrder(id)
	response := sendGet(urlCancel, 10000)
	if response == "" {
		fmt.Println("Ошибка возврата номера")
	} else {
		fmt.Println("Cancel responce : ", response)
	}
}

// статус заявки
func getStatus(id string) (string, string) {
	urlStatus := smsServ.URLGetStatus(id)
	response := sendGet(urlStatus, 10000)
	return smsServ.GetStatus(response)

}

// timeout - время ожидания ответа в миллисекундах
func sendGet(url string, timeout int) string {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	/*
		req.Header.Add("Expires", time.Unix(0, 0).Format(time.RFC1123))
		req.Header.Add("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
		req.Header.Add("Pragma", "no-cache")
		req.Header.Add("X-Accel-Expires", fmt.Sprintf("%d", timeout))
		req.Header.Add("Accept", "*/ /*")
	 */

	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	var err error
	if timeout != 0 {
		err = fasthttp.DoTimeout(req, resp, time.Millisecond*time.Duration(timeout))
	} else {
		err = fasthttp.Do(req, resp)
	}

	if err != nil {
		fmt.Printf("Client get failed: %s\n", err)
		return ""
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		fmt.Printf("Expected status code %d but got %d\n", fasthttp.StatusOK, resp.StatusCode())
		return ""
	}

	body := resp.Body()

	//fmt.Printf("Response body is: %s\n", body)
	return string(body)
}
