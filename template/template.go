package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IncomingTemplateReqs struct {
	Owner              string `json:"owner"`
	Name               string `json:"name"`
	IncludeAllBranches bool   `json:"include_all_branches"`
	Description        string `json:"description"`
	GithubAccessToken  string `json:"githubAccessToken"`
}

type GithubNecessaryReqs struct {
	Owner              string `json:"owner"`
	Name               string `json:"name"`
	IncludeAllBranches bool   `json:"include_all_branches"`
	Description        string `json:"description"`
}

// Create new repo from template
func CreateNewRepo(c *gin.Context) {

	// Get repo template
	// getRepoTemplate(c)
	var temp IncomingTemplateReqs
	// Json req body
	c.BindJSON(&temp)

	// accessible slice to set access token in header
	t := IncomingTemplateReqs{
		temp.Owner,
		temp.Name,
		temp.IncludeAllBranches,
		temp.Description,
		temp.GithubAccessToken,
	}

	// Slice to be sent to github, access token is not part of the request body
	githubReqs := GithubNecessaryReqs{
		temp.Owner,
		temp.Name,
		temp.IncludeAllBranches,
		temp.Description,
	}

	b, err := json.Marshal(githubReqs)

	if err != nil {
		panic(err)
	}

	body := []byte(b)

	req, reqErr := http.NewRequest("POST", "https://api.github.com/repos/AlexShiresRoth/go-template-test/generate", bytes.NewBuffer(body))

	authToken := fmt.Sprintf("Bearer %v", t.GithubAccessToken)

	if reqErr != nil {

		fmt.Println("Error in initial request: ", err)

		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	req.Header.Add("Accept", "application/vnd.github+json")

	// We need to dynamically get the token from the user
	req.Header.Add("Authorization", authToken)

	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}

	resp, reqError := client.Do(req)

	if reqError != nil {

		fmt.Println("Error in request: ", reqError)

		log.Fatal(reqError)
	}

	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)

	if error != nil {

		fmt.Println("Error in reading body: ", error)

		log.Fatal(error)
	}

	c.Data(http.StatusOK, "application/json", body)

}
