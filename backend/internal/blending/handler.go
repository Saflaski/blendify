package blend

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang/glog"
)


type BlendHandler struct {
    frontendUrl string
	sessionIdCookieName string

}

func NewBlendHandler() *BlendHandler {
	return &BlendHandler{}
}

func (h *BlendHandler) GetNewBlend(w http.ResponseWriter, r *http.Request) {
	glog.Infof("New Blend Request with UA: %s, UB: %s", chi.URLParam(r, "UA"), chi.URLParam(r, "UB"))
	glog.Warning("Incomplete Method getNewBlend()")

	
}