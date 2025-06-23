/*
 * Copyright 2023 The RuleGo Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/rulego/rulego/utils/str"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/rulego/rulego-server/config"
	"github.com/rulego/rulego-server/config/logger"
	"github.com/rulego/rulego-server/internal/router"
	"github.com/rulego/rulego-server/internal/service"
	"gopkg.in/ini.v1"

	endpointApi "github.com/rulego/rulego/api/types/endpoint"
	"github.com/rulego/rulego/node_pool"
)

const (
	version = "1.0.0"
)

var (
	//是否是查询版本
	ver bool
	//配置文件
	configFile string
)

func init() {
	flag.StringVar(&configFile, "c", "", "配置文件")
	flag.BoolVar(&ver, "v", false, "打印版本")
}

func main() {
	flag.Parse()

	if ver {
		fmt.Printf("RuleGo-Ci Server v%s", version)
		os.Exit(0)
	}

	// 首先设置默认配置
	c := config.DefaultConfig

	// 如果指定了配置文件，则从配置文件加载配置（环境变量未设置的项）
	if configFile != "" {
		if cfg, err := ini.Load(configFile); err != nil {
			log.Fatal("error:", err)
		} else {
			// 从配置文件加载配置，但保留环境变量已设置的值
			if err := cfg.MapTo(&c); err != nil {
				log.Fatal("error:", err)
			}
			// 加载全局配置
			if section, err := cfg.GetSection("global"); err == nil {
				c.Global = section.KeysHash()
			}
			// 加载用户配置
			if section, err := cfg.GetSection("users"); err == nil {
				c.Users = section.KeysHash()
			}
		}
	}
	// 优先从环境变量加载配置
	c.LoadFromEnv()

	config.Set(c)
	logger.Set(initLogger(c))

	//pprof
	if c.Pprof.Enable {
		addr := c.Pprof.Addr
		if addr == "" {
			addr = "0.0.0.0:6060"
		}
		log.Printf("pprof enabled, addr=%s \n", addr)
		go http.ListenAndServe(addr, nil)
	}

	log.Printf("Get Converter Info: %s \n", str.GetConverterInfo())

	//初始化用户名、密码、apiKey之间的映射
	c.InitUserMap()
	log.Printf("use config file=%s \n", configFile)

	if err := loadNodePool(c); err != nil {
		log.Fatal("loadNodePool error:", err)
	} else {
		log.Printf("loadNodePool file=%s \n", c.NodePoolFile)
	}

	//初始化rulego配置
	router.InitRulegoConfig()
	//创建http服务
	ep, err := router.NewRestServe(c)
	if err != nil {
		log.Fatal("error:", err)
	}
	//启动http服务
	if err := ep.Start(); err != nil {
		log.Fatal("error:", err)
	}
	//创建websocket服务
	if restEp, ok := ep.(endpointApi.HttpEndpoint); ok {
		wsEp, err := router.NewWebsocketServe(c, restEp)
		if err != nil {
			log.Fatal("websocket endpoint creation error:", err)
		}
		if err := wsEp.Start(); err != nil {
			log.Fatal("websocket start error:", err)
		}
	}
	//初始化服务
	if err := service.Setup(c); err != nil {
		log.Fatal("setup service error:", err)
	}

	sigs := make(chan os.Signal, 1)
	// 监听系统信号，包括中断信号和终止信号
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
		if ep != nil {
			ep.Destroy()
		}
		log.Println("stopped server")
		os.Exit(0)
	}
}

// 初始化日志记录器
func initLogger(c config.Config) *log.Logger {
	if c.LogFile == "" {
		return log.New(os.Stdout, "", log.LstdFlags)
	} else {
		f, err := os.OpenFile(c.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		return log.New(f, "", log.LstdFlags)
	}
}

func loadNodePool(c config.Config) error {
	file := c.NodePoolFile
	if file != "" {
		if buf, err := os.ReadFile(file); err != nil {
			return err
		} else {
			_, err = node_pool.DefaultNodePool.Load(buf)
			return err
		}
	}
	return nil
}
