// Copyright 2019 Abdulkadir DILSIZ
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

package api

import (
	"fmt"
	"forgolang_forum/model"
	"github.com/valyala/fasthttp"
	"path/filepath"
)

// UploadController file upload api controller
type UploadController struct {
	Controller
	*API
}

// Create file upload method
func (c UploadController) Create(ctx *fasthttp.RequestCtx) {
	file, err := ctx.FormFile("file")
	if err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: nil,
			Detail: err.Error(),
		}, fasthttp.StatusBadRequest)
		return
	}

	dir := ctx.FormValue("dir")
	if string(dir) == "" {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: map[string]string{
				"dir": "is not nil",
			},
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	err = c.App.Storage.Upload(file, filepath.Join(string(dir), file.Filename), "public-read")
	if err != nil {
		fmt.Println(err)
	}

	resp := make(map[string]string)
	resp["filename"] = file.Filename
	resp["dir"] = string(dir)

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: resp,
	}, fasthttp.StatusCreated)
}
