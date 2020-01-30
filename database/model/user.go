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
	"fmt"
	"forgolang_forum/database"
	"forgolang_forum/utils"
	"gopkg.in/guregu/null.v3/zero"
	"time"
)

// User Authentication/authorization base database model
type User struct {
	database.DBInterface `json:"-"`
	ID                   int64       `db:"id" json:"id"`
	Username             string      `db:"username" json:"username" unique:"users_username_unique_index" validate:"required"`
	PasswordDigest       zero.String `db:"password_digest" json:"-"`
	Password             string      `json:"password"`
	Email                string      `db:"email" json:"email" unique:"users_email_unique_index" validate:"required,email"`
	EmailHidden          bool        `db:"email_hidden" json:"email_hidden"`
	Bio                  zero.String `db:"bio" json:"bio" validate:"lte=10240"`
	Url                  zero.String `db:"url" json:"url" validate:"lte=200"`
	IsActive             bool        `db:"is_active" json:"is_active"`
	Avatar               zero.String `db:"avatar" json:"avatar"`
	Role                 zero.String `db:"role" json:"role,omitempty"`
	RoleAssignmentID     zero.Int    `db:"role_assignment_id" json:"role_assignment_id,omitempty"`
	State                zero.String `db:"state" json:"state,omitempty"`
	TPartyName           zero.String `db:"tparty_name" json:"tparty_name,omitempty"`
	TPartyData           zero.String `db:"tparty_data" json:"tparty_data,omitempty"`
	InsertedAt           time.Time   `db:"inserted_at" json:"inserted_at"`
	UpdatedAt            time.Time   `db:"updated_at" json:"updated_at"`
}

// NewUser user generate with default data
func NewUser(pwd *string) *User {
	if pwd != nil {
		return &User{
			PasswordDigest: zero.StringFrom(utils.HashPassword(*pwd, 11)),
			EmailHidden:    true,
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

// Query generate for user
func (d User) Query(force bool) string {
	roleAssignment := new(UserRoleAssignment)
	role := new(Role)
	userState := new(UserState)
	userComeBack := new(UserComebackApp)
	thirdParty := new(ThirdParty)

	var query string

	query = fmt.Sprintf(`
		SELECT 
			u.id as id,
			u.username as username,
			u.email as email,
			u.is_active as is_active,
			u.avatar as avatar,
			r.code as role,
			ra.id as role_assignment_id,
			us.state as state,
			tp.name as tparty_name,
			uca.data as tparty_data,
			u.inserted_at as inserted_at,
			u.updated_at as updated_at
		FROM %s AS u
	`, d.TableName())

	if force {
		return fmt.Sprintf(`%s
			INNER OUTER JOIN %s AS ra ON u.id = ra.user_id
			LEFT OUTER JOIN %s AS ra2 ON ra.user_id = ra2.user_id and ra.id < ra2.id
			LEFT OUTER JOIN %s AS r on ra.role_id = r.id
			LEFT OUTER JOIN %s AS us ON u.id = us.user_id
			INNER OUTER JOIN %s AS us2 ON us.user_id = us2.user_id and us.id < us2.id
			LEFT OUTER JOIN %s AS uca ON u.id = uca.user_id
			LEFT OUTER JOIN %s AS tp ON uca.tparty_id = tp.id
			WHERE ra2.id IS NULL AND us2.id IS NULL
		`, query,
			roleAssignment.TableName(),
			roleAssignment.TableName(),
			role.TableName(),
			userState.TableName(),
			userState.TableName(),
			userComeBack.TableName(),
			thirdParty.TableName())
	}

	return fmt.Sprintf(`%s
			LEFT OUTER JOIN %s AS ra ON u.id = ra.user_id
			LEFT OUTER JOIN %s AS ra2 ON ra.user_id = ra2.user_id and ra.id < ra2.id
			LEFT OUTER JOIN %s AS r on ra.role_id = r.id
			LEFT OUTER JOIN %s AS us ON u.id = us.user_id
			LEFT OUTER JOIN %s AS us2 ON us.user_id = us2.user_id and us.id < us2.id
			LEFT OUTER JOIN %s AS uca ON u.id = uca.user_id
			LEFT OUTER JOIN %s AS tp ON uca.tparty_id = tp.id
			WHERE ra2.id IS NULL AND us2.id IS NULL
		`, query,
		roleAssignment.TableName(),
		roleAssignment.TableName(),
		role.TableName(),
		userState.TableName(),
		userState.TableName(),
		userComeBack.TableName(),
		thirdParty.TableName())
}
