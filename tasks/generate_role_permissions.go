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

package tasks

import (
	"fmt"
	"forgolang_forum/cmn"
	"forgolang_forum/database/model"
	"forgolang_forum/utils"
	"strings"
)

// GenerateRolePermissions generate role permissions for api controller and methods
func GenerateRolePermissions(app *cmn.App, args interface{}) error {
	apiRoutes := GetArg("Router", args).(map[string]map[string][]string)

	role := model.NewRole()
	var roles []model.Role
	app.Database.QueryWithModel(fmt.Sprintf("SELECT * FROM %s", role.TableName()),
		&roles)
	rolesMap := make(map[string]int64)
	for _, r := range roles {
		rolesMap[r.Code] = r.ID
	}

	route := model.NewRoute()
	var routes []model.Route
	app.Database.QueryWithModel(fmt.Sprintf("SELECT * FROM %s", route.TableName()),
		&routes)

	var names []string
	for _, r := range routes {
		names = append(names, r.Name)
	}

	for k, r := range apiRoutes {
		if exists, _ := utils.InArray(k, names); !exists {
			_r := model.NewRoute()
			_r.Name = k
			err := app.Database.Insert(model.NewRoute(), _r, "id")
			if err != nil {
				panic(err)
			}
			app.Cache.Set(strings.Join([]string{cmn.RedisKeys["routes"].(string), _r.Name}, ":"),
				_r.ToJSON(), 0)
			app.Logger.LogInfo(fmt.Sprintf("Generate %s route", _r.Name))
			for k2, r2 := range r {
				for _, r3 := range r2 {
					rp := model.NewRolePermission(rolesMap[k2])
					rp.Controller = k
					rp.Method = r3
					err := app.Database.Insert(model.NewRolePermission(0), rp, "id")
					if err != nil {
						panic(err)
					}
					app.Logger.LogInfo(fmt.Sprintf("Generate %s: %s/%s permission",
						k2,
						rp.Controller,
						rp.Method))
				}
			}
		}
	}
	return nil
}
