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

package api

import (
	"forgolang_forum/model"
	"github.com/valyala/fasthttp"
	"time"
)

// HeartbeatController api status controller
type HeartbeartController struct {
	Controller
	*API
}

func (c HeartbeartController) Show(ctx *fasthttp.RequestCtx) {
	data := make(map[string]interface{})

	data["status"] = "OK"
	data["time_information"] = map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"zone": "UTC",
	}

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: data,
	}, fasthttp.StatusOK)
}
