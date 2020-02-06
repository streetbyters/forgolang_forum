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

package model

import (
	"forgolang_forum/database"
	"time"
)

// User2fa user 2fa type structure
type User2fa struct {
	database.DBInterface `json:"-"`
	ID                   int64              `db:"id" json:"id"`
	UserID               int64              `db:"user_id" json:"user_id" foreign:"fk_user_2fa_user_id" validate:"required"`
	Type                 database.TwoFactor `db:"type" json:"type"`
	InsertedAt           time.Time          `db:"inserted_at" json:"inserted_at"`
}

// NewUser2fa generate user 2fa structure
func NewUser2fa(userID int64) *User2fa {
	return &User2fa{UserID: userID, Type: database.Email}
}

// TableName user 2fa database
func (m User2fa) TableName() string {
	return "user_2fa"
}

// ToJSON user 2fa structure to json string
func (m User2fa) ToJSON() string {
	return database.ToJSON(m)
}
