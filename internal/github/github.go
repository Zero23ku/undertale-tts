package github

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Commit struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

type Release struct {
	Name       string `json:"name"`
	ZipballUrl string `json:"zipball_url"`
	TarballUrl string `json:"tarball_url"`
	Commit     Commit `json:"commit"`
	NodeId     string `json:"node_id"`
}

var releaseURL = "https://api.github.com/repos/Zero23ku/ac-tts-golang/tags"

func GetLatestReleaseVersion() string {
	res, err := http.Get(releaseURL)

	if err != nil {
		//logging.CreateLog("github - couldn't make HTTP request", err)
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		//logging.CreateLog("github - couldn't deserealize response", err)
		log.Fatal(err)
	}

	var response []Release

	if err := json.Unmarshal(body, &response); err != nil {
		//logging.CreateLog("github - couldn't transform response", err)
		return ""
	}

	if len(response) > 0 {
		return response[0].Name
	}
	return ""
}
