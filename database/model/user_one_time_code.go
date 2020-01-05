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

package model

import (
	"forgolang_forum/database"
	"time"
)

// UserOneTimeCode user 2fa code structure
type UserOneTimeCode struct {
	database.DBInterface `json:"-"`
	ID                   int64        `db:"id" json:"id"`
	UserID               int64        `db:"user_id" json:"user_id" foreign:"fk_user_one_time_codes_user_id" validate:"required"`
	Code                 string       `db:"code" json:"code" unique:"user_one_time_codes_code_unique" validate:"required"`
	Type                 database.OTC `db:"type" json:"type"`
	InsertedAt           time.Time    `db:"inserted_at" json:"inserted_at"`
}

// NewUserOneTimeCode generate user 2fa code structure
func NewUserOneTimeCode(userID int64) *UserOneTimeCode {
	return &UserOneTimeCode{UserID: userID, Type: database.Confirmation}
}

// TableName user one time code database
func (m UserOneTimeCode) TableName() string {
	return "user_one_time_codes"
}

// ToJSON user one time code structure to json string
func (m UserOneTimeCode) ToJSON() string {
	return database.ToJSON(m)
}
