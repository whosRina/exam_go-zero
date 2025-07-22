package main

import (
	"flag"
	"fmt"

	"exam-system/questionBank/internal/config"
	"exam-system/questionBank/internal/handler"
	"exam-system/questionBank/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/question.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// 创建REST服务并启用CORS
	server := rest.MustNewServer(c.RestConf, rest.WithCors("*")) // 允许所有域访问
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
