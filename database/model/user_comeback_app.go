// Copyright 2019 Abdulkadir Dilsiz - Çağatay Yücelen
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
	"github.com/jmoiron/sqlx/types"
	"gopkg.in/guregu/null.v3/zero"
	"time"
)

// UserComebackApp integrated third-party apps user information structure
type UserComebackApp struct {
	database.DBInterface `json:"-"`
	ID                   int64          `db:"id" json:"id"`
	UserID               int64          `db:"user_id" json:"user_id" unique:"user_comeback_apps_user_tparty_unique" validate:"required"`
	TPartyID             int64          `db:"tparty_id" json:"tparty_id" unique:"user_comeback_apps_user_tparty_unique" validate:"required"`
	AccessToken          string         `db:"access_token" json:"access_token" validate:"required"`
	RefreshToken         zero.String    `db:"refresh_token" json:"refresh_token"`
	Expire               zero.Int       `db:"expire" json:"expire"`
	Data                 types.JSONText `db:"data" json:"data"`
	InsertedAt           time.Time      `db:"inserted_at" json:"inserted_at"`
	UpdatedAt            time.Time      `db:"updated_at" json:"updated_at"`
}

// NewUserComebackApp generate user comeback app structure with user and tparty identifier
func NewUserComebackApp(userID, tPartyID int64) *UserComebackApp {
	return &UserComebackApp{UserID: userID, TPartyID: tPartyID}
}

// TableName user comeback app database
func (m UserComebackApp) TableName() string {
	return "user_comeback_apps"
}

// ToJSON user comeback app structure to json string
func (m UserComebackApp) ToJSON() string {
	return database.ToJSON(m)
}
