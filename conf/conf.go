package conf

import (
	"fmt"
	"time"
)

const (
	ALL = iota
	SYMBOL_ID
	OTHER_SYMBOL_IDS
)

var SECTION = 0
var INDEX = 0

func GetAppName() string {
	return fmt.Sprintf("all_%d_%d", SECTION, INDEX)
}

const (
	OPEN = iota
	REST
	CLOSE
)

func MarketState() int {

	var now_t1 = time.Now()
	var AmStartStr = fmt.Sprintf("%d%02d%02d%02d%02d%02d", now_t1.Year(), now_t1.Month(), now_t1.Day(), 9, 10, 00)
	var PmCloseStr = fmt.Sprintf("%d%02d%02d%02d%02d%02d", now_t1.Year(), now_t1.Month(), now_t1.Day(), 15, 10, 00)
	var AmStart, _ = time.ParseInLocation("20060102150405", AmStartStr, time.Local)
	var PmClose, _ = time.ParseInLocation("20060102150405", PmCloseStr, time.Local)

	if AmStart.Unix() <= now_t1.Unix() && now_t1.Unix() <= PmClose.Unix() {
		return OPEN
	} else if now_t1.Unix() > PmClose.Unix() {
		return CLOSE
	} else {
		return REST
	}
}
