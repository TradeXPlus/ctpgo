package main

import (
	"fmt"
	"github.com/TradeXPlus/ctpgo/lib"
	"github.com/TradeXPlus/ctpgo/strategy"
	"github.com/TradeXPlus/ctpgo/utils"
	"github.com/tidwall/gjson"
	"log"
	"os"
)

var (
	gCtp       strategy.CtpClient                           // ctp 句柄及配置项
	gMdSpi     strategy.FtdcMdSpi                           // 行情模块函数 句柄
	gTraderSpi strategy.FtdcTraderSpi                       // 交易模块函数 句柄
	gCfg       strategy.Config                              // 加载配置文件项目
	StreamFile = utils.GetCurrentExePath() + "/StreamFile/" // StreamFile ctp 流文件，绝对路径
)

// LoadCfg 设置交易账号等相关的参数
func LoadCfg(RunMode string) {
	cfg := utils.LoadJson("cfg.json")
	if RunMode == "" {
		RunMode = cfg.Get("USERS.default").String()
	}
	ukey := "USERS." + RunMode
	if !cfg.Get(ukey).Exists() {
		_, err := utils.Println("该模式未设置交易账号信息")
		if err != nil {
		}
		os.Exit(1)
	}
	class := cfg.Get("STRATEGYS.default").String() //策略配置默认参数
	sKey := "STRATEGYS." + class

	gCfg = strategy.Config{
		MdFront:     cfg.Get(ukey + ".m_host").String(),
		TraderFront: cfg.Get(ukey + ".t_host").String(),
		BrokerID:    cfg.Get(ukey + ".bid").String(),
		InvestorID:  cfg.Get(ukey + ".uid").String(),
		Password:    cfg.Get(ukey + ".pwd").String(),
		AppID:       cfg.Get(ukey + ".app_id").String(),
		AuthCode:    cfg.Get(ukey + ".auth_code").String(),

		Class:   class, //策略struct名称
		MaxKlen: cfg.Get(sKey + ".max_klen").Int(),
		Period:  cfg.Get(sKey + ".period").Int(), //策略struct名称
	}

	// 加载策略获取行情订阅
	cfg.Get(sKey + ".symbol").ForEach(func(key, value gjson.Result) bool {
		gCfg.Symbol = append(gCfg.Symbol, value.String())
		return true
	})
	// 加载策略参数
	cfg.Get(sKey + ".params").ForEach(func(key, value gjson.Result) bool {
		gCfg.Params = append(gCfg.Params, value.Int())
		return true
	})
}

func init() {
	// 全局 行情、交易 函数句柄
	RunMode := "product" // 运行模式【运行程序时带上参数可设置】,需要在cfg.json中配置参数
	if len(os.Args) > 1 {
		RunMode = os.Args[1]
	}
	LoadCfg(RunMode) // 设置交易相关参数，账号

	// 检查流文件目录是否存在
	fileExists := utils.IsDirExist(StreamFile)
	if !fileExists {
		err := os.Mkdir(StreamFile, os.ModePerm)
		if err != nil {
			fmt.Println("创建目录失败，请检查是否有操作权限")
			os.Exit(2)
		}
	}

	gCtp = strategy.CtpClient{
		MdApi:     lib.CThostFtdcMdApiCreateFtdcMdApi(StreamFile),
		TraderApi: lib.CThostFtdcTraderApiCreateFtdcTraderApi(StreamFile),
		Config:    gCfg,

		MdRequestId:        1,
		TraderRequestId:    1,
		IsTraderInit:       false,
		IsTraderInitFinish: false,
		IsMdLogin:          false,
		IsTraderLogin:      false,
	}
	//fmt.Printf("%+v\n", gCfg)
	if gCtp.Register(&gCfg, &gTraderSpi) { //注册策略
		gMdSpi = strategy.FtdcMdSpi{CtpClient: &gCtp}
		gTraderSpi = strategy.FtdcTraderSpi{CtpClient: &gCtp}
	} else {
		fmt.Printf("注册策略： %v 失败！\n", gCfg.Class)
		os.Exit(3)
	}
}

func main() {
	log.Println("启动交易程序")

	gCtp.MdApi.RegisterSpi(lib.NewDirectorCThostFtdcMdSpi(&gMdSpi))
	gCtp.MdApi.RegisterFront(gCfg.MdFront)
	gCtp.MdApi.Init()

	gCtp.TraderApi.RegisterSpi(lib.NewDirectorCThostFtdcTraderSpi(&gTraderSpi))
	gCtp.TraderApi.RegisterFront(gCfg.TraderFront)

	gCtp.TraderApi.SubscribePublicTopic(lib.THOST_TERT_QUICK)
	gCtp.TraderApi.SubscribePrivateTopic(lib.THOST_TERT_QUICK)
	gCtp.TraderApi.Init()

	gCtp.TraderApi.Join()
	// .Join() 如果后面有其它需要处理的功能可以不写，但必须保证程序不能退出，Join 就是保证程序不退出的
}
