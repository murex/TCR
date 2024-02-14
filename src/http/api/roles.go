/*
Copyright (c) 2024 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/role"
	"net/http"
)

type roleData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

const (
	startAction string = "start"
	stopAction  string = "stop"
)

// RolesGetHandler handles HTTP GET requests on all TCR roles
func RolesGetHandler(c *gin.Context) {
	tcr := getTCRInstance(c)
	var data []roleData
	current := tcr.GetCurrentRole()
	for _, r := range role.All() {
		data = append(data, roleData{
			Name:        r.Name(),
			Description: r.LongName(),
			Active:      r == current,
		})
	}
	c.IndentedJSON(http.StatusOK, data)
}

// RoleGetHandler handles HTTP GET requests on a TCR role
func RoleGetHandler(c *gin.Context) {
	tcr := getTCRInstance(c)
	var r role.Role
	name := c.Param("name")
	switch name {
	case role.Navigator{}.Name():
		r = role.Navigator{}
	case role.Driver{}.Name():
		r = role.Driver{}
	default:
		report.PostWarning("unrecognized role: ", name)
		c.Status(http.StatusBadRequest)
		return
	}

	data := roleData{
		Name:        r.Name(),
		Description: r.LongName(),
		Active:      r == tcr.GetCurrentRole(),
	}
	c.IndentedJSON(http.StatusOK, data)
}

// RolesPostHandler handles HTTP POST requests on current TCR role
func RolesPostHandler(c *gin.Context) {
	tcr := getTCRInstance(c)
	name := c.Param("name")
	switch name {
	case role.Navigator{}.Name():
		rolePostHandler(c, role.Navigator{}, tcr.RunAsNavigator, tcr.Stop)
	case role.Driver{}.Name():
		rolePostHandler(c, role.Driver{}, tcr.RunAsDriver, tcr.Stop)
	default:
		report.PostWarning("unrecognized role: ", name)
		c.Status(http.StatusBadRequest)
	}
}

func rolePostHandler(c *gin.Context, r role.Role, starter func(), stopper func()) {
	action := c.Param("action")
	switch action {
	case startAction:
		starter()
		data := roleData{Name: r.Name(), Description: r.LongName(), Active: true}
		c.IndentedJSON(http.StatusAccepted, data)
	case stopAction:
		stopper()
		data := roleData{Name: r.Name(), Description: r.LongName(), Active: false}
		c.IndentedJSON(http.StatusAccepted, data)
	default:
		report.PostWarning("unrecognized action: ", action)
		c.Status(http.StatusBadRequest)
	}
}
