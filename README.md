# Bob the Bot

## Steps

### 01) Open this project with Gitpod

> this is a wip 🚧

When the project is opened with Gitpod, all dependencies are installed (see `.gitpod.yml`)

### 02) Generate the project

```bash
cabu generate service-post bob
cd bob
go mod tidy
```
> when **cabu** is run for the first time, it will pull the Capsule Builder Docker image

#### Source code generated

> bob.go
```golang
package main

import (
	hf "github.com/bots-garden/capsule/capsulemodule/hostfunctions"
	/* string to json */
	"github.com/tidwall/gjson"
	/* create json string */
	"github.com/tidwall/sjson"
)

func main() {
	hf.SetHandleHttp(Handle)
}

func Handle(request hf.Request) (response hf.Response, errResp error) {
	
	hf.Log("📝 Body: " + request.Body)
	hf.Log("📝 URI: " + request.Uri)
	hf.Log("📝 Method: " + request.Method)

	name := gjson.Get(request.Body, "name")
	
	headersResp := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	jsondoc := `{"message": ""}`
	jsondoc, _ = sjson.Set(jsondoc, "message", "👋 hello " + name.Str)

	return hf.Response{Body: jsondoc, Headers: headersResp}, nil
}

/*
curl -v -X POST   http://localhost:8080   -H 'content-type: application/json'   -d '{"name": "Bob"}'
*/
```

#### First build

```bash
cd bob
cabu build . bob.go bob.wasm
```

#### First Run

```bash
cd bob
capsule -wasm=./bob.wasm -mode=http -httpPort=8080
```

Call the service: 

```bash
curl -v -X POST http://localhost:8080 -H 'content-type: application/json' -d '{"name": "Bob"}'
```

### 03) Prepare WebHook

- Go to Settings > Webhooks (Content type: application/json) / Issues + Issue comments

https://8080-botsgarden-bobthebot-3uk65iyrzav.ws-eu70.gitpod.io