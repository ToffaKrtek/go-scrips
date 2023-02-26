package main
// ИМЯ_КОНТЕЙНЕРА команда (start, stop)
// Пример вызова: regru-explorer BUILDER start 
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
  "time"
)

const (
  url = "https://api.cloudvps.reg.ru/v1/reglets"
  token = ""
  maxAwait = 80 //Максимальное время ожидания подтверждения операции в секундах
  stepToNewRequest = 4 // Шаг для повторных запросов сверки статуса при подтверждении операции в секундах
)

type Reglet struct {
  Id int `json:"id"`
  Name string `json:"name"`
  Status string `json:"status"`
}
type Body struct {
  Reglets []Reglet `json:"reglets"`
  Reglet Reglet `json:"reglet`
}

func call(url, token string) (Body, error) {
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
      return Body{}, err
  }
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Authorization", "Bearer " + token)
  response, err := http.DefaultClient.Do(req)
  if err != nil {
      return Body{}, err
  }
  defer response.Body.Close()

  if response.StatusCode == http.StatusOK {
    bodyBytes, err := ioutil.ReadAll(response.Body)
    if err != nil {
      return Body{}, err
    }
    /*bodyString := string(bodyBytes)
    fmt.Println(bodyString)*/
    var target Body
    _ = json.Unmarshal(bodyBytes, &target)
    return target, nil
  }
  fmt.Println("Ошибка. Статус ответа.")
  return Body{}, nil
}
//Получение id ВМ
func getVMId(name, token string) string {
  reglets, err := call(url, token)
  if err != nil || len(reglets.Reglets) < 1 {
    panic("Ошибка запроса")
  }
  for _, reglet := range reglets.Reglets {
    if reglet.Name == name {
      return strconv.Itoa(reglet.Id)
    }
  }
  return ""
}

// Запуск ВМ
func startVM(id, token string) *http.Response {
  cur_url := url + "/" + id + "/actions"
  body := []byte(`{
    "type" : "start"
  }`)
  return execActionsReq(cur_url, token, body)
}

// Остановка ВМ
func stopVM(id, token string) *http.Response {
  cur_url := url + "/" + id + "/actions"
  body := []byte(`{
    "type" : "stop"
  }`)
  return execActionsReq(cur_url, token, body)
}

func execActionsReq(cur_url, token string, body []byte) *http.Response {
  req, err := http.NewRequest("POST", cur_url, bytes.NewBuffer(body))
  if err != nil {
    panic(err)
  }
  req.Header.Set("Authorization", "Bearer " + token)
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  res, err := client.Do(req)
  if err != nil {
    panic(err)
  }
  return res
}

func awaitStatusChange(id, token, status string) bool {
  timeFromStart := 0
	for range time.Tick(time.Second * stepToNewRequest) {
		// do the interval task
    cur_status := getStatusVM(id, token)
    if cur_status == status {
      fmt.Println("Готово!")
      return true
    }
    timeFromStart += stepToNewRequest
    if timeFromStart >= maxAwait {
      fmt.Println("Превышено время ожидания при подтверждения операции!")
      return false
    }
  }
  fmt.Println("Ошибка при ожидании подтверждения операции!")
  return false
}

func getStatusVM(id, token string) string {
  cur_url := url + "/" + id
  body,_ := call(cur_url, token)
  return body.Reglet.Status
}

func main() {
  //token := os.Getenv("TOKEN") //Получаем токен
  if len(os.Args) < 2 {
    panic("Не указано имя ВМ")
  }
  name := os.Args[1] //Имя машины
  if name == "" {
    panic("Не указано имя ВМ")
  }
  id := getVMId(name, token)
  if id == "" {
    panic("Не найдена ВМ")
  }
  fmt.Println("ID-сервера -- ", id)

  if len(os.Args) < 3 {
    panic("Не указана команда ВМ")
  }
  arg := os.Args[2] // Операция

  switch arg {
    case "start":
      fmt.Println("Запуск контейнера")
      startVM(id, token)
      awaitStatusChange(id, token, "active")
      return 
    case "stop":
      fmt.Println("Остановка контейнера")
      stopVM(id, token)
      awaitStatusChange(id, token, "off")
      return 
  }
}
