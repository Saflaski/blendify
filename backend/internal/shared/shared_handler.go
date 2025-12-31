package shared

import (
	"net/http"

	"github.com/golang/glog"
)

type SharedHandler struct {
	svc SharedService
}

func NewSharedHandler(s SharedService) *SharedHandler {
	return &SharedHandler{svc: s}
}

func (*SharedHandler) DeleteAllData(w http.ResponseWriter, r *http.Request) {
	glog.Info("Entered Delete All Data")
}
