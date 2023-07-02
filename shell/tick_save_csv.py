import datetime
import json

import pandas
import redis

redis_cli = redis.StrictRedis(host='xxx.xxx.xxx', port=6379, decode_responses=True, db=0)


def get_df(dlist: list) -> pandas.DataFrame:
    obj_list = []
    for json_str in dlist:
        obj = json.loads(json_str)
        obj_list.append(obj)
    return pandas.DataFrame(obj_list)


def save_csv(_df: pandas.DataFrame):
    _df.to_csv(f"./csv_db/{datetime.datetime.now().strftime('%Y%m%d')}.csv", index=False)


if __name__ == '__main__':
    key_list = redis_cli.keys("RT_*")

    df_list = []
    for key in key_list:
        data_list = redis_cli.lrange(key, 0, -1)
        df = get_df(data_list)
        df_list.append(df)

    df_all = pandas.concat(df_list)
    df_all = df_all.sort_values("trade_time")
    save_csv(df_all)
