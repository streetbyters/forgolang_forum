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

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"forgolang_forum/api"
	"forgolang_forum/cmn"
	"forgolang_forum/database"
	"forgolang_forum/model"
	"forgolang_forum/utils"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	configFile string
	devMode    string
	migrate    bool
	reset      bool
	mode       model.MODE
	appPath    string
	dbPath     string
	genSecret  bool
)

func main() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	flag.StringVar(&devMode, "mode", "dev", "Development Mode")
	flag.BoolVar(&migrate, "migrate", false, "Run migrations")
	flag.BoolVar(&reset, "reset", false, "Reset database")
	flag.StringVar(&configFile, "config", "", "Config file")
	flag.StringVar(&dbPath, "dbPath", "", "Database path")
	flag.StringVar(&appPath, "appPath", "", "Application path")
	flag.BoolVar(&genSecret, "genSecretEnv", false, "Generate secret env file")
	flag.Parse()

	if appPath == "" {
		appPath, _ = os.Getwd()
	}
	dirs := strings.SplitAfter(appPath, "forgolang_forum")

	mode = model.MODE(devMode)
	appPath = path.Join(dirs[0])
	dbPath = appPath

	if genSecret {
		genSecretEnv()
		return
	}

	if configFile == "" {
		configFile = string(mode) + ".env"
	}

	logger := utils.NewLogger(string(mode))

	viper.SetConfigName(configFile)
	viper.AddConfigPath(appPath)
	err := viper.ReadInConfig()
	cmn.FailOnError(logger, err)

	file, err := ioutil.ReadFile(filepath.Join(appPath, "secret.env"))
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewReader(file)
	err = viper.MergeConfig(buffer)
	if err != nil {
		panic(err)
	}

	config := &model.Config{
		EnvFile:      configFile,
		Path:         appPath,
		Port:         viper.GetInt("PORT"),
		SecretKey:    viper.GetString("SECRET_KEY"),
		DB:           model.DB(viper.GetString("DB")),
		DBPath:       dbPath,
		DBName:       viper.GetString("DB_NAME"),
		DBHost:       viper.GetString("DB_HOST"),
		DBPort:       viper.GetInt("DB_PORT"),
		DBUser:       viper.GetString("DB_USER"),
		DBPass:       viper.GetString("DB_PASS"),
		DBSsl:        viper.GetString("DB_SSL"),
		RabbitMqHost: viper.GetString("RABBITMQ_HOST"),
		RabbitMqPort: viper.GetInt("RABBITMQ_PORT"),
		RabbitMqUser: viper.GetString("RABBITMQ_USER"),
		RabbitMqPass: viper.GetString("RABBITMQ_PASS"),
		RedisHost:    viper.GetString("REDIS_HOST"),
		RedisPort:    viper.GetInt("REDIS_PORT"),
		RedisPass:    viper.GetString("REDIS_PASS"),
		RedisDB:      viper.GetInt("REDIS_DB"),
	}

	if config.DB == "" {
		panic(errors.New("enter DB conf"))
	}
	if config.Port == 0 {
		panic(errors.New("enter PORT conf"))
	}

	db, err := database.NewDB(config)
	cmn.FailOnError(logger, err)
	db.Logger = logger

	newApp := cmn.NewApp(config, logger)
	newApp.Channel = ch
	newApp.Database = db
	newApp.Mode = mode

	if migrate {
		db.Reset = reset
		if err := database.InstallDB(db); err != nil {
			panic(err)
		}
		return
	}

	newAPI := api.NewAPI(newApp)
	go func() {
		err := newAPI.Router.Server.ListenAndServe(newAPI.Router.Addr)
		cmn.FailOnError(logger, err)
	}()

	<-newApp.Channel
}

func genSecretEnv() {
	body := []byte(fmt.Sprintf(`SECRET_KEY=%s

CDN_URL=

AWS_SES_ACCESS_KEY_ID=
AWS_SES_SECRET_ACCESS_KEY=
AWS_SES_REGION=
AWS_SES_SOURCE=
AWS_S3_ACCESS_KEY_ID=
AWS_S3_SECRET_ACCESS_KEY=
AWS_S3_BUCKET`, "asdasd"))

	var file *os.File
	_, err := os.Stat(filepath.Join(appPath, "secret.env"))
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filepath.Join(appPath, "secret.env"))
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			if _, err := file.Write(body); err != nil {
				log.Fatal(err)
			}

			log.Println("Generated secret environment file")
			return
		}
	} else {
		log.Println("Secret environment file is exists")
	}
}
