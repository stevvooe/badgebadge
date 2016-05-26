package main

import (
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
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

func main() {
	iris.Get("/github.com/:username/:reponame", func(ctx *iris.Context) {
		username := ctx.Param("username")
		reponame := ctx.Param("reponame")
		queryString := ctx.URI().QueryArgs()
		logrus.Debug("branch: %v", string(queryString.Peek("branch")))
		badge, err := checkBadges(username, reponame, string(queryString.Peek("branch")))
		if err != nil {
			ctx.Write(err.Error())
			ctx.SetStatusCode(iris.StatusInternalServerError)
		}
		ctx.ServeFile(badge, false)
	})

	logrus.Info("Server listening on :8080")
	iris.Listen(":8080")

}
