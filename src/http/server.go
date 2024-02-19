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

package http

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/murex/tcr/engine"
	"github.com/murex/tcr/http/api"
	"github.com/murex/tcr/http/ws"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/role"
	"github.com/murex/tcr/runmode"
	"net/http"
	"time"
)

// Server provides a TCR interface implementation over HTTP. It acts
// as a proxy between the TCR engine and HTTP clients
type Server struct {
	tcr              engine.TCRInterface
	params           params.Params
	host             string
	devMode          bool
	websocketTimeout time.Duration
	websockets       *ws.ConnectionPool
}

// New creates a new instance of Server
func New(p params.Params, tcr engine.TCRInterface) *Server {
	server := Server{
		tcr: tcr,
		// host: "0.0.0.0", // To enable connections from a remote host
		host:             "127.0.0.1", // To restrict connections to local host only
		devMode:          true,
		websocketTimeout: 1 * time.Minute, // default timeout value
		websockets:       ws.NewConnectionPool(),
		params:           p,
	}
	tcr.AttachUI(&server, false)
	return &server
}

// Start starts TCR HTTP server
func (s *Server) Start() {
	report.PostInfo("Starting HTTP server on port ", s.params.PortNumber)
	router := s.newGinEngine()

	s.addStaticRoutes(router)
	s.addAPIRoutes(router)
	s.addWebsocketRoutes(router)

	s.startGinEngine(router)
}

func (s *Server) newGinEngine() *gin.Engine {
	// gin.Default() uses gin.Logger() which should be turned off in TCR production version
	router := gin.New()
	router.Use(gin.Recovery())

	gin.SetMode(gin.ReleaseMode)
	if s.InDevMode() {
		gin.SetMode(gin.DebugMode)
		// In development mode we want to see incoming HTTP requests
		router.Use(gin.Logger())
		// Add CORS Middleware in development mode to allow running
		// backend and frontend on separate ports
		router.Use(corsMiddleware())
	}
	return router
}

func (s *Server) startGinEngine(router *gin.Engine) {
	// Start HTTP server on its own goroutine
	go func() {
		err := router.Run(s.GetServerAddress())
		if err != nil {
			report.PostError("could not start HTTP server: ", err.Error())
		}
	}()
}

func (*Server) addStaticRoutes(router *gin.Engine) {
	// Serve frontend static files from embedded filesystem
	router.Use(static.Serve("/", embedFolder(staticFS, "static/webapp/browser")))
	router.NoRoute(func(c *gin.Context) {
		report.PostInfo(c.Request.URL.Path, " doesn't exists, redirecting to /")
		c.Redirect(http.StatusMovedPermanently, "/")
	})
}

func (s *Server) addAPIRoutes(router *gin.Engine) {
	// Add TCR engine to gin context so that it can be accessed by API handlers
	router.Use(api.TCREngineMiddleware(s.tcr))
	// Setup route group for the API
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/build-info", api.BuildInfoGetHandler)
		apiRoutes.GET("/session-info", api.SessionInfoGetHandler)
		apiRoutes.GET("/roles", api.RolesGetHandler)
		apiRoutes.GET("/roles/:name", api.RoleGetHandler)
		apiRoutes.POST("/roles/:name/:action", api.RolesPostHandler)
	}
}

func (s *Server) addWebsocketRoutes(router *gin.Engine) {
	// Add self to gin context so that it can be accessed by web socket handlers
	router.Use(ws.HTTPServerMiddleware(s))
	// Setup websocket route
	router.GET("/ws", ws.WebsocketHandler)
}

// InDevMode indicates if the server is running in dev (development) mode
func (s *Server) InDevMode() bool {
	return s.devMode
}

// GetServerAddress returns the TCP server address that the server is listening to.
func (s *Server) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", s.host, s.params.PortNumber)
}

// GetWebsocketTimeout returns the timeout after which inactive websocket connections
// should be closed
func (s *Server) GetWebsocketTimeout() time.Duration {
	return s.websocketTimeout
}

// RegisterWebsocket register a new websocket connection to the server
func (s *Server) RegisterWebsocket(w ws.WebsocketWriter) {
	s.websockets.Register(w)
}

// UnregisterWebsocket unregister a new websocket connection from the server
func (s *Server) UnregisterWebsocket(w ws.WebsocketWriter) {
	s.websockets.Unregister(w)
}

// ShowRunningMode shows the current running mode
func (s *Server) ShowRunningMode(mode runmode.RunMode) {
	s.websockets.Dispatch(func(w ws.WebsocketWriter) {
		w.ReportTitle(false, "Running in ", mode.Name(), " mode")
	})
}

// NotifyRoleStarting tells the user that TCR engine is starting with the provided role
func (s *Server) NotifyRoleStarting(r role.Role) {
	s.websockets.Dispatch(func(w ws.WebsocketWriter) {
		// ReportRole call is used for role changing trigger message
		w.ReportRole(false, r.Name(), ":", "start")
		// ReportTitle is used for console trace
		w.ReportTitle(false, "Starting with ", r.LongName())
	})
}

// NotifyRoleEnding tells the user that TCR engine is ending the provided role
func (s *Server) NotifyRoleEnding(r role.Role) {
	s.websockets.Dispatch(func(w ws.WebsocketWriter) {
		// ReportRole call is used for role changing trigger message
		w.ReportRole(false, r.Name(), ":", "end")
		// ReportTitle is used for console trace
		w.ReportTitle(false, "Ending ", r.LongName())
	})
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
	report.PostInfo("Using CORS middleware")
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
