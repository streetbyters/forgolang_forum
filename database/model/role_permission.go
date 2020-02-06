// Copyright 2019 StreetByters Community
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

package model

import (
	"forgolang_forum/database"
	"time"
)

// RolePermission authorization permission artifacts
type RolePermission struct {
	database.DBInterface `json:"-"`
	ID                   int64     `db:"id" json:"id"`
	RoleID               int64     `db:"role_id" json:"role_id" foreign:"fk_role_permissions_role_id" validate:"required"`
	Controller           string    `db:"controller" json:"controller" validate:"required,gte=1,lte=200"`
	Method               string    `db:"method" json:"method" validate:"required,gte=1,lte=10"`
	InsertedAt           time.Time `db:"inserted_at" json:"inserted_at"`
}

// NewRolePermission generate role permission structure
func NewRolePermission(roleID int64) *RolePermission {
	return &RolePermission{RoleID: roleID}
}

// TableName role permission databasee
func (m RolePermission) TableName() string {
	return "role_permissions"
}

// ToJSON role permission structure to json string
func (m RolePermission) ToJSON() string {
	return database.ToJSON(m)
}
