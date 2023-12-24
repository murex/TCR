package api

import (
	"github.com/gin-gonic/gin"
	"github.com/murex/tcr/settings"
	"net/http"
)

type buildInfo struct {
	BuildVersion string `json:"version"`
	BuildOs      string `json:"os"`
	BuildArch    string `json:"arch"`
	BuildCommit  string `json:"commit"`
	BuildDate    string `json:"date"`
	BuildAuthor  string `json:"author"`
}

func GetBuildInfo(c *gin.Context) {
	data := buildInfo{
		BuildVersion: settings.BuildVersion,
		BuildOs:      settings.BuildOs,
		BuildArch:    settings.BuildArch,
		BuildCommit:  settings.BuildCommit,
		BuildDate:    settings.BuildDate,
		BuildAuthor:  settings.BuildAuthor,
	}
	c.IndentedJSON(http.StatusOK, data)
}
