package main

import (
	"strconv"
	"strings"

	hf "github.com/bots-garden/capsule/capsulemodule/hostfunctions"
	"github.com/tidwall/gjson"
)

/*
	issue creation: action == opened
	issue edition: action == edited
	issue comment: action == created
*/

const (
	Opened  string = "opened"  // new issue
	Edited         = "edited"  // issue or comment
	Created        = "created" // new comment
)

//export OnLoad
func OnLoad() {
	//botToken, _ := hf.GetEnv("BOT_TOKEN")
	botName, _ := hf.GetEnv("BOT_NAME")
	hf.Log("👋 👋" + botName)
}

func main() {
	hf.SetHandleHttp(Handle)
}

func displayComment(issue gjson.Result, comment gjson.Result) {
	comment_id := comment.Get("id")
	comment_body := comment.Get("body")
	comment_user := comment.Get("user")
	comment_user_login := comment_user.Get("login")

	botName, _ := hf.GetEnv("BOT_NAME")

	hf.Log("Comment:")
	hf.Log("  Id: " + strconv.FormatInt(int64(comment_id.Num), 10))
	hf.Log("  Body: " + comment_body.Str)
	hf.Log("  Login: " + comment_user_login.Str)

	if strings.Contains(comment_body.Str, botName) {
		hf.Log("👋 are you talking to me?")
		sendMessage(issue, "👋 are you talking to me?")
	}
}

func displayIssue(issue gjson.Result) {
	issue_id := issue.Get("id")
	issue_title := issue.Get("title")
	issue_body := issue.Get("body")
	issue_user := issue.Get("user")
	issue_user_login := issue_user.Get("login")

	botName, _ := hf.GetEnv("BOT_NAME")

	hf.Log("Issue: [Id: " + strconv.FormatInt(int64(issue_id.Num), 10) + "]")
	hf.Log("  Title: " + issue_title.Str)
	hf.Log("  Body: " + issue_body.Str)
	hf.Log("  User: " + issue_user_login.Str)

	if strings.Contains(issue_body.Str, botName) {
		hf.Log("👋 are you talking to me?")
		sendMessage(issue, "👋 are you talking to me?")
	}
}

func displayIssueTitle(issue gjson.Result) {
	issue_id := issue.Get("id")
	issue_title := issue.Get("title")
	hf.Log("Issue: " + issue_title.Str + " [Id: " + strconv.FormatInt(int64(issue_id.Num), 10) + "]")
}

func getIssueNumberAsString(issue gjson.Result) string {
	//issue_id := issue.Get("id")
	//return strconv.FormatInt(int64(issue_id.Num), 10)
	issue_number := issue.Get("number")
	return strconv.FormatInt(int64(issue_number.Num), 10)
}

/*
https://docs.github.com/en/rest/issues/comments#create-an-issue-comment

curl \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer <YOUR-TOKEN>" \
  https://api.github.com/repos/OWNER/REPO/issues/ISSUE_NUMBER/comments \
  -d '{"body":"Me too"}'
*/

func sendMessage(issue gjson.Result, message string) {

	botToken, _ := hf.GetEnv("BOT_TOKEN")
	repoName, _ := hf.GetEnv("REPO_NAME")
	repoOwner, _ := hf.GetEnv("REPO_OWNER")
	apiUrl, _ := hf.GetEnv("API_URL")

	commentApiUrl := apiUrl + "/repos/" + repoOwner + "/" + repoName + "/issues/" + getIssueNumberAsString(issue) + "/comments"

	hf.Log("🌍 " + commentApiUrl)

	headers := map[string]string{
		"Accept":        "application/vnd.github+json",
		"Authorization": "Bearer " + botToken,
	}

	jsondoc := `{"body": "` + message + `"}`

	hf.Log("📝 " + jsondoc)

	ret, err := hf.Http(commentApiUrl, "POST", headers, jsondoc)
	if err != nil {
		hf.Log("😡 error:" + err.Error())
	} else {
		hf.Log("👋 Return value from GitHub: " + ret)
	}

}

func Handle(request hf.Request) (response hf.Response, errResp error) {

	/*
		hf.Log("📝 URI: " + request.Uri)
		hf.Log("📝 Method: " + request.Method)
		hf.Log("📝 Body: " + request.Body)
	*/

	action := gjson.Get(request.Body, "action")

	issue := gjson.Get(request.Body, "issue")

	// Create an issue
	if action.Str == Opened {
		hf.Log("📝 Issue Creation:")
		displayIssue(issue)
	}

	// Add comment to the issue
	if action.Str == Created {
		hf.Log("📝 Comment added to the issue:")
		displayIssueTitle(issue)
		comment := gjson.Get(request.Body, "comment")
		displayComment(issue, comment)
	}

	if action.Str == Edited {
		comment := gjson.Get(request.Body, "comment")
		if comment.Exists() { // update of a comment
			hf.Log("📝 Comment updated:")
			displayIssueTitle(issue)
			displayComment(issue, comment)
		} else { // update of an issue
			hf.Log("📝 Issue updated:")
			displayIssue(issue)
		}
	}

	headersResp := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	jsondoc := `{"message": "😀"}`
	//jsondoc, _ = sjson.Set(jsondoc, "message", "👋 hello " + name.Str)

	return hf.Response{Body: jsondoc, Headers: headersResp}, nil
}

/*
curl -v -X POST   http://localhost:8080   -H 'content-type: application/json'   -d '{"name": "Bob"}'
*/

/*
https://docs.github.com/en/rest/issues/comments#create-an-issue-comment

curl \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer <YOUR-TOKEN>" \
  https://api.github.com/repos/OWNER/REPO/issues/ISSUE_NUMBER/comments \
  -d '{"body":"Me too"}'
*/
