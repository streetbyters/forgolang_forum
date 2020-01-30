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

// PostTag special classifications for posts
type PostTag struct {
	database.DBInterface `json:"-"`
	ID                   int64     `db:"id" json:"id"`
	PostID               int64     `db:"post_id" json:"post_id" foreign:"fk_post_tags_post_id" unique:"post_tags_post_tag_unique" validate:"required"`
	TagID                int64     `db:"tag_id" json:"tag_id" foreign:"fk_post_tags_tag_id" unique:"post_tags_post_tag_unique" validate:"required"`
	SourceUserID         zero.Int  `db:"source_user_id" json:"source_user_id,omitempty"`
	InsertedAt           time.Time `db:"inserted_at" json:"inserted_at"`
}

// NewPostTag generate post tag structure with post identifier
func NewPostTag(postID int64) *PostTag {
	return &PostTag{PostID: postID}
}

// TableName post tag database
func (m PostTag) TableName() string {
	return "post_tags"
}

// ToJSON post tag structure to json string
func (m PostTag) ToJSON() string {
	return database.ToJSON(m)
}
