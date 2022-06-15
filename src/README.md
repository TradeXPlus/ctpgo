# goctp
上海期货交易所 ctp 接口 Golang版 (for linux64)

## 修改配置
    修改 bin/cfg.json，写上对应的环境账号,和使用的策略即可。
    可以配置多个策略，但是当前只能运行一个策略："default": "策略名称",
    "STRATEGYS":{                   策略配置区域
    "default": "Strategy3k",        默认启动的策略
    "Strategy3k": {                 Strategy3k 策略相对应的配置参数
      "symbol": ["FG209", "MA209"],
      "max_klen": 1500,
      "period": 300
      }
    }
    
## 订阅行情
    "StrategyEMA": {
      "symbol": ["FG209", "MA2209"],    此处填写StrategyEMA策略要订阅的行情
      "max_klen": 1500,
      "period": 60
    }

## 开仓示例
    var Input InputOrderStruct

    Input.InstrumentID = "FG209"
    Input.Direction    = OrderBuy
    Input.Price        = 1800
    Input.Volume       = 1

    TraderSpi.OrderOpen(Input)

## 平仓示例
    var Input InputOrderStruct

    Input.InstrumentID = "FG209"
    Input.Direction    = OrderBuy
    Input.Price        = 3600
    Input.Volume       = 1

    TraderSpi.OrderClose(Input)

## 撤单示例
    TraderSpi.OrderCancel("FG209", "报单编号")