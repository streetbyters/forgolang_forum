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

// ThirdParty integrated third-party systems structure
type ThirdParty struct {
	database.DBInterface `json:"-"`
	ID                   int64           `db:"id" json:"-"`
	Name                 string          `db:"name" json:"name" validate:"required,gte=2,lte=200"`
	Code                 string          `db:"code" json:"code" unique:"third_party_code_unique" validate:"required"`
	Type                 database.TParty `db:"type" json:"type" validate:"required"`
	IsActive             bool            `db:"is_active" json:"is_active"`
	InsertedAt           time.Time       `db:"inserted_at" json:"inserted_at"`
	UpdatedAt            time.Time       `db:"updated_at" json:"updated_at"`
}

// NewThirdParty generate third-party structure
func NewThirdParty() *ThirdParty {
	return &ThirdParty{IsActive: true}
}

// TableName third-party database
func (m ThirdParty) TableName() string {
	return "third_party"
}

// ToJSON third-party structure to json string
func (m ThirdParty) ToJSON() string {
	return database.ToJSON(m)
}
