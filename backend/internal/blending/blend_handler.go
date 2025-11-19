package blend

import (
	"fmt"
	"net/http"
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

	//Service to return blend percentage number
	//MOCK DATA
	//mock_data := 42

	category := blendCategory(blendReq.category)
	timeDuration := blendTimeDuration(blendReq.timeDuration)

	blendNumber, err := h.svc.GetBlend(blendReq.user, category, timeDuration)
	if err != nil {
		http.Error(w, "Error calculating blend", http.StatusInternalServerError)
		return
	}
	responseString := fmt.Sprintf(`{"blend_percentage": %d}`, blendNumber)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, responseString)

}
