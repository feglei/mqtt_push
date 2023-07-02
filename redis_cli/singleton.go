package redis_cli

var RedisCli = Cache{}

func Init() {
	RedisCli.StartAndGC("{\"conn\":\"xxx.xxx.xxx:6379\", \"password\":\"\"}")
}
