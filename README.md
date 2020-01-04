 <h1 align="center">Forgolang.com</h1>
 <p align="center">
   <a href="https://travis-ci.org/akdilsiz/forgolang_forum">
    <img src="https://travis-ci.org/akdilsiz/forgolang_forum.svg?branch=master"/>
   </a>
   <a href="https://github.com/akdilsiz/forgolang_forum/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/forgolang_forum/agente"/>
   </a>
   <a href="https://codecov.io/gh/akdilsiz/forgolang_forum">
     <img src="https://codecov.io/gh/akdilsiz/forgolang_forum/branch/master/graph/badge.svg" />
   </a>
   <a href="https://goreportcard.com/report/github.com/akdilsiz/forgolang_forum">
    <img src="https://goreportcard.com/badge/github.com/akdilsiz/forgolang_forum"/>
   </a>
 </p>

Forgolang.com's open-source forum system

## Requirements
 - Go > 1.11
 - PostgreSQL
 - Redis 
 - RabbitMQ

## Development
```shell script
git clone -b develop https://github.com/akdilsiz/forgolang_forum

go mod vendor

go run ./cmd -genSecretEnv

# Development Mode
go run ./cmd -mode dev -migrate -reset
go run ./cmd -mode dev

# Test Mode
go run ./cmd -mode test -migrate -reset
go run ./cmd -mode test
```

## Contribution
I would like to accept any contributions to make Forgolang.com better and feature rich.\
[See detail](docs/contributions.md)

## LICENSE

Copyright 2019 Abdulkadir DILSIZ

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
