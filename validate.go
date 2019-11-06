package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/kataras/golog"
)

var badgeMinimum int

const inadequateBadge = "badges/inadequate.svg"
const adequateBadge = "badges/adequate.svg"

var badgeRex = regexp.MustCompile(`(?i)(\[!\[[a-z0-9_ .]+\]\([0-9a-z.\/:&?=-]+\)\]\([0-9a-z.\/:&?-]+\))`)

func checkBadges(username string, reponame string, branch string) ([]string, error) {
	githubRepo := fmt.Sprintf("github.com/%s/%s", username, reponame)
	if branch == "" {
		branch = "master"
	}
	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/README.md", username, reponame, branch))
	if err != nil {
		// handle error
		golog.Debugf("could not fetch github repo: %s - err: %v", githubRepo, err)
		return []string{}, err
	}
	golog.Debug("resp: ", resp)
	if resp.StatusCode != 200 {
		golog.Debugf("Github repo %s not found", githubRepo)
		return []string{}, errors.New("github repo/branch not found")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		golog.Debugf("could not read body from %s: %v", githubRepo, err)
		return []string{}, err
	}
	golog.Debug("body: ", string(body))

	// The `All` variants of these functions apply to all
	// matches in the input, not just the first. For
	// example to find all matches for a regexp.
	badgeMatch := badgeRex.FindAllString(string(body), -1)
	golog.Debug(badgeMatch)
	golog.Infof("Badge count: %v - min: %v", len(badgeMatch), badgeMinimum)
	if len(badgeMatch) >= badgeMinimum {
		golog.Info("Congrats you haz all badges")
		return append([]string{adequateBadge}, badgeMatch...), nil
	}
	golog.Info("Needz moar badges")
	return append([]string{inadequateBadge}, badgeMatch...), nil
}
