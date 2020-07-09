package api

type Message struct {
	Item    string `json:"item"`
	Subject string `json:"subject"`
	Type    int    `json:"type"` //#0发起 1反馈
	Status  bool   `json:"status"`
	GameId int64  `json:"id"`
	Id      int64  `json:"id"`
	Data    []byte `json:"data"`
}

type Account struct {
	Balance  int64 `json:"balance"`  //余额
	Unsetled int64 `json:"unsetled"` //未结算余额
}

// 定义错误码
type Errno struct {
	Code    int
	Message string
}

//各系统模块错误码范围
//1000-1999 系统框架
//2000-2999 游戏大厅
//3000-3999 游戏
//4000-4999 基础后台
var (
	OK = &Errno{Code: 0, Message: "OK"}

	// 系统错误, 前缀为 100
	SysInternalServerError = &Errno{Code: 1001, Message: "内部服务器错误"}
	SysErrEncrypt          = &Errno{Code: 1002, Message: "加密用户密码时发生错误"}
	SysConfigErr           = &Errno{1003, "没有找到系统配置文件"}

	// 数据库错误, 前缀为 101
	DBMrgErr  = &Errno{Code: 1011, Message: "数据库管理器生成异常"}
	DBConnErr = &Errno{Code: 1012, Message: "数据库连接异常"}
	DBRunErr  = &Errno{Code: 1013, Message: "数据库运行时异常"}
	DBLogErr  = &Errno{Code: 1014, Message: "数据库日志文件生成异常"}

	// 认证错误, 前缀是 102
	AuthErr      = &Errno{Code: 1021, Message: "验证失败"}
	AuthExp      = &Errno{Code: 1022, Message: "token过期"}
	AuthParseErr = &Errno{Code: 1023, Message: "token解析异常"}

	//redis错误
	RedisConnErr = &Errno{Code: 1031, Message: "redis连接异常"}
	RedisRun     = &Errno{Code: 1032, Message: "redis运行时异常"}

	//mq错误
	MQConnErr = &Errno{Code: 1041, Message: "MQ连接异常"}
	MQRun     = &Errno{Code: 1042, Message: "MQ运行时异常"}

	//token
	TokenInvidErr = &Errno{Code: 1051, Message: "token无效"}
	TokenNilErr   = &Errno{Code: 1052, Message: "token为空"}

	//加解密
	RSADecERR = &Errno{Code: 1061, Message: "解密异常"}
	RSAEncERR = &Errno{Code: 1062, Message: "加密异常"}

	//http
	HTTPErr      = &Errno{Code: 1071, Message: "http异常"}
	HTTPUidErr   = &Errno{Code: 1072, Message: "uid无效"}
	HTTPParamErr = &Errno{Code: 1073, Message: "参数异常"}

	//google
	GoogleAuthGetErr    = &Errno{Code: 1081, Message: "获取动态码异常"}
	GoogleAuthVerifyErr = &Errno{Code: 1082, Message: "验证动态码异常"}

	//mq
	MQSyncErr    = &Errno{Code: 1091, Message: "需要同步"}
	MQProcessErr = &Errno{Code: 1092, Message: "消息处理异常"}
)

const (
	PREF_USER = "user:USER_"
	PREF_WALL = "wallet:WALLET_"
	PREF_TOKEN = "token:USER_"
)
