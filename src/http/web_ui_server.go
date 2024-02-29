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
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/murex/tcr/engine"
	"github.com/murex/tcr/http/api"
	"github.com/murex/tcr/http/ws"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/runmode"
	"net/http"
	"time"
)

// WebUIServer provides a TCR interface implementation over HTTP. It acts
// as a proxy between the TCR engine and HTTP clients
type WebUIServer struct {
	tcr              engine.TCRInterface
	params           params.Params
	host             string
	devMode          bool
	router           *gin.Engine
	httpServer       *http.Server
	websocketTimeout time.Duration
}

// New creates a new instance of WebUIServer
func New(p params.Params, tcr engine.TCRInterface) *WebUIServer {
	webUIServer := WebUIServer{
		tcr: tcr,
		// host: "0.0.0.0", // To enable connections from a remote host
		host:             "127.0.0.1", // To restrict connections to local host only
		devMode:          true,
		router:           nil,
		httpServer:       nil,
		websocketTimeout: 1 * time.Minute, // default timeout value
		params:           p,
	}
	tcr.AttachUI(&webUIServer, false)
	return &webUIServer
}

// Start starts TCR HTTP server
func (webUIServer *WebUIServer) Start() {
	report.PostInfo("Starting HTTP server on port ", webUIServer.params.PortNumber)
	webUIServer.initGinEngine()
	webUIServer.addStaticRoutes()
	webUIServer.addAPIRoutes()
	webUIServer.addWebsocketRoutes()
	webUIServer.startGinEngine()
}

func (webUIServer *WebUIServer) initGinEngine() {
	// gin.Default() uses gin.Logger() which should be turned off in TCR production version
	webUIServer.router = gin.New()
	webUIServer.router.Use(gin.Recovery())

	gin.SetMode(gin.ReleaseMode)
	if webUIServer.InDevMode() {
		gin.SetMode(gin.DebugMode)
		// In development mode we want to see incoming HTTP requests
		webUIServer.router.Use(gin.Logger())
		// Add CORS Middleware in development mode to allow running
		// backend and frontend on separate ports
		webUIServer.router.Use(corsMiddleware())
	}
}

func (webUIServer *WebUIServer) startGinEngine() {
	// Create HTTP server instance
	webUIServer.httpServer = &http.Server{ //nolint:gosec
		Addr:    webUIServer.GetServerAddress(),
		Handler: webUIServer.router,
	}

	// Start HTTP server on its own goroutine
	go func() {
		err := webUIServer.httpServer.ListenAndServe()
		if err != nil {
			report.PostError("could not start HTTP server: ", err.Error())
		}
	}()
}

// stopGinEngine is provided for testing purpose, so that we can shutdown
// the HTTP server when needed
func (webUIServer *WebUIServer) stopGinEngine() {
	report.PostInfo("Stopping HTTP server")
	if err := webUIServer.httpServer.Shutdown(context.Background()); err != nil {
		report.PostError("could not stop HTTP server: ", err.Error())
	}
}

func (webUIServer *WebUIServer) addStaticRoutes() {
	// Serve frontend static files from embedded filesystem
	webUIServer.router.Use(static.Serve("/", embedFolder(staticFS, "static/webapp/browser")))
	webUIServer.router.NoRoute(func(c *gin.Context) {
		report.PostInfo(c.Request.URL.Path, " doesn't exists, redirecting to /")
		c.Redirect(http.StatusMovedPermanently, "/")
	})
}

func (webUIServer *WebUIServer) addAPIRoutes() {
	// Add TCR engine to gin context so that it can be accessed by API handlers
	webUIServer.router.Use(api.TCREngineMiddleware(webUIServer.tcr))
	// Setup route group for the API
	apiRoutes := webUIServer.router.Group("/api")
	{
		apiRoutes.GET("/build-info", api.BuildInfoGetHandler)
		apiRoutes.GET("/session-info", api.SessionInfoGetHandler)
		apiRoutes.GET("/roles", api.RolesGetHandler)
		apiRoutes.GET("/roles/:name", api.RoleGetHandler)
		apiRoutes.POST("/roles/:name/:action", api.RolesPostHandler)
	}
}

func (webUIServer *WebUIServer) addWebsocketRoutes() {
	// Add self to gin context so that it can be accessed by web socket handlers
	webUIServer.router.Use(ws.HTTPServerMiddleware(webUIServer))
	// Setup websocket route
	webUIServer.router.GET("/ws", ws.WebsocketHandler)
}

// InDevMode indicates if the server is running in dev (development) mode
func (webUIServer *WebUIServer) InDevMode() bool {
	return webUIServer.devMode
}

// GetServerAddress returns the TCP server address that the server is listening to.
func (webUIServer *WebUIServer) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", webUIServer.host, webUIServer.params.PortNumber)
}

// GetWebsocketTimeout returns the timeout after which inactive websocket connections
// should be closed
func (webUIServer *WebUIServer) GetWebsocketTimeout() time.Duration {
	return webUIServer.websocketTimeout
}

// ShowRunningMode shows the current running mode
func (*WebUIServer) ShowRunningMode(_ runmode.RunMode) {
	// Not needed: Runmode query will be handled by the client through a GET request
}

// ShowSessionInfo shows main information related to the current TCR session
func (*WebUIServer) ShowSessionInfo() {
	// With HTTP server this operation is triggered by the client through
	// a GET request. There is nothing to do here
}

// Confirm asks the user for confirmation
func (*WebUIServer) Confirm(_ string, _ bool) bool {
	// Always return true until there is a need for this function
	return true
}

// StartReporting tells HTTP server to start reporting information
func (*WebUIServer) StartReporting() {
	// Not needed: subscription is managed by each websocket handler instance
}

// StopReporting tells HTTP server to stop reporting information
func (*WebUIServer) StopReporting() {
	// Not needed: subscription is managed by each websocket handler instance
}

// MuteDesktopNotifications allows preventing desktop Notification popups from being displayed.
// Used for test automation at the moment. Could be turned into a feature later if there is need for it.
func (*WebUIServer) MuteDesktopNotifications(_ bool) {
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
