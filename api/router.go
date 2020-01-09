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

package api

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"forgolang_forum/model"
	errors2 "github.com/akdilsiz/agente/errors"
	"github.com/akdilsiz/limiterphi"
	"github.com/fate-lovely/phi"
	"github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// Router api router structure
type Router struct {
	API     *API
	Server  *fasthttp.Server
	Addr    string
	Handler *phi.Mux
	Routes  map[string]map[string][]string
}

var (
	prefix           string
	reqID            uint64
	allowHeaders     = "authorization"
	allowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	allowOrigin      = "*"
	allowCredentials = "true"
)

// NewRouter building api router
func NewRouter(api *API) *Router {
	router := &Router{
		API: api,
	}
	router.Routes = make(map[string]map[string][]string)

	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])

	rate, _ := limiter.NewRateFromFormatted("50-S")
	store, err := sredis.NewStoreWithOptions(api.App.Cache, limiter.StoreOptions{
		Prefix:   "forgolang.com",
		MaxRetry: 4,
	})
	if err != nil {
		panic(err)
	}
	rateMiddleware := limiterphi.NewMiddleware(limiter.New(store, rate))

	r := phi.NewRouter()

	r.Use(router.requestID)
	r.Use(router.recover)
	r.Use(router.logger)
	r.Use(router.cors)
	r.Use(rateMiddleware.Handle)

	r.NotFound(router.notFound)
	r.MethodNotAllowed(router.methodNotAllowed)

	hC := HomeController{API: api}
	r.Get("/", hC.Index)

	routerPrefix := strings.Join([]string{api.App.Config.Prefix, "v1"}, "/")

	r.Route(routerPrefix, func(r phi.Router) {
		// Auth routes
		r.Route("/auth", func(r phi.Router) {
			r.Post("/sign_in", LoginController{API: api}.Create)
			r.Post("/token", TokenController{API: api}.Create)
			r.Post("/register", RegisterController{API: api}.Create)
			r.Post("/confirmation/{userID}/{code}", ConfirmationController{API: api}.Create)

			// Third-party routes
			r.Get("/github", AuthController{API: api}.Github)
			r.Get("/github/callback", AuthController{API: api}.GithubCallback)
		})

		r.Group(func(r phi.Router) {
			r.Use(api.JWTAuth.Verify)

			uC := UploadController{API: api}

			r.Post("/upload", uC.Create)
			router.Routes["UploadController"] = make(map[string][]string)
			router.Routes["UploadController"]["superadmin"] = []string{
				"Create",
			}

			//User Routes
			r.Group(func(r phi.Router) {
				uC := UserController{API: api}
				r.With(UserPolicy{API: api}.Index).Get("/user", uC.Index)
				r.With(UserPolicy{API: api}.Create).Post("/user", uC.Create)
				r.Route("/user/{userID}", func(r phi.Router) {
					r.With(UserPolicy{API: api}.Show).Get("/", uC.Show)
					r.With(UserPolicy{API: api}.Update).Put("/", uC.Update)
					r.With(UserPolicy{API: api}.Delete).Delete("/", uC.Delete)

					// Role assignment routes
					r.With(UserRoleAssignmentPolicy{API: api}.Create).
						Post("/role_assignment", UserRoleAssignmentController{API: api}.Create)
					router.Routes["UserRoleAssignmentController"] = make(map[string][]string)
					router.Routes["UserRoleAssignmentController"]["superadmin"] = []string{
						"Create",
					}
				})
				router.Routes["UserController"] = make(map[string][]string)
				router.Routes["UserController"]["superadmin"] = []string{
					"Index",
					"Show",
					"Create",
					"Update",
					"Delete",
				}
				router.Routes["UserController"]["user"] = []string{
					"Show",
					"Update",
				}
			})

			// Category routes
			r.Group(func(r phi.Router) {
				cC := CategoryController{API: api}
				r.Get("/category", cC.Index)
				r.With(CategoryPolicy{API: api}.Create).Post("/category", cC.Create)
				r.Route("/category/{categoryID}", func(r phi.Router) {
					r.Get("/", cC.Show)
					r.With(CategoryPolicy{API: api}.Update).Put("/", cC.Update)
					r.With(CategoryPolicy{API: api}.Delete).Delete("/", cC.Delete)
				})
				router.Routes["CategoryController"] = make(map[string][]string)
				router.Routes["CategoryController"]["superadmin"] = []string{
					"Index",
					"Create",
					"Show",
					"Update",
					"Delete",
				}
				router.Routes["CategoryController"]["moderator"] = []string{
					"Index",
					"Show",
					"Update",
				}
				router.Routes["CategoryController"]["user"] = []string{
					"Index",
					"Show",
				}
			})
		})
	})

	router.Server = &fasthttp.Server{
		Handler:            r.ServeFastHTTP,
		ReadTimeout:        10 * time.Second,
		MaxRequestBodySize: 1 * 1024 * 1024 * 1024,
		Logger:             api.App.Logger,
	}
	router.Addr = fmt.Sprintf("%s:%d", api.App.Config.Host, api.App.Config.Port)
	router.Handler = r

	return router
}

func (r Router) notFound(ctx *fasthttp.RequestCtx) {
	r.API.JSONResponse(ctx, model.ResponseError{
		Errors: nil,
		Detail: "not found",
	}, http.StatusNotFound)
}

func (r Router) methodNotAllowed(ctx *fasthttp.RequestCtx) {
	r.API.JSONResponse(ctx, model.ResponseError{
		Errors: nil,
		Detail: "method not allowed",
	}, http.StatusMethodNotAllowed)
}

// Reference: https://github.com/go-chi/chi/blob/master/middleware/request_id.go
func (r Router) requestID(next phi.HandlerFunc) phi.HandlerFunc {
	return func(ctx *fasthttp.RequestCtx) {
		id := atomic.AddUint64(&reqID, 1)
		requestID := fmt.Sprintf("%s-%06d", prefix, id)
		ctx.SetUserValue("requestID", requestID)
		ctx.Response.Header.Add("x-request-id", requestID)
		next(ctx)
	}
}

// Reference: https://github.com/go-chi/chi/blob/master/middleware/recoverer.go
func (r Router) recover(next phi.HandlerFunc) phi.HandlerFunc {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if rvr := recover(); rvr != nil {
				var err error
				switch x := rvr.(type) {
				case *errors2.PluggableError:
					e := rvr.(*errors2.PluggableError)
					r.API.JSONResponse(ctx, model.ResponseError{
						Errors: e.Errors,
						Detail: e.Error(),
					}, e.Status)

					defer func() {
						r.API.App.Logger.LogError(e, "Pluggable error")
					}()
					return
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("unknown panic")
				}

				if r.API.App.Mode == model.Test {
					panic(rvr)
				}

				r.API.App.Logger.LogError(err, "router recover")
				r.API.JSONResponse(ctx, model.ResponseError{
					Errors: nil,
					Detail: http.StatusText(http.StatusInternalServerError),
				}, http.StatusInternalServerError)
				return
			}
		}()

		next(ctx)
	}
}

func (r Router) logger(next phi.HandlerFunc) phi.HandlerFunc {
	return func(ctx *fasthttp.RequestCtx) {
		next(ctx)
		defer func() {
			if r.API.App.Mode != model.Test {
				if r.API.App.Mode == model.Prod {
					r.API.App.Logger.LogInfo("Path: " + string(ctx.Path()) +
						" Method: " + string(ctx.Method()) +
						" - " + strconv.Itoa(ctx.Response.StatusCode()))
				} else {
					r.API.App.Logger.LogDebug("Path: " + string(ctx.Path()) +
						" Method: " + string(ctx.Method()) +
						" - " + strconv.Itoa(ctx.Response.StatusCode()))
				}
			}
		}()
	}
}

func (r Router) cors(next phi.HandlerFunc) phi.HandlerFunc {
	return func(ctx *fasthttp.RequestCtx) {
		if string(ctx.Request.Header.Method()) == "OPTIONS" {
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", allowCredentials)
			ctx.Response.Header.Set("Access-Control-Allow-Headers", allowHeaders)
			ctx.Response.Header.Set("Access-Control-Allow-Methods", allowMethods)
			ctx.Response.Header.Set("Access-Control-Allow-Origin", allowOrigin)
			ctx.Response.Header.Set("Accept", "application/json")
			ctx.Response.Header.Set("Accept", "multipart/form-data")

			ctx.SetStatusCode(http.StatusNoContent)
			return
		}
		next(ctx)
	}
}
