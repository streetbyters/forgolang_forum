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
	"gopkg.in/guregu/null.v3/zero"
	"time"
)

// UserRoleAssignment users authorization artifacts
type UserRoleAssignment struct {
	database.DBInterface `json:"-"`
	ID                   int64     `db:"id" json:"id"`
	UserID               int64     `db:"user_id" json:"user_id" foreign:"fk_user_role_assignments_user_id" validate:"required"`
	RoleID               int64     `db:"role_id" json:"role_id" foreign:"fk_user_role_assignments_role_id" validate:"required"`
	SourceUserID         zero.Int  `db:"source_user_id" foreign:"fk_user_role_assignments_source_user_id" json:"source_user_id"`
	InsertedAt           time.Time `db:"inserted_at" json:"inserted_at"`
}

// NewUserRoleAssignment generate user role assignment structure
func NewUserRoleAssignment(userID, roleID int64) *UserRoleAssignment {
	return &UserRoleAssignment{UserID: userID, RoleID: roleID}
}

// TableName user role assignment database
func (m UserRoleAssignment) TableName() string {
	return "user_role_assignments"
}

// ToJSON user role assignment structure to json string
func (m UserRoleAssignment) ToJSON() string {
	return database.ToJSON(m)
}
