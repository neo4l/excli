package ex

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"
)

const (
	BUY = 1 + iota
	SELL
	BUY_MARKET
	SELL_MARKET
)

type EXJob struct {
	API             goex.API
	CurrencyPair    goex.CurrencyPair
	OrderType       int
	OrderAmount     float64
	OrderCnt        int
	OrderChangeRate float64
	MaxPrice        float64
	MinPrice        float64
	PricePrecision  int
	AmountPrecision int
	Interval        time.Duration
}

func NewJob() *EXJob {
	return &EXJob{
		OrderCnt:        1,
		OrderChangeRate: 0.005,
		AmountPrecision: 4,
		Interval:        10 * time.Second,
	}
}

func NewJob2(api goex.API, currency string, orderType int, orderAmount float64) *EXJob {
	return &EXJob{
		API:             api,
		CurrencyPair:    goex.NewCurrencyPair2(currency),
		OrderType:       orderType,
		OrderAmount:     orderAmount,
		OrderCnt:        1,
		OrderChangeRate: 0.005,
		AmountPrecision: 4,
		Interval:        10 * time.Second,
	}
}

func (self *EXJob) Run() {
	if err := self.GetPrecision(); err != nil {
		fmt.Printf("Get precision fail: %s", err)
		return
	}
	for {
		Println("run: .........start...........")
		unOrders, err := self.GetUnfinishOrders()
		if err != nil {
			Println("GetUnfinishOrders error:", err)
			continue
		}
		newOrderCnt := self.OrderCnt - len(unOrders)
		if newOrderCnt > 0 {
			self.NewOrders(newOrderCnt)
		} else {
			self.CancelOrders(unOrders)
		}
		Println("run: .........end...........")
		time.Sleep(self.Interval)
	}
}

func (self *EXJob) GetPrecision() error {
	ticker, err := self.GetTicker()
	if err != nil {
		return err
	}
	Println(ToJson(ticker))
	self.PricePrecision = len(strings.Split(strconv.FormatFloat(ticker.Last, 'f', -1, 64), ".")[1])

	amounts := strings.Split(strconv.FormatFloat(self.OrderAmount, 'f', -1, 64), ".")

	self.AmountPrecision = 0
	if len(amounts) > 1 {
		self.AmountPrecision = len(amounts[1])
	}

	fmt.Printf("Precision: %d, %d, %f, %s\n", self.PricePrecision, self.AmountPrecision, ticker.Vol, strings.Split(strconv.FormatFloat(ticker.Last, 'f', -1, 64), ".")[1])
	return nil
}

func (self *EXJob) NewOrders(orderCnt int) error {
	Println("method: NewOrders.......")
	ticker, err := self.GetTicker()
	if err != nil {
		return err
	}
	Println(ToJson(ticker))

	//self.NewOrderMarket(ticker)

	for i := 0; i < orderCnt; i++ {

		amount := Float64ToString(self.OrderAmount+float64(rand.Intn(50)), self.AmountPrecision)
		price, err := self.GetPrice(ticker, i)
		if err != nil {
			continue
		}

		fmt.Printf("new order: %d, %s, %s\n", self.OrderType, amount, price)
		order, err := self.NewOrder(self.OrderType, amount, price)
		if err != nil {
			Println("new order error:", err)
			return err
		}
		fmt.Printf("new order: %s, %s, %f/%f, %f, %s\n", order.OrderID2, order.Side, order.DealAmount, order.Amount, order.Price, time.Now().Format("2006-01-02 15:04:05"))
	}
	return nil
}

func (self *EXJob) CancelOrders(unOrders []goex.Order) {
	for _, order := range unOrders {
		fmt.Printf("opened order: %s, %s, %f/%f, %f, %s\n", order.OrderID2, order.Side, order.DealAmount, order.Amount, order.Price, time.Unix(int64(order.OrderTime/1000), 0).Format("2006-01-02 15:04:05"))
		if time.Now().Unix()-int64(order.OrderTime/1000) > 10 {
			b, err := self.CancelOrder(order.OrderID2)
			Println("cancle order:", order.OrderID2, b, err)
		}
	}
}

func (self *EXJob) GetPrice(ticker *goex.Ticker, seq int) (string, error) {
	var price float64
	switch self.OrderType {
	case 1, 3:
		price = ticker.Buy * (1.0001 - self.OrderChangeRate*float64(seq))
	case 2, 4:
		price = ticker.Sell * (0.9999 + self.OrderChangeRate*float64(seq))
	default:
		price = 0
	}
	if price <= 0 || (self.MaxPrice > 0 && price > self.MaxPrice) || (self.MinPrice > 0 && price < self.MinPrice) {
		return "", fmt.Errorf("Price not allowed: %f, Min %f, Max %f\n", price, self.MinPrice, self.MaxPrice)
	}
	return Float64ToString(price, self.PricePrecision), nil
}

func (self *EXJob) NewOrder(orderType int, amount, price string) (*goex.Order, error) {
	//fmt.Printf("new order: %d, %s, %s\n", orderType, amount, price)
	switch orderType {
	case 1:
		return self.API.LimitBuy(amount, price, self.CurrencyPair)
	case 2:
		return self.API.LimitSell(amount, price, self.CurrencyPair)
	case 3:
		return self.API.MarketBuy(amount, price, self.CurrencyPair)
	case 4:
		return self.API.MarketSell(amount, price, self.CurrencyPair)
	default:
		return nil, errors.New("unknow order type")
	}
}

func ShowOrder(o goex.Order) {
	amount := fmt.Sprintf("%f(%f)", o.Amount, o.DealAmount)
	price := fmt.Sprintf("%f(%f)", o.Price, o.AvgPrice)
	orderTime := time.Unix(int64(o.OrderTime), 0).Format("2006-01-02 15:04:05")
	Println("order: ", o.Side, o.OrderID2, amount, price, o.Status, orderTime)
}

func NewAPI(key, secretKey, exchange string) goex.API {
	return builder.NewAPIBuilder().APIKey(key).APISecretkey(secretKey).Build(exchange)
}

func NewAPIWithProxy(key, secretKey, exchange, proxy string) goex.API {
	//HttpProxy("socks5://127.0.0.1:1086")
	return builder.NewAPIBuilder().APIKey(key).APISecretkey(secretKey).HttpProxy(proxy).Build(exchange)
}

func (self *EXJob) GetTicker() (*goex.Ticker, error) {
	return self.API.GetTicker(self.CurrencyPair)
}

func (self *EXJob) GetUnfinishOrders() ([]goex.Order, error) {
	return self.API.GetUnfinishOrders(self.CurrencyPair)
}

func (self *EXJob) CancelOrder(orderId string) (bool, error) {
	return self.API.CancelOrder(orderId, self.CurrencyPair)
}

func Println(args ...interface{}) {
	fmt.Println(args...)
}

func ToJson(obj interface{}) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func Float64ToString(v float64, precision int) string {
	return strconv.FormatFloat(v, 'f', precision, 64)
}

func FloatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}
