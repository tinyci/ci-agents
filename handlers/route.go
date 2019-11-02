package handlers

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"golang.org/x/net/websocket"
)

// Routes is a collection of Route structs.
type Routes map[string]map[string]*Route

// Route is a route management struct for generated services.
type Route struct {
	Method             string
	UseCORS            bool
	UseAuth            bool
	WebsocketProcessor WebsocketFunc
	Handler            func(*H, *gin.Context, HandlerFunc) *errors.Error
	Processor          HandlerFunc
	ParamValidator     func(*H, *gin.Context) *errors.Error
	Capability         model.Capability
	TokenScope         string
}

// HandlerFunc is the basic kind of HandlerFunc.
type HandlerFunc func(context.Context, *H, *gin.Context) (interface{}, int, *errors.Error)

// WebsocketFunc is the controller for websocket operations.
type WebsocketFunc func(context.Context, *H, *gin.Context, *websocket.Conn) *errors.Error

// SetProcessor allows you to more simply set the processor for a given route.
func (r Routes) SetProcessor(route string, method string, processor HandlerFunc) {
	r[TransformSwaggerRoute(route)][strings.ToLower(method)].Processor = processor
}

// SetWebsocketProcessor allows you to more simply set the processor for a given route.
func (r Routes) SetWebsocketProcessor(route string, processor WebsocketFunc) {
	r[TransformSwaggerRoute(route)]["get"].WebsocketProcessor = processor
}
