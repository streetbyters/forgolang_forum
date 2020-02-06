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
	"time"
)

// PostVotesUp Users' ratings on discussion topics
type PostVotesUp struct {
	database.DBInterface `json:"-"`
	ID                   int64     `db:"id" json:"id"`
	PostID               int64     `db:"post_id" json:"post_id" foreign:"fk_post_votes_up_post_id" unique:"post_votes_up_post_user_unique" validate:"required"`
	UserID               int64     `db:"user_id" json:"user_id" foreign:"fk_post_votes_up_user_id" unique:"post_votes_up_post_user_unique" validate:"required"`
	InsertedAt           time.Time `db:"inserted_at" json:"inserted_at"`
}

// NewPostVotesUp generate post votes up structure
func NewPostVotesUp(postID, userID int64) *PostVotesUp {
	return &PostVotesUp{PostID: postID, UserID: userID}
}

// TableName post votes up database
func (m PostVotesUp) TableName() string {
	return "post_votes_up"
}

// ToJSON post votes up structure to json string
func (m PostVotesUp) ToJSON() string {
	return database.ToJSON(m)
}
