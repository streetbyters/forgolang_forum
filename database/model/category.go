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

// Category base post artifacts
type Category struct {
	database.DBInterface `json:"-"`
	ID                   int64       `db:"id" json:"id"`
	Title                string      `db:"title" json:"title" validate:"required,gte=3,lte=128"`
	Description          zero.String `db:"description" json:"description" validate:"lte=240"`
	Slug                 string      `db:"slug" json:"slug" unique:"categories_slug_unique"`
	InsertedAt           time.Time   `db:"inserted_at" json:"inserted_at"`
	UpdatedAt            time.Time   `db:"updated_at" json:"updated_at"`
}

// NewCategory generate category struct
func NewCategory() *Category {
	return &Category{}
}

// TableName category database
func (m Category) TableName() string {
	return "categories"
}

// ToJSON category struct to json string
func (m Category) ToJSON() string {
	return database.ToJSON(m)
}
