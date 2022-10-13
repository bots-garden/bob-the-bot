package main

import (
	"strconv"
	"strings"

	hf "github.com/bots-garden/capsule/capsulemodule/hostfunctions"
	"github.com/tidwall/gjson"
)

const (
	Open   string = "open"   // new issue
	Close         = "close"  // issue or comment
	Reopen        = "reopen" // issue
	Update        = "update" // issue or comment
)

const (
	Issue = "issue"
	Note  = "note"
)

func main() {
	hf.SetHandleHttp(Handle)
}

func addNoteToTheIssue(issueIid float64, projectId float64, message string) (string, error) {

	botToken, _ := hf.GetEnv("BOT_TOKEN")
	apiUrl, _ := hf.GetEnv("API_URL")

	issueNumber := strconv.FormatInt(int64(issueIid), 10)
	projectNumber := strconv.FormatInt(int64(projectId), 10)

	hf.Log("🖐 issueNumber: " + issueNumber)
	hf.Log("🖐 projectNumber: " + projectNumber)

	noteApiUrl := apiUrl + "/projects/" + projectNumber + "/issues/" + issueNumber + "/notes"

	headers := map[string]string{
		"Content-Type":  "application/json; charset=utf-8",
		"PRIVATE-TOKEN": botToken,
	}

	jsondoc := `{"body": "` + message + `"}`

	return hf.Http(noteApiUrl, "POST", headers, jsondoc)
}

func Handle(request hf.Request) (response hf.Response, errResp error) {
	// project: https://gitlab.com/k33g_org/bob-the-bot
	var resp string
	var err error

	object := gjson.Get(request.Body, "object_kind") // it should be issue or note
	objectAttributes := gjson.Get(request.Body, "object_attributes") // issue or comment
	projectId := gjson.Get(request.Body, "project.id").Num
	userName := gjson.Get(request.Body, "user.username").Str 

	botName, _ := hf.GetEnv("BOT_NAME")

	if object.Str == Issue {

		action := objectAttributes.Get("action") // open close reopen update (only for issue)

		if action.Str == Open || action.Str == Update {
			issueDescription := objectAttributes.Get("description").Str
			issueIid := objectAttributes.Get("iid").Num

			if strings.Contains(issueDescription, botName) {
				resp, err = addNoteToTheIssue(issueIid, projectId, "👋 @"+userName+" what's up? 😄")
			}
		}
	}

	// a comment is added to the issue or an existing comment is updated
	if object.Str == Note {

		note := objectAttributes.Get("note").Str
		issueIid := gjson.Get(request.Body, "issue.iid").Num

		if strings.Contains(note, botName) {
			resp, err = addNoteToTheIssue(issueIid, projectId, "🤔 @"+userName+" are you talking to me?")
		}
	}

	headersResp := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	if err != nil {
		hf.Log("😡 error:" + err.Error())
	}

	return hf.Response{Body: resp, Headers: headersResp}, err
}
