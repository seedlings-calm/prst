package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/seedlings-calm/prst/app"
	"github.com/seedlings-calm/prst/app/router"
	"github.com/seedlings-calm/prst/common"
	cfg "github.com/seedlings-calm/prst/config"
	"github.com/seedlings-calm/prst/db"
	"github.com/seedlings-calm/prst/middleware"
)

var AppRouters = make([]func(r *gin.Engine, mw *common.GinJWTMiddleware), 0)

func init() {
	//  注册路由 fixme 其他应用的路由，在本目录新建文件放在init方法
	AppRouters = append(AppRouters, router.InitRouter)
}

// @title prst API
// @version 0.0.1
// @description gin框架API
func main() {
	// 定义已知值
	base := 1 + 0.2
	power := 90.0
	result := 100.0

	// 计算 x
	x := result / math.Pow(base, power)

	// 输出结果
	fmt.Printf("x = %.14e\n", x)
	fmt.Printf("x = %.14f\n", x)
	return
	//初始化配置信息
	_ = cfg.Setup()
	//初始化jwt
	jwtMW, err := app.JWTInit()
	if err != nil {
		panic("初始化jwt失败")
	}
	//初始化db
	db.GormMysql()

	//初始化zaplogger
	common.LoggerInit()

	gin.SetMode(cfg.ModelSwitchGinModel())

	r := middleware.InitEngine()

	//加载所有路由
	for _, f := range AppRouters {
		f(r, jwtMW)
	}
	confg := cfg.GetGlobalConf()
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", confg.App.Host, confg.App.Port),
		Handler: r,
	}

	go func() {
		// 服务连接
		if confg.App.Enable {
			if err := srv.ListenAndServeTLS(confg.Ssl.Pem, confg.Ssl.Key); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("listen: ", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("listen: ", err)
			}
		}
	}()
	fmt.Println(common.Green("Server run at:"))
	fmt.Printf("-  Local:   %s://localhost:%d/ \r\n", "http", confg.App.Port)
	fmt.Printf("-  Local:   %s://localhost:%d/ \r\n", "https", confg.App.Port)

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
