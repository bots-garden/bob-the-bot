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
	Issue         = "issue"
	Note          = "note"
)

func main() {
	hf.SetHandleHttp(Handle)
}

/*
POST /projects/:id/issues/:issue_iid/notes
https://docs.gitlab.com/ee/api/notes.html#create-new-issue-note
curl --request POST --header "PRIVATE-TOKEN: <your_access_token>" "https://gitlab.example.com/api/v4/projects/5/issues/11/notes?body=note"
*/
func addNoteToTheIssue(issueAttributes gjson.Result, projectId float64, message string) (string, error) {

	botToken, _ := hf.GetEnv("BOT_TOKEN")
	//repoName, _ := hf.GetEnv("REPO_NAME")
	//repoOwner, _ := hf.GetEnv("REPO_OWNER")
	apiUrl, _ := hf.GetEnv("API_URL")

	issueIid := issueAttributes.Get("iid").Num
	issueNumber := strconv.FormatInt(int64(issueIid), 10)
	projectNumber := strconv.FormatInt(int64(projectId), 10)

	hf.Log("🖐 issueNumber: " + issueNumber)
	hf.Log("🖐 projectNumber: " + projectNumber)

	noteApiUrl := apiUrl + "/projects/" + projectNumber + "/issues/" + issueNumber + "/notes"

	headers := map[string]string{
		"Content-Type":        "application/json; charset=utf-8",
		"PRIVATE-TOKEN": botToken,
	}

	jsondoc := `{"body": "` + message + `"}`

	return hf.Http(noteApiUrl, "POST", headers, jsondoc)
}

func Handle(request hf.Request) (response hf.Response, errResp error) {
	// project: https://gitlab.com/k33g_org/bob-the-bot
	var resp string
	var err error

	//hf.Log("🤖 " + request.Body)

	/*
			https://docs.gitlab.com/ee/user/project/integrations/webhook_events.html#issue-events

		  "object_kind": "issue",
		  "event_type": "issue",

			object_attributes.action

				open
				close
				reopen
				update
		Comments: https://docs.gitlab.com/ee/user/project/integrations/webhook_events.html#comment-on-an-issue
	*/

	object := gjson.Get(request.Body, "object_kind") // it should be issue or note
	//event := gjson.Get(request.Body, "event_type")   // it should be issue or note

	objectAttributes := gjson.Get(request.Body, "object_attributes") // issue or comment

	projectId := gjson.Get(request.Body, "project.id").Num

	userName := gjson.Get(request.Body, "user.username").Str

	botName, _ := hf.GetEnv("BOT_NAME")

	// an issue is created or updated
	if object.Str == Issue {

		action := objectAttributes.Get("action") // open close reopen update (only for issue)

		if action.Str == Open || action.Str == Update {
			issueDescription := objectAttributes.Get("description").Str
			// we don't need to pass all the attributes, the iid is enough
			if strings.Contains(issueDescription, botName) {
				resp, err = addNoteToTheIssue(objectAttributes, projectId, "👋 @"+userName+" are you talking to me?")
			}
		}
	}

	// a comment is added to the issue or an existing comment is updated
	if object.Str == Note {
		
		note := objectAttributes.Get("note").Str
		issueAttributes := gjson.Get(request.Body, "issue")
		// we don't need to pass all the attributes, the iid is enough
		if strings.Contains(note, botName) {
			resp, err = addNoteToTheIssue(issueAttributes, projectId, "👋 @"+userName+" are you talking to me?")
		}
	}

	headersResp := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	if err != nil {
		hf.Log("😡 error:" + err.Error())
	} 
	/*
	else {
		hf.Log("🦊 GitLab response: " + resp)
	}
	*/

	return hf.Response{Body: resp, Headers: headersResp}, err
}
