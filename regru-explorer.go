package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

const (
  url = "https://api.cloudvps.reg.ru/v1/reglets/"
)


//Получение id ВМ
func getVMId(name, token string) string {
  req, err := http.NewRequest("GET", url)
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
  defer res.Body.Close()
  ///Здесь парсить ответ, возвращать id по имени
  // пример ответа:
  /* "reglets": [
    {
      "backups_enabled": false,
      "billed_until": "2023-02-19 12:43:12",
      "created_at": "2023-02-18 07:02:54",
      "disk": 80,
      "disk_usage": 0.0,
      "external_application": null,
      "hostname": "95-163-236-226.cloudvps.regruhosting.ru",
      "id": 2623965,
      "image": {
        "created_at": "2020-04-07 11:29:07",
        "distribution": "ubuntu-20.04",
        "id": 306495,
        "min_disk_size": "5",
        "name": "Ubuntu 20.04 LTS",
        "private": false,
        "region_slug": "msk1",
        "size_gigabytes": "2.4",
        "slug": "ubuntu-20-04-amd64",
        "type": "distribution"
      },
      "image_id": 306495,
      "ip": "95.163.236.226",
      "ipv6": "2a00:f940:2:4:2::5808",
      "last_backup_date": null,
      "locked": 0,
      "memory": 4096,
      "name": "GITLAB",
      "ptr": "95-163-236-226.cloudvps.regruhosting.ru",
      "region_slug": "msk1",
      "service_id": 68156613,
      "size": {
        "archived": 0,
        "disk": 80,
        "id": 1123,
        "memory": 4096,
        "name": "Base-4",
        "price": "2.6",
        "price_month": 1750,
        "slug": "base-4",
        "unit": "hour",
        "vcpus": 4,
        "weight": 40
      },
      "size_slug": "base-4",
      "status": "off",
      "sub_status": null,
      "vcpus": 4,
      "vpcs": []
    },
*/
  return ""
}

// Запуск ВМ
func startVM(id, token string) *http.Response {
  cur_url := url + id + "/actions"
  body := []byte(`{
    "type" : "start"
  }`)
  return execActionsReq(cur_url, token, body)
}

// Остановка ВМ
func stopVM(id, token string) *http.Response {
  cur_url := url + id + "/actions"
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


func main() {
  token := os.Getenv("TOKEN") //Получаем токен
  name := os.Args[1] //Имя машины
  id := getVMId(name)
  if id == "" {
    panic("Undefined")
  }
  arg := os.Args[2] // Операция

  switch arg {
    case "start":
      fmt.Println("Запуск контейнера")
      startVM(id, token)
    case "stop":
      fmt.Println("Запуск контейнера")
      stopVM(id, token)
  }
}
