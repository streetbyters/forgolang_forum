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
	"forgolang_forum/model"
	"time"
)

// TagLanguage tag localization(L10N) structure
type TagLanguage struct {
	model.Model  `json:"-"`
	ID           int64     `db:"id" json:"id"`
	TagID        int64     `db:"tag_id" json:"tag_id" foreign:"fk_tag_languages_tag_id" validate:"required"`
	LanguageID   int64     `db:"language_id" json:"language_id" foreign:"fk_tag_languages_language_id" validate:"required"`
	SourceUserID int64     `db:"source_user_id" json:"source_user_id" foreign:"fk_tag_languages_source_user_id"`
	Name         string    `db:"name" json:"name" validate:"required,gte=3,lte=32"`
	InsertedAt   time.Time `db:"inserted_at" json:"inserted_at"`
}

// NewTagLanguage generate tag language structure with tag identifier
func NewTagLanguage(tagID int64) *TagLanguage {
	return &TagLanguage{TagID: tagID}
}

// TableName tag language database
func (m TagLanguage) TableName() string {
	return "tag_languages"
}

// ToJSON tag language structure to json string
func (m TagLanguage) ToJSON() string {
	return database.ToJSON(m)
}
