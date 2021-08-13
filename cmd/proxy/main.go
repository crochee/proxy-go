package main

import (
	"flag"
	"log"

	"github.com/crochee/proxy-go"
	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/pkg/logger"
)

func main() {
	configFile := flag.String("c", "./conf/config.yml", "")
	flag.Parse()

	// 1.如果我们在main()函数里直接传递一个非零值退出码给os.Exit(),
	// 那么倘若我们在main()函数里事先已经登记了defer函数，那么这个defer函数将不会被执行
	// 2.直接调用os.Exit的话, 如果正好同一时刻正好有其他goroutine在执行且恰好它触发了panic, 那么这个panic将无法被呈现出来

	// 初始化配置
	if err := config.InitConfig(*configFile); err != nil {
		log.Fatal(err)
	}

	if err := proxygo.Server(); err != nil {
		logger.Fatal(err.Error())
	}

	// 正常退出时，返回码0
	// panic为2
}
