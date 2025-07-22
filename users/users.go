package main

import (
	"context"
	"exam-system/users/internal/config"
	"exam-system/users/internal/handler"
	"exam-system/users/internal/svc"
	"exam-system/users/model"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"log"
)

var configFile = flag.String("f", "etc/users.yaml", "the config file")

func main() {
	stat.SetReporter(nil) // 禁用统计日志
	flag.Parse()

	var c config.Config
	// 加载配置文件
	conf.MustLoad(*configFile, &c)

	// 创建REST服务并启用CORS
	server := rest.MustNewServer(c.RestConf, rest.WithCors("*")) // 允许所有域访问
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 初始化数据库连接
	conn := sqlx.NewMysql(c.DataSource)
	usersModel := model.NewUsersModel(conn)

	if err := usersModel.InitAdminUser(context.Background()); err != nil {
		log.Printf("初始化管理员用户失败，但服务仍将继续运行: %v", err)
	}

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
