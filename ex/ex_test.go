package ex

import (
	"testing"
	"time"

	"github.com/nntaoli-project/goex"
)

func Test_Run(t *testing.T) {
	exJob := &EXJob{}
	exJob.API = NewAPI("", "", goex.BINANCE)
	exJob.CurrencyPair = goex.NewCurrencyPair2("KSM_USDT")
	exJob.OrderType = 1
	exJob.OrderAmount = 1000
	exJob.OrderCnt = 2
	exJob.OrderChangeRate = 0.003
	exJob.MinPrice = 0.006
	exJob.MaxPrice = 0.0064
	exJob.Interval = 10 * time.Second
	//Println(goex.NewCurrencyPair2("MIX_KRW"))
	//Println(exJob.API.GetAccount())
	ticker, _ := exJob.GetTicker()
	Println(ToJson(ticker))
	// exJob.Run()

}
