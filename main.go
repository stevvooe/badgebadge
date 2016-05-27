package main

import (
	"html/template"
	"os"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func init() {
	logLevel := os.Getenv("DEBUG")
	logrus.Infof("Logs: %v", logLevel)
	if logLevel != "" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	badgeMinimum, _ = strconv.Atoi(os.Getenv("BADGE_COUNT"))
	if badgeMinimum == 0 {
		badgeMinimum = 3
	}
}

type Badge struct {
	BadgesHTML template.HTML
	Badge      string
}

func main() {
	iris.Static("/badges", "./badges/", 1)
	iris.Get("/github.com/:username/:reponame", func(ctx *iris.Context) {
		username := ctx.Param("username")
		reponame := ctx.Param("reponame")
		queryString := ctx.URI().QueryArgs()
		branch := string(queryString.Peek("branch"))
		debug := string(queryString.Peek("debug"))

		logrus.Debug("branch: %v - debug: %v", branch, debug)
		badges, err := checkBadges(username, reponame, branch)
		if err != nil {
			ctx.Write(err.Error())
			ctx.SetStatusCode(iris.StatusInternalServerError)
		}
		if debug == "true" {
			htmlBadges := blackfriday.MarkdownBasic([]byte(strings.Join(badges[1:], "\n")))
			htmlBadges = bluemonday.UGCPolicy().SanitizeBytes(htmlBadges)
			if len(htmlBadges) == 0 {
				htmlBadges = []byte("<p><em>No badges found in README.md</em></p>")
			}

			if err := ctx.Render("index.html", Badge{
				BadgesHTML: template.HTML(string(htmlBadges)),
				Badge:      badges[0],
			}); err != nil {
				logrus.Panic(err)
			}
		} else {
			ctx.ServeFile(badges[0], false)
		}
	})

	logrus.Info("Server listening on :8080")
	iris.Listen(":8080")

}
