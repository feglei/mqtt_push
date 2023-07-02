package rt_mod

import (
	"encoding/json"
	"errors"
	"fmt"
	"mqtt_push/log"
	"mqtt_push/mqtt_cli"
	"mqtt_push/redis_cli"
	"strconv"
	"strings"
	"time"
)

type RtTickMod struct {
	SymbolID     string  `json:"symbol_id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	TradeTime    string  `json:"trade_time"`
	PrePrice     float64 `json:"pre_price"`
	Price        float64 `json:"price"`
	OpenPrice    float64 `json:"open_price"`
	HighPrice    float64 `json:"high_price"`
	LowPrice     float64 `json:"low_price"`
	ClosePrice   float64 `json:"close_price"`
	OpenInterest float64 `json:"open_interest"`
	Volume       float64 `json:"volume"`
	Amount       float64 `json:"amount"`
	Num          float64 `json:"num"`
	AskPrice1    float64 `json:"ask_price1"`
	AskVolume1   float64 `json:"ask_volume1"`
	BidPrice1    float64 `json:"bid_price1"`
	BidVolume1   float64 `json:"bid_volume1"`
	AskPrice2    float64 `json:"ask_price2"`
	AskVolume2   float64 `json:"ask_volume2"`
	BidPrice2    float64 `json:"bid_price2"`
	BidVolume2   float64 `json:"bid_volume2"`
	AskPrice3    float64 `json:"ask_price3"`
	AskVolume3   float64 `json:"ask_volume3"`
	BidPrice3    float64 `json:"bid_price3"`
	BidVolume3   float64 `json:"bid_volume3"`
	AskPrice4    float64 `json:"ask_price4"`
	AskVolume4   float64 `json:"ask_volume4"`
	BidPrice4    float64 `json:"bid_price4"`
	BidVolume4   float64 `json:"bid_volume4"`
	AskPrice5    float64 `json:"ask_price5"`
	AskVolume5   float64 `json:"ask_volume5"`
	BidPrice5    float64 `json:"bid_price5"`
	BidVolume5   float64 `json:"bid_volume5"`

	SpeedUp3   float64 `json:"speed_up3"`
	Amount1    float64 `json:"amount_1"`
	AmountAvg1 float64 `json:"amount_avg_1"`
	SaveTime   int64   `json:"save_time"`
}

func (r *RtTickMod) ToJsonStr() (jsonStr string, err error) {
	byte_data, err := json.Marshal(r)
	return string(byte_data), err
}

func (r *RtTickMod) Logic() error {
	llen, err := r.GetListLen(r.SymbolID)
	if err != nil {
		log.Error("GetListLen is error", err)
	}

	if llen <= 0 {
		r.SaveDB()
		return nil
	}
	nowMod := GetLNowRt(r.SymbolID)
	if Time2Int(r.TradeTime) <= Time2Int(nowMod.TradeTime) {
		return nil
	}

	var pre3Mod RtTickMod
	if llen >= 66 {
		pre3Mod = GetIndexRt(r.SymbolID, 65)
	} else {
		pre3Mod = GetLLastRt(r.SymbolID)
		pre3Mod.Price = nowMod.OpenPrice
	}

	if pre3Mod.Price > 0 {
		//计算三分钟涨速
		r.SpeedUp3 = (r.Price - pre3Mod.Price) / pre3Mod.Price
		//计算3分钟平均量
		r.AmountAvg1 = (r.Amount - pre3Mod.Amount) / 3
	}

	//计算一分钟量
	var pre1Mod RtTickMod
	if llen >= 22 {
		pre1Mod = GetIndexRt(r.SymbolID, 21)
		r.Amount1 = r.Amount - pre1Mod.Amount
	} else {
		r.Amount1 = r.Amount - pre3Mod.Amount
	}

	//保存任务
	r.SaveDB()
	llen, err = r.GetListLen(r.SymbolID)
	return nil
}

func (r *RtTickMod) SaveDB() (err error) {
	r.Amount1 = Decimal(r.Amount1)
	r.AmountAvg1 = Decimal(r.AmountAvg1)
	r.SpeedUp3 = Decimal(r.SpeedUp3 * 100)

	v, err := r.ToJsonStr()
	redis_cli.RedisCli.LPUSH(r.SymbolID, v)
	r.Push(r.SymbolID, v)
	return err
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func (r *RtTickMod) Push(topic string, message string) {
	mqtt_cli.ECli.Publish("stock/"+topic, message)
}

func (r *RtTickMod) GetListLen(key string) (int64, error) {
	d := redis_cli.RedisCli.LLEN(r.SymbolID)
	if d == nil {
		return 0, errors.New("Logic, llen is nil:" + r.SymbolID)
	}
	return d.(int64), nil
}

func Time2Int(t1 string) int {
	t1 = strings.Replace(t1, ":", "", -1)
	tNum, err := strconv.Atoi(t1)
	if err != nil {
		return -1
	}
	return tNum
}

func GetLLastRt(key string) (rt_mod RtTickMod) {
	return GetIndexRt(key, -1)
}

func GetLNowRt(key string) (rt_mod RtTickMod) {
	return GetIndexRt(key, 0)
}

func GetIndexRt(key string, index int) (rt_mod RtTickMod) {
	lastObj := redis_cli.RedisCli.LINDEX(key, index)
	if lastObj == nil {
		fmt.Println("lastObj is null ", key, index)
		return rt_mod
	}
	lastBytes := lastObj.([]byte)
	json.Unmarshal(lastBytes, &rt_mod)
	return rt_mod
}

func CreateRtTickMod(keys []string, item []string) (rt_mod RtTickMod) {
	rt_mod = RtTickMod{
		Code:         item[0],
		Name:         item[1],
		TradeTime:    item[2],
		PrePrice:     Str2Float64(item[3]),
		Price:        Str2Float64(item[4]),
		OpenPrice:    Str2Float64(item[5]),
		HighPrice:    Str2Float64(item[6]),
		LowPrice:     Str2Float64(item[7]),
		ClosePrice:   Str2Float64(item[8]),
		OpenInterest: Str2Float64(item[9]),
		Volume:       Str2Float64(item[10]),
		Amount:       Str2Float64(item[11]),
		Num:          Str2Float64(item[12]),
		AskPrice1:    Str2Float64(item[13]),
		AskVolume1:   Str2Float64(item[14]),
		BidPrice1:    Str2Float64(item[15]),
		BidVolume1:   Str2Float64(item[16]),
		AskPrice2:    Str2Float64(item[17]),
		AskVolume2:   Str2Float64(item[18]),
		BidPrice2:    Str2Float64(item[19]),
		BidVolume2:   Str2Float64(item[20]),
		AskPrice3:    Str2Float64(item[21]),
		AskVolume3:   Str2Float64(item[22]),
		BidPrice3:    Str2Float64(item[23]),
		BidVolume3:   Str2Float64(item[24]),
		AskPrice4:    Str2Float64(item[25]),
		AskVolume4:   Str2Float64(item[26]),
		BidPrice4:    Str2Float64(item[27]),
		BidVolume4:   Str2Float64(item[28]),
		AskPrice5:    Str2Float64(item[29]),
		AskVolume5:   Str2Float64(item[30]),
		BidPrice5:    Str2Float64(item[31]),
		BidVolume5:   Str2Float64(item[32]),
	}
	rt_mod.SymbolID = rt_mod.Code[:len(rt_mod.Code)-3]
	rt_mod.SaveTime = time.Now().Unix()
	return rt_mod
}

func Str2Float64(s string) (f float64) {
	var err error
	f, err = strconv.ParseFloat(s, 64)
	if err != nil {
		return float64(0)
	}
	return f
}
