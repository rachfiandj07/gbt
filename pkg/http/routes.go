package http

import (
	"github.com/golang-base-template/util/middleware"
	"github.com/julienschmidt/httprouter"
)

func AssignRoutes(router *httprouter.Router) {
	router.GET("/get-data/:source", middleware.ChainReq(GetData, middleware.InitContext, middleware.SetHeader))
}
