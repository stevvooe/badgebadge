package main

import (
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/view"
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

// Badge holds the html content for the report
type Badge struct {
	HTML    template.HTML
	Badge   string
	Title   string
	Verdict string
	Color   string
	Code    string
}

func main() {
	app := iris.New()
	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())

	tmpl := view.HTML("./templates", ".html")
	tmpl.Layout("layout.html")

	app.Adapt(tmpl)

	app.StaticWeb("/badges", "./badges/")
	app.StaticWeb("/assets", "./assets/")
	app.Get("/github.com/:username/:reponame", handleReport)
	app.Get("/report/github.com/:username/:reponame", handleReport)
	app.Get("/report", showReport)
	app.Get("/", handleHome)
	// logrus.Debug("Server listening on :8080")
	app.Listen(":8080")
}

func handleHome(ctx *iris.Context) {
	ctx.MustRender("index.html", nil)
}

func showReport(ctx *iris.Context) {
	repo := ctx.URLParam("repo")
	ctx.Redirect("/report/" + repo)
}

func handleReport(ctx *iris.Context) {
	username := ctx.Param("username")
	reponame := ctx.Param("reponame")
	branch := ctx.URLParam("branch")
	debug := ctx.URLParam("debug")
	uri := ctx.RequestPath(true)
	pathSlice := strings.Split(uri, "/")

	if pathSlice[1] == "report" {
		debug = "true"
	}

	logrus.Debug("branch: %v - debug: %v", branch, debug)

	badges, title, err := checkBadges(username, reponame, branch)
	if err != nil {
		ctx.Render("error.html", struct{ Message string }{Message: err.Error()})
		ctx.SetStatusCode(iris.StatusInternalServerError)
		return
	}
	if debug == "true" {
		htmlBadges := blackfriday.MarkdownBasic([]byte(strings.Join(badges[1:], "\n")))
		htmlBadges = bluemonday.UGCPolicy().SanitizeBytes(htmlBadges)
		if len(htmlBadges) == 0 {
			htmlBadges = []byte("<p><em>No badges found in README.md</em></p>")
		}
		verdict := "no"
		color := "red"
		if title == "adequate" {
			verdict = "yes"
			color = ""
		}
		githubRepo := fmt.Sprintf("github.com/%s/%s", username, reponame)
		code := "[![Badge Badge](http://doyouevenbadge.com/github.com/" + githubRepo + ")](http://doyouevenbadge.com/report/github.com/" + githubRepo + ")"

		if err := ctx.Render("report.html", Badge{
			HTML:    template.HTML(string(htmlBadges)),
			Badge:   badges[0],
			Title:   title,
			Verdict: verdict,
			Color:   color,
			Code:    code,
		}); err != nil {
			logrus.Panic(err)
		}
	} else {
		ctx.ServeFile(badges[0], false)
	}
}
