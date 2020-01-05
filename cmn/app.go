// Copyright 2019 Abdulkadir DILSIZ
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
	"forgolang_forum/utils"
	"github.com/go-redis/redis"
	"os"
)

// RedisKeys application cache keys
var RedisKeys = make(map[string]interface{})

// App structure
type App struct {
	Database *database.Database
	Channel  chan os.Signal
	Config   *model.Config
	Logger   *utils.Logger
	Mode     model.MODE
	Storage  *aws.S3
	Email    *aws.SES
	Cache    *redis.Client
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
	}

	return app
}

// FailOnError panic error with logger
func FailOnError(logger *utils.Logger, err error) {
	if err != nil {
		logger.Panic().Err(err)
		panic(err)
	}
}

// GetRedisKeys get a key with given keys
func GetRedisKey(keys ...string) string {
	if len(keys) > 1 {
		k := RedisKeys[keys[0]]
		return k.(map[string]string)[keys[1]]
	}
	return RedisKeys[keys[0]].(string)
}