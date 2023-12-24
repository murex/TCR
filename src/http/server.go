/*
Copyright (c) 2023 Murex

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

package http

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/murex/tcr/engine"
	"github.com/murex/tcr/http/api"
	"github.com/murex/tcr/role"
	"github.com/murex/tcr/runmode"
	"github.com/murex/tcr/ui"
	"net/http"
	"time"
)

// Server is the user interface implementation when using the HTTP server
type Server struct {
	//reportingChannel chan bool
	tcr  engine.TCRInterface
	port int
	//params           params.Params
}

// Start starts TCR HTTP server
func (s Server) Start() {
	startHttpServer(s.tcr, s.port)
}

func (s Server) ShowRunningMode(mode runmode.RunMode) {
	//TODO implement me
	panic("implement me")
}

func (s Server) NotifyRoleStarting(r role.Role) {
	//TODO implement me
	panic("implement me")
}

func (s Server) NotifyRoleEnding(r role.Role) {
	//TODO implement me
	panic("implement me")
}

func (s Server) ShowSessionInfo() {
	//TODO implement me
	panic("implement me")
}

func (s Server) Confirm(message string, def bool) bool {
	//TODO implement me
	panic("implement me")
}

func (s Server) StartReporting() {
	//TODO implement me
	panic("implement me")
}

func (s Server) StopReporting() {
	//TODO implement me
	panic("implement me")
}

func (s Server) MuteDesktopNotifications(muted bool) {
	//TODO implement me
	panic("implement me")
}

// New creates a new instance of terminal
func New(port int, tcr engine.TCRInterface) ui.UserInterface {
	return Server{
		//reportingChannel: nil,
		tcr:  tcr,
		port: port,
		//params:           params.Params{},
	}
}

// startHttpServer runs TCR HTTP server, listening on the provided port number
// and attaching it to the provided TCR instance
func startHttpServer(tcr engine.TCRInterface, port int) {
	api.SetTCRInstance(tcr)

	router := gin.Default()

	if gin.Mode() != gin.ReleaseMode {
		// Add CORS Middleware in development mode to allow running
		// backend and frontend on separate ports
		router.Use(corsMiddleware())
	}

	// TODO: improvements: 2 modes - development vs production
	// When in production:
	// - using env:   export GIN_MODE=release
	// - or using code:  gin.SetMode(gin.ReleaseMode)
	// - provide a way to configure HTTP port number externally
	// - turn off CORS middleware (should not be necessary)
	// When in development:
	// - use the real static FileSystem instead of the embedded one

	// Serve frontend static files from embedded filesystem
	router.Use(static.Serve("/", embedFolder(staticFS, "static/webapp/browser")))
	router.NoRoute(func(c *gin.Context) {
		// TODO replace fmt.Printf
		fmt.Printf("%s doesn't exists, redirecting to /\n", c.Request.URL.Path)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	// Setup route group for the API
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/build-info", api.GetBuildInfo)
		apiRoutes.GET("/session-info", api.GetSessionInfo)
	}

	// Start HTTP server
	go func() {
		// TODO handle error
		_ = router.Run(fmt.Sprintf("127.0.0.1:%d", port))
	}()

	// TODO - deal with opening of webapp page in a browser
	// Open application page in browser
	// err := browser.OpenURL("http://" + addr)
	// if err != nil {
	//	 fmt.Printf("Failed to open browser: %v\n", err.Error())
	// }
	//
	// for {
	//	 time.Sleep(1 * time.Second)
	// }
}

func corsMiddleware() gin.HandlerFunc {
	// TODO replace fmt.Printf
	fmt.Printf("- Plugging CORS middleware\n")
	return cors.New(cors.Config{
		AllowWildcard:    true,
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Origin"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
		AllowWebSockets:  false,
		MaxAge:           12 * time.Hour,
	})
}
