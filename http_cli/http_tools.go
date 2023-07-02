package http_cli

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"mqtt_push/conf"
)

func GetString(url string) string {
	s, _ := req.Get(url)
	return s.String()
}

func Get(url string, result interface{}) {
	s, _ := req.Get(url)
	json.Unmarshal(s.Bytes(), result)
}

func Post(url string, repo interface{}, result interface{}) {
	req.MustGet(conf.API_SERVER_URL)
	client := req.C()
	resp, err := client.R().
		SetBody(&repo).
		SetResult(&result).
		Post(url)

	if err != nil {
		fmt.Println("error:", err, resp)
	}
}
