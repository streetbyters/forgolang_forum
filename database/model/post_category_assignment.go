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
	"gopkg.in/guregu/null.v3/zero"
	"time"
)

// PostCategoryAssignment discussion topic category assignment
type PostCategoryAssignment struct {
	database.DBInterface `json:"-"`
	ID                   int64     `db:"id" json:"id"`
	PostID               int64     `db:"post_id" json:"post_id" foreign:"fk_post_category_assignments_post_id" unique:"post_category_assignments_post_category_unique" validate:"required"`
	CategoryID           int64     `db:"category_id" json:"category_id" foreign:"fk_post_category_assignments_category_id" unique:"post_category_assignments_post_category_unique" validate:"required"`
	SourceUserID         zero.Int  `db:"source_user_id" json:"source_user_id" foreign:"fk_post_category_assignments_source_user_id"`
	InsertedAt           time.Time `db:"inserted_at" json:"inserted_at"`
}

// NewPostCategoryAssignment generate post category assignment structure
func NewPostCategoryAssignment(postID,
	categoryID, sourceUserID int64) *PostCategoryAssignment {
	return &PostCategoryAssignment{PostID: postID,
		CategoryID:   categoryID,
		SourceUserID: zero.IntFrom(sourceUserID)}
}

// TableName post category assignment database
func (m PostCategoryAssignment) TableName() string {
	return "post_category_assignments"
}

// ToJSON post category assignment structure to json string
func (m PostCategoryAssignment) ToJSON() string {
	return database.ToJSON(m)
}
