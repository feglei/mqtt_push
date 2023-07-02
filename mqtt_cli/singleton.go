package mqtt_cli

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"mqtt_push/conf"
	"os"
	"time"
)

type MQTTCli struct {
	Broker string
	Port   int

	ClientID string
	UserName string
	Password string

	EClient            mqtt.Client
	Opts               *mqtt.ClientOptions
	messagePubHandler  mqtt.MessageHandler
	connectHandler     mqtt.OnConnectHandler
	connectLostHandler mqtt.ConnectionLostHandler
}

var ECli = MQTTCli{
	Broker: "xxx.xxx.xxx",
	Port:     1883,
	ClientID: "rt_tick_",
	UserName: "sys_rttick",
	Password: "xxxxxxxxxx",
}

func (m *MQTTCli) Connect() {
	m.Opts = mqtt.NewClientOptions()
	m.Opts.AddBroker(fmt.Sprintf("tcp://%s:%d", m.Broker, m.Port))
	m.Opts.SetClientID(m.ClientID)
	fmt.Println(m.ClientID)
	m.Opts.SetUsername(m.UserName)
	m.Opts.SetPassword(m.Password)
	m.Opts.SetDefaultPublishHandler(messagePubHandler)
	m.Opts.OnConnect = connectHandler
	m.Opts.OnConnectionLost = connectLostHandler
	m.EClient = mqtt.NewClient(m.Opts)
	if token := m.EClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
		os.Exit(-1)
	}
}

func (m *MQTTCli) Subscribe(topic string) {
	// 订阅主题
	if token := m.EClient.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func (m *MQTTCli) Publish(topic string, message string) {
	// 发布消息
	token := m.EClient.Publish(topic, 0, false, message)
	//fmt.Println(token, topic, message)
	token.Wait()
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//fmt.Println(msg.Topic(), string(msg.Payload()))
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func Init() {
	ECli.ClientID = fmt.Sprintf("rt_tick_%d_%d", conf.SECTION, conf.INDEX)
	ECli.Connect()
	go loop_state()
}

func loop_state() {
	for {
		state := getStateStr()
		now_time := time.Now()
		message := fmt.Sprintf("%s %s", state, now_time.Format("2006.01.02 15:04:05"))
		ECli.Publish("robot/state", message)
		time.Sleep(500 * time.Millisecond)
	}
}

func getStateStr() string {

	var now_t1 = time.Now()
	var AmStartStr = fmt.Sprintf("%d%02d%02d%02d%02d%02d", now_t1.Year(), now_t1.Month(), now_t1.Day(), 9, 00, 00)
	var AmClosetStr = fmt.Sprintf("%d%02d%02d%02d%02d%02d", now_t1.Year(), now_t1.Month(), now_t1.Day(), 11, 30, 00)
	var PmStartStr = fmt.Sprintf("%d%02d%02d%02d%02d%02d", now_t1.Year(), now_t1.Month(), now_t1.Day(), 13, 00, 00)
	var PmCloseStr = fmt.Sprintf("%d%02d%02d%02d%02d%02d", now_t1.Year(), now_t1.Month(), now_t1.Day(), 15, 00, 00)
	var AmStart, _ = time.ParseInLocation("20060102150405", AmStartStr, time.Local)
	var AmCloset, _ = time.ParseInLocation("20060102150405", AmClosetStr, time.Local)
	var PmStart, _ = time.ParseInLocation("20060102150405", PmStartStr, time.Local)
	var PmClose, _ = time.ParseInLocation("20060102150405", PmCloseStr, time.Local)

	isAm := AmStart.Unix() <= now_t1.Unix() && now_t1.Unix() <= AmCloset.Unix()
	isPm := PmStart.Unix() <= now_t1.Unix() && now_t1.Unix() <= PmClose.Unix()

	if now_t1.Unix() < AmStart.Unix() {
		return "未开盘"
	} else if isAm || isPm {
		return "开盘"
	} else if now_t1.Unix() > PmClose.Unix() {
		return "收盘"
	} else {
		return "休市"
	}
}
