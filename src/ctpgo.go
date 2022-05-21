package main

import (
	"ctpgo"
	"fmt"
	"log"
	"os"
)

var (
	Ctp CtpClient // ctp 句柄及配置项

	MdSpi              FtdcMdSpi     // 行情模块函数 句柄
	TraderSpi          FtdcTraderSpi // 交易模块函数 句柄
	MapInstrumentInfos Map           // 交易所合约详情列表 InstrumentInfoStruct

	MapOrderList Map // 报单列表（已成交、未成交、撤单等状态）的列表数据 OrderListStruct

	MdFront     string // ctp 服务器，及交易账号
	TraderFront string
	BrokerID    string
	InvestorID  string
	Password    string

	AppID    string // AppID 客户端认证
	AuthCode string

	StreamFile      = GetCurrentExePath() + "/StreamFile/" // StreamFile ctp 流文件，绝对路径
	OrderBuy   byte = '0'                                  // 买卖方向：买
	OrderSell  byte = '1'                                  //  买卖方向：卖
)

// FtdcMdSpi Ctp 行情 spi 回调函数
type FtdcMdSpi struct {
	CtpClient
}

// FtdcTraderSpi Ctp 交易 spi 回调函数
type FtdcTraderSpi struct {
	CtpClient
}

// CtpClient Ctp 客户端 行情、交易模块 全局变量
type CtpClient struct {
	MdApi     ctpgo.CThostFtdcMdApi     // 行情模块 api
	TraderApi ctpgo.CThostFtdcTraderApi // 交易模块 api

	BrokerID   string // 期货公司代码，用户ID，密码
	InvestorID string
	Password   string

	AppID    string // 客户端认证
	AuthCode string

	TradingDay      string // 当前交易日期
	TradeMonth      string // 当前交易月份
	MdRequestId     int    // 行情请求编号
	TraderRequestId int    // 交易请求编号
	IsTraderInit    bool   // 交易系统是否已经初始化了

	// 交易程序是否初始化完成（自动完成如下动作：交易账号登陆、结算单确认、查询合约、查询资金账户、查询用户报单、查询用户持仓 后算完成）
	IsTraderInitFinish bool
	IsTraderLogin      bool // 交易程序是否已登录过
	IsMdLogin          bool // 行情程序是否已登录过
}

// SetTradeAccount 设置交易账号
func SetTradeAccount(RunMode string) {

	cfg := LoadJson("cfg.json")
	ukey := "USERS." + RunMode
	if !cfg.Get(ukey).Exists() {
		_, err := Println("该模式未设置交易账号信息")
		if err != nil {

		}
		os.Exit(1)
	}

	MdFront = cfg.Get(ukey + ".m_host").String()
	TraderFront = cfg.Get(ukey + ".t_host").String()
	BrokerID = cfg.Get(ukey + ".bid").String()
	InvestorID = cfg.Get(ukey + ".uid").String()
	Password = cfg.Get(ukey + ".pwd").String()
	AppID = cfg.Get(ukey + ".app_id").String()
	AuthCode = cfg.Get(ukey + ".auth_code").String()
}

func init() {
	// 全局 行情、交易 函数句柄
	MdSpi = FtdcMdSpi{}
	TraderSpi = FtdcTraderSpi{}
}

func main() {
	RunMode := "test" // 运行模式【运行程序时带上参数可设置】,需要在cfg.json中配置参数
	if len(os.Args) >= 2 {
		RunMode = os.Args[1]
	}
	SetTradeAccount(RunMode) // 设置交易账号

	log.Println("启动交易程序")

	// 检查流文件目录是否存在
	fileExists := IsDirExist(StreamFile)
	if !fileExists {
		err := os.Mkdir(StreamFile, os.ModePerm)
		if err != nil {
			fmt.Println("创建目录失败，请检查是否有操作权限")
		}
	}

	Ctp = CtpClient{
		MdApi:              ctpgo.CThostFtdcMdApiCreateFtdcMdApi(StreamFile),
		TraderApi:          ctpgo.CThostFtdcTraderApiCreateFtdcTraderApi(StreamFile),
		BrokerID:           BrokerID,
		InvestorID:         InvestorID,
		Password:           Password,
		AppID:              AppID,
		AuthCode:           AuthCode,
		MdRequestId:        1,
		TraderRequestId:    1,
		IsTraderInit:       false,
		IsTraderInitFinish: false,
		IsMdLogin:          false,
		IsTraderLogin:      false,
	}

	Ctp.MdApi.RegisterSpi(ctpgo.NewDirectorCThostFtdcMdSpi(&FtdcMdSpi{Ctp}))
	Ctp.MdApi.RegisterFront(MdFront)
	Ctp.MdApi.Init()

	Ctp.TraderApi.RegisterSpi(ctpgo.NewDirectorCThostFtdcTraderSpi(&FtdcTraderSpi{Ctp}))
	Ctp.TraderApi.RegisterFront(TraderFront)

	Ctp.TraderApi.SubscribePublicTopic(ctpgo.THOST_TERT_QUICK)
	Ctp.TraderApi.SubscribePrivateTopic(ctpgo.THOST_TERT_QUICK)
	Ctp.TraderApi.Init()
	Ctp.TraderApi.Join()

	// .Join() 如果后面有其它需要处理的功能可以不写，但必须保证程序不能退出，Join 就是保证程序不退出的
}
