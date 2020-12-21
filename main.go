package main

import (
	"flag"
	"fmt"
	"time"

	"excli/ex"

	"github.com/nntaoli-project/goex"
)

func main() {

	var accesskey string
	var sercetkey string
	var exchange string
	var currencyPair string
	var orderType int
	var orderAmount float64
	var orderCnt int
	var orderChangeRate float64
	var minPrice float64
	var maxPrice float64
	var interval int

	flag.StringVar(&accesskey, "accesskey", "", "交易所api Accesskey")
	flag.StringVar(&sercetkey, "sercetkey", "", "交易所api Sercetkey")
	flag.StringVar(&exchange, "exchange", goex.BINANCE, "交易所名称")
	flag.StringVar(&currencyPair, "currencyPair", "", "交易对名称")
	flag.IntVar(&orderType, "orderType", 1, "订单类型：1限价买入 2限价卖出")
	flag.IntVar(&orderCnt, "orderCnt", 3, "最大订单数量")
	flag.Float64Var(&orderAmount, "orderAmount", 50, "单笔订单代币数量")
	flag.Float64Var(&orderChangeRate, "orderChangeRate", 0.003, "订单价格变化率")
	flag.Float64Var(&minPrice, "minPrice", 0, "最低价格,0无限制")
	flag.Float64Var(&maxPrice, "maxPrice", 0, "最高价格,0无限制")
	flag.IntVar(&interval, "interval", 5, "循环下单间隔多少秒")
	flag.Parse()

	exJob := ex.NewJob()
	exJob.API = ex.NewAPI(accesskey, sercetkey, exchange)
	exJob.CurrencyPair = goex.NewCurrencyPair2(currencyPair)
	exJob.OrderType = orderType
	exJob.OrderAmount = orderAmount
	exJob.OrderCnt = orderCnt
	exJob.OrderChangeRate = orderChangeRate
	exJob.MinPrice = minPrice
	exJob.MaxPrice = maxPrice
	exJob.Interval = time.Duration(interval) * time.Second
	fmt.Printf("exchange: %s\n", exchange)
	fmt.Printf("accesskey: %s\n", accesskey)
	fmt.Printf("sercetkey: %s\n", sercetkey)
	fmt.Println(ex.ToJson(exJob))
	exJob.Run()

}
