package url_conf

import (
	"mqtt_push/conf"
	"mqtt_push/http_cli"
	"time"
)

func IsOpenDay() bool {
	isStr := http_cli.GetString(conf.DB_STOCKS_ISOPENDAY)
	if isStr == "1" {
		return true
	}
	if isStr == "0" {
		return false
	}
	time.Sleep(1 * time.Second)
	return IsOpenDay()
}
