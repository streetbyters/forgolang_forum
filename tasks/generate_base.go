// Copyright 2019 Street Byters Community
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

package tasks

import (
	"context"
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database/model"
	"time"
)

const elasticBody = `{"settings":{"analysis":{"analyzer":{"default":{"tokenizer":"standard","filter":["ascii"]}},"filter":{"ascii":{"type":"asciifolding","preserve_original":true}}}}}`

// GenerateBase artifacts
func GenerateBase(app *cmn.App, args interface{}) error {
	var language model.Language
	app.Logger.LogInfo("Start generate base artifacts")

	reset := GetArg("Reset", args).(bool)

	if reset {
		app.Cache.FlushDB()
		_, err := app.ElasticClient.DeleteIndex("users", "posts", "comments").Do(context.TODO())
		if err != nil {
			panic(err)
		}
		app.Logger.LogInfo("Reset elasticsearch indexes")

		result := app.Database.Query(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE",
			language.TableName()))
		if result.Error != nil {
			panic(result.Error)
		}
		app.Logger.LogInfo("Reset languages")
	}

	//Generate elasticsearch indexes
	_, err := app.ElasticClient.
		CreateIndex("users").
		Body(elasticBody).
		Do(context.TODO())
	if err != nil {
		panic(err)
	}
	_, err = app.ElasticClient.
		CreateIndex("posts").
		Body(elasticBody).
		Do(context.TODO())
	if err != nil {
		panic(err)
	}
	_, err = app.ElasticClient.
		CreateIndex("comments").
		Body(elasticBody).
		Do(context.TODO())
	if err != nil {
		panic(err)
	}
	app.Logger.LogInfo("Generate elasticsearch indexes")

	languageTR := model.NewLanguage()
	languageTR.Code = "tr-TR"
	languageTR.Name = "Turkce"
	languageTR.DateFormat.SetValid(time.RFC3339Nano)
	err = app.Database.Insert(new(model.Language), languageTR, "id")
	if err != nil {
		panic(err)
	}

	languageEN := model.NewLanguage()
	languageEN.Code = "en-US"
	languageEN.Name = "English (U.S)"
	languageEN.DateFormat.SetValid(time.RFC3339Nano)
	err = app.Database.Insert(new(model.Language), languageEN, "id")
	if err != nil {
		panic(err)
	}

	app.Logger.LogInfo("Finish generate base artifacts")
	return nil
}
