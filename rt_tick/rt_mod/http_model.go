package rt_mod

type Repo struct {
	ApiName string    `json:"api_name"`
	Token   string    `json:"token"`
	Params  ParamsMod `json:"params"`
	Fields  string    `json:"fields"`
}
type ParamsMod struct {
	TsCode string `json:"ts_code"`
}

type Result struct {
	RequestId string  `json:"request_id"`
	Code      int     `json:"code"`
	Msg       string  `json:"msg"`
	Data      DataMod `json:"data"`
}
type DataMod struct {
	Fields  []string   `json:"fields"`
	Items   [][]string `json:"items"`
	HasMore bool       `json:"has_more"`
}
