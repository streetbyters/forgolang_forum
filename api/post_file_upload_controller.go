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

package api

import (
	"fmt"
	"forgolang_forum/database/model"
	model2 "forgolang_forum/model"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
	"path/filepath"
	"strconv"
)

type PostFileController struct {
	Controller
	*API
}

func (c PostFileController) Create(ctx *fasthttp.RequestCtx) {
	var post model.Post
	c.GetDB().QueryRowWithModel(fmt.Sprintf(`
		SELECT * FROM %s AS p
		WHERE p.id = $1
	`, post.TableName()),
		&post,
		phi.URLParam(ctx, "postID")).Force()

	file, err := ctx.FormFile("file")
	if err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: nil,
			Detail: err.Error(),
		}, fasthttp.StatusBadRequest)
		return
	}

	c.App.Storage.Upload(file,
		filepath.Join("posts", strconv.FormatInt(post.ID, 10), "files",
			file.Filename),
		"public-read")

	resp := make(map[string]string)
	resp["filename"] = file.Filename
	resp["path"] = filepath.Join("posts", strconv.FormatInt(post.ID, 10), "files")

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: resp,
	}, fasthttp.StatusCreated)
}
