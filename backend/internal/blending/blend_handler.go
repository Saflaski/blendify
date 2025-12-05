package blend

import (
	"backend-lastfm/internal/utility"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

type BlendHandler struct {
	frontendUrl         string
	sessionIdCookieName string
	svc                 BlendService
}

type UUID string

func NewBlendHandler(frontendUrl, sidName string, service BlendService) *BlendHandler {
	return &BlendHandler{
		frontendUrl:         frontendUrl,
		sessionIdCookieName: sidName,
		svc:                 service}
}

type BlendRequest struct {
	category     string
	timeDuration string
	user         string
}

func (h *BlendHandler) GenerateNewLink(w http.ResponseWriter, r *http.Request) {
	//Extract cookie?
	glog.Info("Entered GenerateNewLink")
	cookie, err := r.Cookie(h.sessionIdCookieName)
	if err != nil {
		//Something must have gone wrong during runtime for this error to happen
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error during post-validation cookie extraction. Contact Admin")
		glog.Error("Error during post-validation cookie extraction, %w", err)
	}

	userA := userid(cookie.Value)

	newURL, err := h.svc.GenerateNewLinkAndAssignToUser(r.Context(), userA)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error during generating new link. Contact Admin")
		glog.Error("Error during generating new link, %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"linkId": string(newURL)}
	json.NewEncoder(w).Encode(response)

}

func (h *BlendHandler) GetBlendPercentage(w http.ResponseWriter, r *http.Request) { //TODO: Change the name to soemthing else including children functions in S and R
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

	userA := UUID(cookie.Value)
	userB := blendReq.user
	category := blendCategory(blendReq.category) //artist
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

type responseStruct struct {
	Value string `json:"value"`
}

// This is where frontend consumes an invite link and expects a blend id in response.
// An auth association is made with user and blend id
func (h *BlendHandler) AddBlendFromInviteLink(w http.ResponseWriter, r *http.Request) {

	glog.Info("Entered AddBlendFromInviteLink")

	blendResponse, err := utility.DecodeRequest[responseStruct](r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not decode Invite Link for new blend")
	}

	blendLinkValue := blendLinkValue(blendResponse.Value)
	if blendResponse.Value == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Did not add blendlink value")
	}

	cookie, err := r.Cookie(h.sessionIdCookieName)
	if err != nil {
		//Something must have gone wrong during runtime for this error to happen
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error during post-validation cookie extraction. Contact Admin")
		glog.Error("Error during post-validation cookie extraction, %w", err)
	}
	userA := userid(cookie.Value)

	//Validate link?

	blendId, err := h.svc.AddOrMakeBlendFromLink(r.Context(), userA, blendLinkValue)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, " Could not add/make blend from link: %s", err)
		glog.Errorf(" Could not add make/blend from link : %s from user :%s", blendLinkValue, userA)
	}

	if blendId == "0" { //Code for user trying to make blend with themselves
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, " Cannot make blend with yourself: %s", err)
		glog.Errorf(" User tried to make blend with themselves : %s :%s", blendLinkValue, userA)

	}

	// _ = blendId
	glog.Infof("Blend Link Value: %s, User: %s", blendLinkValue, userA)

	w.WriteHeader(http.StatusOK)
	resp := map[string]string{
		"blendId": string(blendId),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
