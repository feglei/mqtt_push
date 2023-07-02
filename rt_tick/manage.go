package rt_tick

import (
	"mqtt_push/conf"
	"mqtt_push/http_cli"
	"mqtt_push/log"
	"mqtt_push/rt_tick/rt_mod"
)

func Init() {
	//初始股票ID对照表
	http_cli.Get(conf.DB_STOCKS_LIST, &rt_mod.StockMapKeys)
	log.Info("init stock map", len(rt_mod.StockMapKeys), " | ", rt_mod.StockMapKeys)
}

var loopNum = 0

func UpdateCode(section int, index int) {
	var temp_tscode string
	if section <= 0 {
		temp_tscode = rt_mod.StockAll()
	} else {
		temp_tscode = rt_mod.GetSection(section, index)
	}
	if len(temp_tscode) <= 9 {
		log.Error("tscode update is error.", temp_tscode)
		return
	}
	rt_mod.TsCode = temp_tscode
	log.Info("ts_code = ", rt_mod.TsCode)
}

func GetRtMod() rt_mod.Result {
	var repo = rt_mod.Repo{ApiName: conf.API_NAME_RTTICK, Token: conf.TOKEN, Params: rt_mod.ParamsMod{TsCode: rt_mod.TsCode}}
	var result rt_mod.Result
	http_cli.Post(conf.API_SERVER_URL, &repo, &result)
	return result
}
