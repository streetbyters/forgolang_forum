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

// PostComment Users comments on discussion topics
type PostComment struct {
	database.DBInterface `json:"-"`
	ID                   int64     `db:"id" json:"id"`
	PostID               int64     `db:"post_id" json:"post_id" foreign:"fk_post_comments_post_id" validate:"required"`
	UserID               int64     `db:"user_id" json:"user_id" foreign:"fk_post_comments_user_id" validate:"required"`
	ParentID             zero.Int  `db:"parent_id" json:"parent_id" foreign:"fk_post_comments_parent_id"`
	InsertedAt           time.Time `db:"inserted_at" json:"inserted_at"`
}

// NewPostComment generate post comment structure
func NewPostComment(postID, userID int64) *PostComment {
	return &PostComment{PostID: postID, UserID: userID}
}

// TableName post comments database
func (m PostComment) TableName() string {
	return "post_comments"
}

// ToJSON post comments structure to json string
func (m PostComment) ToJSON() string {
	return database.ToJSON(m)
}
