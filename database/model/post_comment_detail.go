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
	"time"
)

// PostCommentDetail Users comment detail on discussion topics
type PostCommentDetail struct {
	database.DBInterface `json:"-"`
	ID                   int64     `db:"id" json:"id"`
	CommentID            int64     `db:"comment_id" json:"comment_id" foreign:"fk_post_comment_details_comment_id" validate:"required"`
	Comment              string    `db:"comment" json:"comment" validate:"required,gte=5,lte=10240"`
	InsertedAt           time.Time `db:"inserted_at" json:"inserted_at"`
}

// NewPostCommentDetail generate post comment detail structure
func NewPostCommentDetail(commentID int64) *PostCommentDetail {
	return &PostCommentDetail{CommentID: commentID}
}

// TableName post comment detail database
func (m PostCommentDetail) TableName() string {
	return "post_comment_details"
}

// ToJSON post comment detail structure to json string
func (m PostCommentDetail) ToJSON() string {
	return database.ToJSON(m)
}
