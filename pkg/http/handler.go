package http

import (
	basicctrl "github.com/golang-base-template/internal/basic_module/controller"
	"github.com/golang-base-template/util/response"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var (
	basicCtrl basicctrl.BasicController
)

func Init() {
	if basicCtrl == nil {
		basicCtrl = basicctrl.NewBasicController()
	}
}

func GetData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	res := response.New(r.Header.Get("Origin"), "true")
	source := p.ByName("source")
	result, err := basicCtrl.GetData(ctx)
	result += "." + source
	if err != nil {
		res.WriteError(w, http.StatusInternalServerError, []string{"error when get data"}, err.Error())
		return
	}

	res.WriteResponse(w, result)
	return
}
