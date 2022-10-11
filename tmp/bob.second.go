package main

import (
	"strconv"
	"strings"

	hf "github.com/bots-garden/capsule/capsulemodule/hostfunctions"
	"github.com/tidwall/gjson"
)

const (
	Opened  string = "opened"  // new issue
	Edited         = "edited"  // issue or comment
	Created        = "created" // new comment
)


func main() {
	hf.SetHandleHttp(Handle)
}

func addNoteToTheIssue(issue gjson.Result, message string) (string, error) {

	botToken, _ := hf.GetEnv("BOT_TOKEN")
	repoName, _ := hf.GetEnv("REPO_NAME")
	repoOwner, _ := hf.GetEnv("REPO_OWNER")
	apiUrl, _ := hf.GetEnv("API_URL")


	issue_number := issue.Get("number").Num
	issueNumber := strconv.FormatInt(int64(issue_number), 10)

	commentApiUrl := apiUrl + "/repos/" + repoOwner + "/" + repoName + "/issues/" + issueNumber + "/comments"

	headers := map[string]string{
		"Accept":        "application/vnd.github+json",
		"Authorization": "Bearer " + botToken,
	}

	jsondoc := `{"body": "` + message + `"}`

	return hf.Http(commentApiUrl, "POST", headers, jsondoc)
}

func Handle(request hf.Request) (response hf.Response, errResp error) {

	var resp string
	var err error

	action := gjson.Get(request.Body, "action")
	issue := gjson.Get(request.Body, "issue")
	comment := gjson.Get(request.Body, "comment")

	botName, _ := hf.GetEnv("BOT_NAME")

	// an issue is created or updated
	if action.Str == Opened || (action.Str == Edited && comment.Exists() == false) { 
		issue_user_login := issue.Get("user").Get("login").Str
		issue_body := issue.Get("body").Str

		if strings.Contains(issue_body, botName) {
			resp, err = addNoteToTheIssue(issue, "👋 @"+issue_user_login+" are you talking to me?")
		}
	}

	// a comment is added to the issue or an existing comment is updated
	if action.Str == Created || (action.Str == Edited && comment.Exists()) {
		comment := gjson.Get(request.Body, "comment")
		comment_user_login := comment.Get("user").Get("login").Str
		comment_body := comment.Get("body").Str

		if strings.Contains(comment_body, botName) {
			resp, err = addNoteToTheIssue(issue, "👋 @"+comment_user_login+" are you talking to me?")
		}
	}

	headersResp := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	if err != nil {
		hf.Log("😡 error:" + err.Error())
	} else {
		hf.Log("🙂 GitHub response: " + resp)
	}

	return hf.Response{Body: resp, Headers: headersResp}, err
}
