package blend

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

type BlendHandler struct {
	frontendUrl         string
	sessionIdCookieName string
	svc                 BlendService
}

func NewBlendHandler(service BlendService) *BlendHandler {
	return &BlendHandler{svc: service}
}

type BlendRequest struct {
	category     string
	timeDuration string
	user         string
}

func (h *BlendHandler) GetNewBlend(w http.ResponseWriter, r *http.Request) {
	response := r.URL.Query()

	blendReq := BlendRequest{
		category:     response.Get("category"),
		timeDuration: response.Get("timeDuration"),
		user:         response.Get("user"),
	}
	defer r.Body.Close()

	if blendReq.category == "" || blendReq.timeDuration == "" || blendReq.user == "" {
		http.Error(w, "Missing required query parameters", http.StatusBadRequest)
		return
	}

	//Need to extract userA aka user client who is sending the blend request
	//As cookie has already been validated during auth, we don't need to cookie check
	cookie, err := r.Cookie(h.sessionIdCookieName)
	if err != nil {
		//Something must have gone wrong during runtime for this error to happen
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error during post-validation cookie extraction. Contact Admin")
		glog.Error("Error during post-validation cookie extraction, %w", err)
	}

	userA := cookie.Value
	userB := blendReq.user
	category := blendCategory(blendReq.category)
	timeDuration := blendTimeDuration(blendReq.timeDuration)

	blendNumber, err := h.svc.GetBlend(userA, userB, category, timeDuration)
	if err != nil {
		http.Error(w, "Error calculating blend", http.StatusInternalServerError)
		return
	}
	responseString := fmt.Sprintf(`{"blend_percentage": %d}`, blendNumber)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, responseString)

}
