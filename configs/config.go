package configs

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/fuloge/basework/api"
)

type config struct {
	Desc      string
	Pgsql     *PgSqlConfig
	Redis     *RedisConfig
	Rabbitmq  *RabbitmqConfig
	Log       *LogConfig
	Authkey   *AuthkeyConfig
	Fillter   *FillterConfig
	ES        *ESConfig
	RunMode   int
	Concurnum int
}

type ESConfig struct {
	Url     string
	User    string
	Passwd  string
	LogFile string
}

type PgSqlConfig struct {
	Hosts    []string
	Ports    []int
	User     string
	Password string
	Dbname   string
}

type RedisConfig struct {
	Hosts          []string
	Password       string
	DB             int
	MaxActive      int //最大连接数
	MaxIdle        int //最大空闲连接数
	IdleTimeoutSec int //最大空闲连接时间
}

type RabbitmqConfig struct {
	Url       string
	Qgame     string
	Qgamefeed string
	//User     string   // mq user
	//Password string   // mq password
	//Ip       []string // mq ip
	//Port     []int    // mq port
	//Vhost    []string // vhost
	//QuName   []string // 队列名称
	//RtKey    []string // key值
	//ExName   []string // 交换机名称
	//ExType   []string // 交换机类型
}

type LogConfig struct {
	Logfile string
	Sqlog   string
}

type AuthkeyConfig struct {
	Key        string
	Subject    string
	PrivateKey string
	Publickey  string
}

type FillterConfig struct {
	Array []string
}

var (
	confPath  string
	env       string
	logfile   string
	sqlfile   string
	EnvConfig *config
	WhiteList map[string]string
)

func init() {
	flag.StringVar(&env, "env", "dev", "set running env")
	flag.StringVar(&logfile, "logfile", "", "set log file")
	flag.StringVar(&sqlfile, "sqllog", "", "set sql log file")

	confPath = "./configs/datasources-" + env + ".toml"

	_, err := toml.DecodeFile(confPath, &EnvConfig)
	if err != nil {
		panic(api.SysConfigErr)
	} else {
		WhiteList = make(map[string]string)

		if logfile != "" {
			EnvConfig.Log.Logfile = logfile
		}

		if sqlfile != "" {
			EnvConfig.Log.Sqlog = sqlfile
		}

		fmt.Println(EnvConfig.Desc)

		//
		for _, path := range EnvConfig.Fillter.Array {
			WhiteList[path] = path
		}
	}
}
