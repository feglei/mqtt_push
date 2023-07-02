import sys
import time

import pandas
import datetime
from paho.mqtt import client as mqtt_client

broker = 'xxx.xxx.xxx'
port = 1883
keepalive = 60
topic = "/python/mqtt"  # 消息主题
client_id = f'play_tick-{datetime.datetime.now().strftime("H%M%S")}'  # 客户端id不能重复


def connect_mqtt():
    def on_connect(client, userdata, flags, rc):
        # 响应状态码为0表示连接成功
        if rc == 0:
            print("Connected to MQTT OK!")
        else:
            print("Failed to connect, return code %d\n", rc)

    # 连接mqtt代理服务器，并获取连接引用
    client = mqtt_client.Client(client_id)
    client.on_connect = on_connect
    client.connect(broker, port, keepalive)
    return client


if __name__ == '__main__':
    if len(sys.argv) != 2:
        print("参数: 需要播放的tick csv文件路径")
        exit(1)

    file_path = sys.argv[1]
    print("回播文件:", file_path)
    df = pandas.read_csv(file_path, dtype={"symbol_id": str})
    df = df.sort_values("trade_time")
    print(df)

    client = connect_mqtt()
    for i, row in df.iterrows():
        ts_code = row["symbol_id"]
        json_str = row.to_json()

        topic = f"test/stock/{ts_code}"
        result = client.publish(topic, json_str)
        print(topic, json_str)
        status = result[0]
        if status != 0:
            print(f"Failed to send message to topic {topic}")
        # time.sleep(0.5)




