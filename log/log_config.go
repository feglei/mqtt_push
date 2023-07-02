package log

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
)

type LogsConfig struct {

	//filename 保存的文件名
	FileName string `json:"filename"`

	//maxlines 每个文件保存的最大行数，默认值 1000000
	MaxLines int `json:"maxlines"`

	//maxsize 每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
	MaxSize int `json:"maxsize"`

	//daily 是否按照每天 logrotate，默认是 true
	Daily bool `json:"daily"`

	//maxdays 文件最多保存多少天，默认保存 7 天
	MaxDays int `json:"maxdays"`

	//rotate 是否开启 logrotate，默认是 true
	Rotate bool `json:"rotate"`

	//level 日志保存的时候的级别，默认是 Trace 级别
	Level int `json:"level"`
}

func Init(fileName string) {
	config := LogsConfig{
		FileName: "./log_file/" + fileName + ".log",
		MaxLines: 0,
		MaxSize:  0,
		Daily:    true,
		MaxDays:  7,
		Rotate:   true,
		Level:    7,
	}

	configStr, _ := json.Marshal(config)
	logs.SetLogger(logs.AdapterFile, string(configStr))
}

func Debug(f interface{}, v ...interface{}) {
	logs.Debug(f, v)
}

func Info(f interface{}, v ...interface{}) {
	logs.Info(f, v)
}

func Warn(f interface{}, v ...interface{}) {
	logs.Warn(f, v)
}

func Error(f interface{}, v ...interface{}) {
	logs.Error(f, v)
}

func Critical(f interface{}, v ...interface{}) {
	logs.Critical(f, v)
}
