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
	"github.com/akdilsiz/agente/utils"
	"gopkg.in/guregu/null.v3/zero"
	"time"
)

// User Authentication/authorization base database model
type User struct {
	database.DBInterface `json:"-"`
	ID                   int64       `db:"id" json:"id"`
	Username             string      `db:"username" json:"username" unique:"users_username_unique_index" validate:"required"`
	PasswordDigest       zero.String `db:"password_digest" json:"-"`
	Password             string      `db:"-" json:"password" validate:"required"`
	Email                string      `db:"email" json:"email" unique:"users_email_unique_index" validate:"required,email"`
	IsActive             bool        `db:"is_active" json:"is_active"`
	Avatar               zero.String `db:"avatar" json:"avatar"`
	InsertedAt           time.Time   `db:"inserted_at" json:"inserted_at"`
	UpdatedAt            time.Time   `db:"updated_at" json:"updated_at"`
}

// NewUser user generate with default data
func NewUser(pwd *string) *User {
	if pwd != nil {
		return &User{
			PasswordDigest: zero.StringFrom(utils.HashPassword(*pwd, 11)),
			IsActive:       true,
		}
	}

	return &User{IsActive: false}
}

// TableName user database table name
func (d User) TableName() string {
	return "users"
}

// ToJSON User database model to json string
func (d User) ToJSON() string {
	return database.ToJSON(d)
}

// Timestamps generate timestamp fields
func (d User) Timestamps() bool {
	return true
}
