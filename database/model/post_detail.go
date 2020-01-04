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
	"gopkg.in/guregu/null.v3/zero"
	"time"
)

// PostDetail discussion topic content details
type PostDetail struct {
	database.DBInterface `json:"-"`
	ID                   int64       `db:"id" json:"id"`
	PostID               int64       `db:"post_id" json:"post_id" foreign:"fk_post_details_post_id" validate:"required"`
	SourceUserID         zero.Int    `db:"source_user_id" json:"source_user_id" foreign:"fk_post_details_source_user_id"`
	Title                string      `db:"title" json:"title" validate:"required,gte=3,lte=200"`
	Description          zero.String `db:"description" json:"description"`
	Content              string      `db:"content" json:"content" validate:"required" validate:"gte=5,lte=10240"`
	InsertedAt           time.Time   `db:"inserted_at" json:"inserted_at"`
}

// NewPostDetail generate post detail struct
func NewPostDetail(postID, sourceUserID int64) *PostDetail {
	return &PostDetail{PostID: postID, SourceUserID: zero.IntFrom(sourceUserID)}
}

// TableName post detail database
func (m PostDetail) TableName() string {
	return "post_details"
}

// ToJSON post detail structure to json string
func (m PostDetail) ToJSON() string {
	return database.ToJSON(m)
}
