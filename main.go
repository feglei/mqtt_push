package main

import (
	"flag"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"mqtt_push/conf"
	"mqtt_push/conf/url_conf"
	"mqtt_push/log"
	"mqtt_push/mqtt_cli"
	"mqtt_push/redis_cli"
	"mqtt_push/rt_tick"
	"mqtt_push/rt_tick/rt_mod"
	"os"
	"sync"
	"time"
)

func wrapper(rtObj *rt_mod.RtTickMod, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()
		rtObj.Logic()
		//fmt.Println(e)
	}
}

func RunTask() {
	var wg sync.WaitGroup

	now_time := time.Now()
	rt := rt_tick.GetRtMod()
	end := time.Since(now_time)
	log.Debug("RunTask GetRtMod time:", end)
	wg.Add(len(rt.Data.Items))
	for i := 0; i < len(rt.Data.Items); i++ {
		rtObj := rt_mod.CreateRtTickMod(rt.Data.Fields, rt.Data.Items[i])
		ants.Submit(wrapper(&rtObj, &wg))
	}
	end = time.Since(now_time)
	log.Debug("RunTask add task time", end, len(rt.Data.Items))

	wg.Wait()
	end = time.Since(now_time)
	log.Debug("RunTask save rtmod", end)
}

var SECTION = flag.Int("section", 0, "全部ID时，分成的总份数")
var INDEX = flag.Int("index", 0, "全部ID时，更新第几份")

func main() {
	defer ants.Release()
	flag.Parse()

	conf.SECTION = *SECTION
	conf.INDEX = *INDEX
	log.Init(conf.GetAppName())

	log.Info("stock_rttick start... ids_type=", conf.GetAppName())

	if url_conf.IsOpenDay() == false {
		log.Info("今天休市，不需要启动程序. 程序休眠6小时。", time.Now())
		time.Sleep(6 * time.Hour)
		os.Exit(0)
	}

	time.Sleep(time.Duration(conf.INDEX) * time.Second)

	//初始化基础数据
	rt_tick.Init()
	mqtt_cli.Init()
	redis_cli.Init()
	rt_tick.UpdateCode(*SECTION, *INDEX)

	isLoop := true
	for isLoop {
		switch conf.MarketState() {
		case conf.OPEN:
			t1 := time.Now()
			RunTask()
			end := time.Since(t1)
			timeConsuming(end)
			if 1*time.Second > end {
				time.Sleep(1*time.Second - end)
			}
			break

		case conf.REST:
			log.Debug("休市时间...", time.Now())
			fmt.Println("休市时间...", time.Now())
			offSecond := 60 - time.Now().Second()
			time.Sleep(time.Duration(offSecond) * time.Second)
			break
		case conf.CLOSE:
			log.Debug("收盘了...", time.Now())
			fmt.Println("收盘了...", time.Now())
			offSecond := 60 - time.Now().Second()
			time.Sleep(time.Duration(offSecond) * time.Second)
			break
		}

	}

	//延迟退出，保证日志写完
	time.Sleep(5 * time.Second)
}

var sampling = 0
var sampling_sum = 1

func timeConsuming(endTime time.Duration) {
	if sampling == 0 {
		log.Debug("run time consuming:", endTime)
	}
	sampling++
	if sampling >= sampling_sum {
		sampling = 0
	}
}
