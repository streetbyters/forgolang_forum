// Copyright 2019 Forgolang Community
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

// PostSlug post seo links
type PostSlug struct {
	database.DBInterface `json:"-"`
	ID                   int64     `db:"id" json:"id"`
	PostID               int64     `db:"post_id" json:"post_id" foreign:"fk_post_slugs_post_id" validate:"required"`
	SourceUserID         zero.Int  `db:"source_user_id" json:"source_user_id" foreign:"fk_post_slugs_source_user_id"`
	Slug                 string    `db:"slug" json:"slug" unique:"post_slugs_slug_index" validate:"required"`
	InsertedAt           time.Time `db:"inserted_at" json:"inserted_at"`
}

// NewPostSlug generate post slug structure
func NewPostSlug(postID, sourceUserID int64) *PostSlug {
	return &PostSlug{PostID: postID, SourceUserID: zero.IntFrom(sourceUserID)}
}

// TableName post slug database
func (m PostSlug) TableName() string {
	return "post_slugs"
}

// ToJSON post slug structure to json string
func (m PostSlug) ToJSON() string {
	return database.ToJSON(m)
}
