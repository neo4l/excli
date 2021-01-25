# excli
Client for automatic trading of centralized exchange

# Usage of ./excli:
#   -accesskey string
#         交易所api Accesskey
#   -currencyPair string
#         交易对名称
#   -exchange string
#         交易所名称 (default "binance.com")
#   -interval int
#         循环下单间隔多少秒 (default 5)
#   -maxPrice float
#         最高价格,0无限制
#   -minPrice float
#         最低价格,0无限制
#   -orderAmount float
#         订单数量，随机波动：0-5% (default 50)
#   -orderChangeRate float
#         订单价格变化率 (default 0.003)
#   -orderCnt int
#         每次订单数量 (default 3)
#   -orderType int
#         订单类型：1限价买入 2限价卖出 (default 1)
#   -sercetkey string
#         交易所api Sercetkey



# 安全提升
# 1.强烈建议交易所API只开通交易功能，非必要情况不要开通API自动提现权限（有固定IP的，可以增加IP白名单限制）
# 2.保护好自己的API Key，避免泄露（如使用脚本启动，保护好自己的run.sh脚本）

# 终端启动

# 修改命令行中的 交易所API Key，交易对等参数，然后终端执行如下命令即可
# 示例：在币安上用USDT购买Kava，每次3笔挂单，订单数量 800个，价格变化率 0.001，最大价格 2 USDT，每5秒执行一次
# ./excli --exchange binance.com --accesskey aaa --sercetkey bbb --currencyPair KAVA_USDT --orderType 1 --orderCnt 3 --orderAmount 800 --orderChangeRate 0.001 --minPrice 0 --maxPrice 2.0 --interval 5


# 脚本启动

# 修改脚本文件run.sh中的执行参数，然后每次启动都执行脚本即可
# ./run.sh
