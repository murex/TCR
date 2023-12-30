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
	"github.com/murex/tcr/utils"
	"net/http"
	"time"
)

// Server provides a TCR interface implementation over HTTP. It acts
// as a proxy between the TCR engine and HTTP clients
type Server struct {
	tcr        engine.TCRInterface
	port       int
	host       string
	devMode    bool
	websockets []*websocketMessageReporter
	// params           params.Params
}

// New creates a new instance of Server
func New(port int, tcr engine.TCRInterface) *Server {
	server := Server{
		tcr:  tcr,
		host: "0.0.0.0", // To enable connections from a remote host
		// host: "127.0.0.1", // To restrict connections to local host only
		port:       port,
		devMode:    true,
		websockets: []*websocketMessageReporter{},
		// params:           params.Params{},
	}
	tcr.AttachUI(&server, false)
	ServerInstance = &server
	return &server
}

var (
	// ServerInstance is HTTP Server singleton instance
	ServerInstance *Server
)

func (s *Server) registerWebSocket(ws *websocketMessageReporter) {
	s.websockets = append(s.websockets, ws)
	// report.PostInfo("websockets: ", len(s.websockets))
}

func (s *Server) unregisterWebSocket(ws *websocketMessageReporter) {
	for i, registered := range s.websockets {
		if ws == registered {
			s.websockets = append(s.websockets[:i], s.websockets[i+1:]...)
			return
		}
	}
}

// Start starts TCR HTTP server
func (s *Server) Start() {
	utils.Trace("Starting HTTP server on port ", s.port)
	// gin.Default() uses gin.Logger() which should be turned off in TCR production version
	router := gin.New()
	router.Use(gin.Recovery())

	gin.SetMode(gin.ReleaseMode)
	if s.devMode {
		gin.SetMode(gin.DebugMode)
		// In development mode we want to see incoming HTTP requests
		router.Use(gin.Logger())
		// Add CORS Middleware in development mode to allow running
		// backend and frontend on separate ports
		router.Use(corsMiddleware())
	}

	// Serve frontend static files from embedded filesystem
	router.Use(static.Serve("/", embedFolder(staticFS, "static/webapp/browser")))
	router.NoRoute(func(c *gin.Context) {
		utils.Trace(c.Request.URL.Path, " doesn't exists, redirecting to /")
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	// Setup route group for the API
	api.SetTCRInstance(s.tcr)
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/build-info", api.BuildInfoGetHandler)
		apiRoutes.GET("/session-info", api.SessionInfoGetHandler)
	}

	// Setup websocket route
	router.GET("/ws", webSocketHandler)

	// Start HTTP server
	go func() {
		// TODO handle error
		_ = router.Run(s.getServerAddress())
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

// getServerAddress returns the TCP server address that the server is listening to.
func (s *Server) getServerAddress() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

// ShowRunningMode shows the current running mode
func (s *Server) ShowRunningMode(mode runmode.RunMode) {
	for _, ws := range s.websockets {
		ws.ReportTitle(false, "Running in ", mode.Name(), " mode")
	}
}

// NotifyRoleStarting tells the user that TCR engine is starting with the provided role
func (s *Server) NotifyRoleStarting(r role.Role) {
	for _, ws := range s.websockets {
		ws.ReportTitle(false, "Starting with ", r.LongName())
	}
}

// NotifyRoleEnding tells the user that TCR engine is ending the provided role
func (s *Server) NotifyRoleEnding(r role.Role) {
	for _, ws := range s.websockets {
		ws.ReportInfo(false, "Ending ", r.LongName())
	}
}

// ShowSessionInfo shows main information related to the current TCR session
func (*Server) ShowSessionInfo() {
	// With HTTP server this operation is triggered by the client through
	// a GET request. There is nothing to do here
}

// Confirm asks the user for confirmation
func (*Server) Confirm(_ string, _ bool) bool {
	// Always return true until there is a need for this function
	return true
}

// StartReporting tells HTTP server to start reporting information
func (*Server) StartReporting() {
	// Not needed: subscription is managed by each websocket handler instance
}

// StopReporting tells HTTP server to stop reporting information
func (*Server) StopReporting() {
	// Not needed: subscription is managed by each websocket handler instance
}

// MuteDesktopNotifications allows preventing desktop Notification popups from being displayed.
// Used for test automation at the moment. Could be turned into a feature later if there is need for it.
func (*Server) MuteDesktopNotifications(_ bool) {
	// With HTTP server this operation should be triggered by the client though
	// a GET request. There is nothing to do here
}

// corsMiddleware opens CORS connections to HTTP server instance.
// So far this is required only during development (mainly during frontend development
// where frontend application is running on its own HTTP server instance)
// Depending on future evolutions there could be a need to open CORS in production
// too (may require finer tuning in this case to limit CORS to what is needed only)
func corsMiddleware() gin.HandlerFunc {
	utils.Trace("Enabling CORS middleware")
	return cors.New(cors.Config{
		AllowWildcard:    true,
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Origin"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
		AllowWebSockets:  true,
		MaxAge:           12 * time.Hour,
	})
}
