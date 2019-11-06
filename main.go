package main

import (
	"html/template"
	"os"
	"strconv"
	"strings"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func init() {
	logLevel := os.Getenv("DEBUG")
	if logLevel != "" {
		golog.Infof("Logs: %v", logLevel)
		// default golog instance is shared with Iris too.
		golog.SetLevel("debug")
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
	app := iris.New()
	app.RegisterView(iris.HTML("./templates", ".html"))
	app.HandleDir("/badges", "./badges/")

	app.Get("/github.com/{username}/{reponame}", handleReport)
	app.Get("/report/github.com/{username}/{reponame}", handleReport)

	app.Run(iris.Addr(":8080"))
}

func handleReport(ctx iris.Context) {
	username := ctx.Params().Get("username")
	reponame := ctx.Params().Get("reponame")
	branch := ctx.URLParam(("branch"))
	debug := ctx.URLParam("debug")
	pathSlice := strings.Split(ctx.Path(), "/")
	if pathSlice[1] == "report" {
		debug = "true"
	}

	golog.Debug("branch: %v - debug: %v", branch, debug)
	badges, err := checkBadges(username, reponame, branch)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
	}
	if debug == "true" {
		htmlBadges := blackfriday.Run([]byte(strings.Join(badges[1:], "\n")))
		htmlBadges = bluemonday.UGCPolicy().SanitizeBytes(htmlBadges)
		if len(htmlBadges) == 0 {
			htmlBadges = []byte("<p><em>No badges found in README.md</em></p>")
		}

		if err := ctx.View("index.html", Badge{
			BadgesHTML: template.HTML(string(htmlBadges)),
			Badge:      badges[0],
		}); err != nil {
			golog.Fatal(err)
		}
	} else {
		ctx.ServeFile(badges[0], false)
	}
}
