// Copyright 2019 StreetByters Community
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

// Post discussion topics created by users
type Post struct {
	database.DBInterface `json:"-"`
	ID                   int64     `db:"id" json:"id"`
	AuthorID             int64     `db:"author_id" json:"author_id" foreign:"fk_posts_author_id" validate:"required"`
	InsertedAt           time.Time `db:"inserted_at" json:"inserted_at"`
}

// NewPost generate post struct
func NewPost(authorID int64) *Post {
	return &Post{AuthorID: authorID}
}

// TableName post database
func (m Post) TableName() string {
	return "posts"
}

// ToJSON post structure to json string
func (m Post) ToJSON() string {
	return database.ToJSON(m)
}

// PostDEP all post fields
type PostDEP struct {
	database.DBInterface `json:"-"`
	ID                   int64                     `db:"id" json:"id"`
	AuthorID             int64                     `db:"author_id" json:"author_id" validate:"required"`
	AuthorUsername       zero.String               `db:"author_username" json:"author_username,omitempty"`
	Slug                 zero.String               `db:"slug" json:"slug,omitempty"`
	Title                zero.String               `db:"title" json:"title,omitempty" validate:"required,gte=3,lte=120"`
	Description          zero.String               `db:"description" json:"description,omitempty"`
	Content              zero.String               `db:"content" json:"content,omitempty" validate:"required,gte=20,lte=10240"`
	CategoryAssignments  *[]PostCategoryAssignment `db:"category_assignments" json:"category_assignments,omitempty"`
	InsertedAt           time.Time                 `db:"inserted_at" json:"inserted_at"`
}
