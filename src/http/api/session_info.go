package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type sessionInfo struct {
	BaseDir           string `json:"baseDir"`
	WorkDir           string `json:"workDir"`
	LanguageName      string `json:"language"`
	ToolchainName     string `json:"toolchain"`
	VCSName           string `json:"vcsName"`
	VCSSessionSummary string `json:"vcsSession"`
	CommitOnFail      bool   `json:"commitOnFail"`
	GitAutoPush       bool   `json:"gitAutoPush"`
	MessageSuffix     string `json:"messageSuffix"`
}

func GetSessionInfo(c *gin.Context) {
	info := tcr.GetSessionInfo()
	data := sessionInfo{
		BaseDir:           info.BaseDir,
		WorkDir:           info.WorkDir,
		LanguageName:      info.LanguageName,
		ToolchainName:     info.ToolchainName,
		VCSName:           info.VCSName,
		VCSSessionSummary: info.VCSSessionSummary,
		CommitOnFail:      info.CommitOnFail,
		GitAutoPush:       info.GitAutoPush,
		MessageSuffix:     info.MessageSuffix,
	}
	c.IndentedJSON(http.StatusOK, data)
}
