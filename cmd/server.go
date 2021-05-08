package cmd

import (
	_ "exporter-center/docs"
	"exporter-center/group"
	"exporter-center/logs"
	"exporter-center/tool"
	tr "exporter-center/tool/toolRabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gorm.io/gorm"
)

var (
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "Start Api server",
		Example: "no example",
		PreRun: func(cmd *cobra.Command, args []string) {
			usage()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	dbConnection *gorm.DB
	logger       *logrus.Logger
)

func init() {
	// db connection
	dbConnection = tool.GetDbConnection()
	logger = logs.Logger("logs/exporter-center.log")
}

func usage() {
	usageStr := `starting exporter-center server`
	logger.Infoln(usageStr)
}

func run() error {
	router := gin.Default()
	// 使用日志中间件
	router.Use(logs.LoggerHandler())
	logger.Infoln("开始启动rabbitmq consumer")
	// 启动rabbitmq consumer
	go tr.GetMsgFromRabbitmq(logger, dbConnection)
	// es组
	esRouter := router.Group("/es")
	group.EsGroup(esRouter)
	// 电话组
	phoneRouter := router.Group("/phone")
	group.PhoneGroup(phoneRouter, dbConnection)
	// auth
	authGroup := router.Group("/auth")
	group.AuthGroup(authGroup)
	router.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	return router.Run(":8080")
}
