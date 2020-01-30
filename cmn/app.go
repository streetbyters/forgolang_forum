// Copyright 2019 Forgolang Community
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmn

import (
	"fmt"
	"forgolang_forum/database"
	"forgolang_forum/model"
	"forgolang_forum/thirdparty/aws"
	"forgolang_forum/thirdparty/github"
	"forgolang_forum/utils"
	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	"github.com/olivere/elastic/v7"
	"github.com/streadway/amqp"
	"net/url"
	"os"
	"strconv"
)

// RedisKeys application cache keys
var RedisKeys = make(map[string]interface{})

// App structure
type App struct {
	Database      *database.Database
	Channel       chan os.Signal
	Config        *model.Config
	Logger        *utils.Logger
	Mode          model.MODE
	Storage       *aws.S3
	Email         *aws.SES
	Cache         *redis.Client
	Amqp          *amqp.Connection
	Queue         *Queue
	Github        *github.Github
	HttpClient    *resty.Client
	ElasticClient *elastic.Client
}

// NewApp building new app
func NewApp(config *model.Config, logger *utils.Logger) *App {
	app := &App{
		Config: config,
		Logger: logger,
	}

	awsConfig, err := aws.Config(app.Config.Path, app.Config.EnvFile)
	FailOnError(logger, err)
	app.Storage, err = aws.NewS3(awsConfig)
	FailOnError(logger, err)
	app.Email, err = aws.NewSES(awsConfig)
	FailOnError(logger, err)

	app.Cache = redis.NewClient(&redis.Options{
		Network:     "tcp",
		Addr:        fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password:    config.RedisPass,
		DB:          config.RedisDB,
		MaxRetries:  3,
		PoolSize:    10,
		PoolTimeout: 15000,
		IdleTimeout: 15000,
	})

	ping := app.Cache.Ping()
	FailOnError(logger, ping.Err())
	logger.LogInfo("Started redis connection")

	uri := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(app.Config.RabbitMqUser, app.Config.RabbitMqPass),
		Host:   app.Config.RabbitMqHost + ":" + strconv.Itoa(app.Config.RabbitMqPort),
	}
	app.Amqp, err = amqp.Dial(uri.String())
	FailOnError(logger, err)
	logger.LogInfo("Started rabbitMQ connection")

	app.HttpClient = resty.New()
	app.ElasticClient, _ = elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%d", app.Config.ElasticHost, app.Config.ElasticPort)),
	)
	RedisKeys["permissions"] = "permissions"
	RedisKeys["routes"] = "routes"
	RedisKeys["user"] = map[string]string{
		"one":         "user",
		"permissions": "user:permissions",
		"permission":  "user:permission",
	}
	RedisKeys["category"] = map[string]string{
		"all": "categories",
		"one": "category",
		"slug": "category:slug",
	}
	RedisKeys["tag"] = map[string]string{
		"all": "tags",
		"one": "tag",
		"count": "tag:count",
	}
	RedisKeys["post"] = map[string]string{
		"all": "posts",
		"one": "post",
		"count": "post:count",
	}

	app.Queue = NewQueue(app).StartAll()
	app.Github = github.NewGithub(config)

	return app
}

// FailOnError panic error with logger
func FailOnError(logger *utils.Logger, err error) {
	if err != nil {
		logger.Panic().Err(err)
		panic(err)
	}
}

// GetRedisKey get a key with given keys
func GetRedisKey(keys ...string) string {
	if len(keys) > 1 {
		k := RedisKeys[keys[0]]
		return k.(map[string]string)[keys[1]]
	}
	return RedisKeys[keys[0]].(string)
}
