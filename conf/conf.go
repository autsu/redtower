package conf

import (
	"errors"
	"flag"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var GlobalConf = &Conf{}
var src = flag.String("c", "", "conf path")

func init() {
	flag.Parse()
	if *src == "" {
		cf := Default()
		GlobalConf = cf
	} else {
		cf, err := New(*src)
		if err != nil {
			panic(err)
		}
		GlobalConf = cf
	}
	if GlobalConf.OpenDebug {
		log.Println(GlobalConf)
	}
}

type Conf struct {
	MaxConnNum               int  `yaml:"MaxConnNum"`
	MaxPackageSize           int  `yaml:"MaxPackageSize"`
	GoroutineMaxNum          int  `yaml:"GoroutineMaxNum"`
	WorkPoolQueueSize        int  `yaml:"WorkPoolQueueSize"`
	HeartBeatSendingInterval int  `yaml:"HeartBeatSendingInterval"`
	HeartBeatDeadline        int  `yaml:"HeartBeatDeadline"`
	OpenDebug                bool `yaml:"OpenDebug"`
}

func New(src string) (cf *Conf, err error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cf = &Conf{}
	if err = yaml.NewDecoder(file).Decode(cf); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("conf file [conf.yaml] not exist, DO NOT MOVE or RENAME")
		}
		return
	}
	return
}

func Default() *Conf {
	return &Conf{
		MaxConnNum:               1000,
		MaxPackageSize:           66666666,
		GoroutineMaxNum:          1000,
		WorkPoolQueueSize:        100,
		HeartBeatSendingInterval: 10,
		HeartBeatDeadline:        30,
		OpenDebug:                true,
	}
}

//var GlobalServerConfig = &ServerConf{}
//var GlobalClientConfig = &ClientConf{}
//
//type ServerConf struct {
//	MaxConnNum               uint64        // 最大连接数
//	MaxPackageSize           uint64        // 包的最大大小
//	GoroutineMaxNum          uint64        // 最多开启 goroutine 的数量
//	WorkPoolQueueSize        uint64        // workpool 的队列大小
//	HeartBeatSendingInterval time.Duration // 心跳包每隔多少时间发送一次
//	HeartBeatDeadline        time.Duration // 多长时间没收到心跳包则判定为失活，将断开连接
//	OpenDebug                bool          // 是否打印日志
//}
//
//const (
//	defaultMaxConnNum               = uint64(1000)
//	defaultMaxPackageSize           = uint64(65535)
//	defaultGoroutineMaxNum          = uint64(1000)
//	defaultWorkPoolQueueSize        = uint64(100)
//	defaultHeartBeatSendingInterval = time.Second * 10
//	defaultHeartBeatDeadline        = time.Second * 30
//)
//
//// SetDefault 检查调用方是否设置值，如果没有则设置为默认值
//func (s *ServerConf) SetDefault() {
//	if s.MaxConnNum == 0 {
//		s.MaxConnNum = defaultMaxConnNum
//	}
//	if s.MaxPackageSize == 0 {
//		s.MaxPackageSize = defaultMaxPackageSize
//	}
//	if s.GoroutineMaxNum == 0 {
//		s.GoroutineMaxNum = defaultGoroutineMaxNum
//	}
//	if s.WorkPoolQueueSize == 0 {
//		s.WorkPoolQueueSize = defaultWorkPoolQueueSize
//	}
//	if s.HeartBeatSendingInterval == 0 {
//		s.HeartBeatSendingInterval = defaultHeartBeatSendingInterval
//	}
//	if s.HeartBeatDeadline == 0 {
//		s.HeartBeatDeadline = defaultHeartBeatDeadline
//	}
//}
//
//type ClientConf struct {
//	HeartBeatSendingInterval time.Duration // 心跳包每隔多少时间发送一次
//}
//
//func (c *ClientConf) SetDefault() {
//	if c.HeartBeatSendingInterval == 0 {
//		c.HeartBeatSendingInterval = defaultHeartBeatSendingInterval
//	}
//}
