### mqtt_push 主要获取tick数据，后将数据推送到emqtt服务节点。
#### 在 rt_tick 包里面计算了3分钟涨速 和 1分钟量等数据。



#### 接收数据示例代码：
call_back.go（接收mqtt消息后，分发到需要使用数据的模块）：

```
var ChanList = make([]chan StockMod, 0)
func RegisterChan(c chan StockMod) {
	ChanList = append(ChanList, c)
}
func Allocate(topic string, payload []byte) {
	//fmt.Println(topic, payload)
	if strings.HasPrefix(topic, "stock/") == false {
		return
	}
	sObj := StockMod{}
	err := json.Unmarshal(payload, &sObj)
	if err != nil {
		fmt.Println(err)
		return
	}
	if sObj.SymbolId == "" {
		fmt.Println("symbol_id is empty")
		return
	}
	//发送数据到所有的 channel
	for c := range ChanList {
		ChanList[c] <- sObj
	}
}
```

#### 业务代码：
```
var Chan = make(chan mqtt_cli.StockMod, 100)
func Init() {
  //注册一个 chan 给mqtt管理器。
	mqtt_cli.RegisterChan(Chan)
	go loop()
}
func loop() {
	for {
		select {
		case s := <-Chan:
      //业务逻辑代码段。。。
  }
}
```
