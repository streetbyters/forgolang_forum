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

// UserRoleAssignmentInvalidation user authorization invalidation structure
type UserRoleAssignmentInvalidation struct {
	database.DBInterface `json:"-"`
	AssignmentID         int64     `db:"assignment_id" json:"assignment_id" foreign:"fk_user_role_assignment_invalidations_assignment_id" unique:"user_role_assignment_invalidation_pkey" validate:"required"`
	SourceUserID         zero.Int  `db:"source_user_id" json:"source_user_id" foreign:"fk_user_role_assignment_invalidations_source_user_id"`
	InsertedAt           time.Time `db:"inserted_at" json:"inserted_at"`
}

// NewUserRoleAssignmentInvalidation generate structure
func NewUserRoleAssignmentInvalidation(assignmentID int64) *UserRoleAssignmentInvalidation {
	return &UserRoleAssignmentInvalidation{AssignmentID: assignmentID}
}

// TableName user role assignment invalidation database
func (m UserRoleAssignmentInvalidation) TableName() string {
	return "user_role_assignment_invalidations"
}

// ToJSON user role assignment structure to json string
func (m UserRoleAssignmentInvalidation) ToJSON() string {
	return database.ToJSON(m)
}
