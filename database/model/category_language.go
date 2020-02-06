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
	"gopkg.in/guregu/null.v3/zero"
	"time"
)

// CategoryLanguage category localization(L10N) structure
type CategoryLanguage struct {
	model.Model  `json:"-"`
	ID           int64       `db:"id" json:"id"`
	CategoryID   int64       `db:"category_id" json:"category_id" foreign:"fk_category_languages_category_id" validate:"required"`
	LanguageID   int64       `db:"language_id" json:"language_id" foreign:"fk_category_languages_language_id" validate:"required"`
	SourceUserID zero.Int    `db:"source_user_id" json:"source_user_id" foreign:"fk_category_languages_source_user_id"`
	Title        string      `db:"title" json:"title" validate:"required,gte=3,lte=128"`
	Description  zero.String `db:"description" json:"description"`
	Slug         string      `db:"slug" json:"slug" validate:"required,gte=3,lte=200"`
	InsertedAt   time.Time   `db:"inserted_at" json:"inserted_at"`
}

// NewCategoryLanguage generate category language structure with category identifier
func NewCategoryLanguage(categoryID int64) *CategoryLanguage {
	return &CategoryLanguage{CategoryID: categoryID}
}

// TableName category language database
func (m CategoryLanguage) TableName() string {
	return "category_languages"
}

// ToJSON category language structure to json string
func (m CategoryLanguage) ToJSON() string {
	return database.ToJSON(m)
}
