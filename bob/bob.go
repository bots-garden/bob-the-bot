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
