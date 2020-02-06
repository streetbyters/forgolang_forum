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
	"forgolang_forum/model"
	"gopkg.in/guregu/null.v3/zero"
	"time"
)

// Language Base Internationalization(I18N) / Localization(L10N) structure
type Language struct {
	model.Model `json:"-"`
	ID          int64       `db:"id" json:"id"`
	Name        string      `db:"name" json:"name" validate:"required,gte=2,lte=64"`
	Code        string      `db:"code" json:"code" unique:"languages_code" validate:"required,gte=2,lte=10"`
	DateFormat  zero.String `db:"date_format" json:"date_format"`
	InsertedAt  time.Time   `db:"inserted_at" json:"inserted_at"`
	UpdatedAt   time.Time   `db:"updated_at" json:"updated_at"`
}

// NewLanguage generate language structure
func NewLanguage() *Language {
	return &Language{}
}

// TableName language database
func (m Language) TableName() string {
	return "languages"
}

// ToJSON tag structure to json string
func (m Language) ToJSON() string {
	return database.ToJSON(m)
}

// Timestamps generate timestamps fields
func (m Language) Timestamps() bool {
	return true
}

