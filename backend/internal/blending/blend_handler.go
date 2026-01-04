package blend

import (
	"backend-lastfm/internal/auth"
	"backend-lastfm/internal/utility"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

type BlendHandler struct {
	frontendUrl         string
	sessionIdCookieName string
	userKey             contextKey
	svc                 BlendService
}

type contextKey string

type UUID string

func NewBlendHandler(frontendUrl, sidName string, service BlendService, userKey string) *BlendHandler {
	return &BlendHandler{
		frontendUrl:         frontendUrl,
		sessionIdCookieName: sidName,
		svc:                 service,
		userKey:             contextKey(userKey),
	}
}

type BlendRequest struct {
	category     string
	timeDuration string
	user         string
}

func (h *BlendHandler) GetUserBlends(w http.ResponseWriter, r *http.Request) {
	userA, err := h.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, " could not validate session id during generating new link. Contact Admin")
		glog.Error("Error during generating new link, %w", err)
		return
	}

	blends, err := h.svc.GetUserBlends(r.Context(), userA)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, " could not get blends. Contact Admin")
		glog.Error("Error during generating blends for user:%s, %w", userA, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(blends)

}

func (h *BlendHandler) GenerateNewLink(w http.ResponseWriter, r *http.Request) {
	//Extract cookie?
	glog.Info("Entered GenerateNewLink")
	// cookie, err := r.Cookie(h.sessionIdCookieName)
	// if err != nil {
	// 	//Something must have gone wrong during runtime for this error to happen
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintf(w, "Error during post-validation cookie extraction. Contact Admin")
	// 	glog.Error("Error during post-validation cookie extraction, %w", err)
	// }

	userA, err := h.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, " could not validate session id during generating new link. Contact Admin")
		glog.Error("Error during generating new link, %w", err)
	}

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

func (h *BlendHandler) GetUserIdFromContext(ctx context.Context) (userid, error) {
	// t := r.Context().Value(auth.UserKey).(string)

	//As we use auth.UserKey of type contextKey which is owned by auth package, we have to import it
	// Apparently this is idiomatic so we use it. "The package that adds to context owns it"
	t, err := auth.GetUserIDFromContext(ctx)
	return userid(t), err
}

func (h *BlendHandler) GetBlendHealth(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetUserIdFromContext(r.Context())
	if err != nil {
		glog.Errorf("Could not find id during blend health check: %s, %w", id, err)
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"id": string(id)}
	json.NewEncoder(w).Encode(response)
}

func (h *BlendHandler) GetBlendedEntryData(w http.ResponseWriter, r *http.Request) {
	// /dataentry?blendId=&category=&timeDuration=&type=
	response := r.URL.Query()

	blendId := blendId(response.Get("blendId"))
	category := blendCategory(response.Get("category"))
	timeDuration := blendTimeDuration(response.Get("duration"))

	id, err := h.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not interpret userid from request. Either try deleting all cookies and trying again or contact admin")
		glog.Errorf("Could not parse userid from context, %s", id)
		return
	}

	blendedData, err := h.svc.GetBlendEntryByBlendId(r.Context(), blendId, category, timeDuration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not get blended entry data ")
		glog.Errorf(" could not get blended entry data %s for user %s with error: %w", blendId, id, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blendedData)
}

func (h *BlendHandler) GetBlendPageData(w http.ResponseWriter, r *http.Request) {

	response := r.URL.Query()

	blendId := blendId(response.Get("blendId"))
	id, err := h.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not interpret userid from request. Either try deleting all cookies and trying again or contact admin")
		glog.Errorf("Could not parse userid from context, %s", id)
		return
	}

	ok, err := h.svc.AuthoriseBlend(r.Context(), blendId, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not interpret find userid. Either try deleting all cookies and trying again or contact admin")
		glog.Errorf("Could not find userid in repo, %s -> %s", id, blendId)
		return
	}
	defer r.Body.Close()

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, " User does not exist in this blend. Please try accepting invite again.")
		glog.Infof("User unauth access %s -> %s", id, blendId)
		return
	}

	blendData, err := h.svc.GetDuoBlendData(r.Context(), blendId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not get blend data ")
		glog.Errorf(" could not get blend data %s for user %s with error: %w", blendId, id, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blendData)

}

// func (h *BlendHandler) GetBlendPercentage(w http.ResponseWriter, r *http.Request) { //TODO: Change the name to soemthing else including children functions in S and R
// 	response := r.URL.Query()

// 	blendReq := BlendRequest{
// 		category:     response.Get("category"),
// 		timeDuration: response.Get("timeDuration"),
// 		user:         response.Get("user"),
// 	}
// 	defer r.Body.Close()

// 	if blendReq.category == "" || blendReq.timeDuration == "" || blendReq.user == "" {
// 		http.Error(w, "Missing required query parameters", http.StatusBadRequest)
// 		return
// 	}

// 	//Need to extract userA aka user client who is sending the blend request
// 	//As cookie has already been validated during auth, we don't need to cookie check
// 	cookie, err := r.Cookie(h.sessionIdCookieName)
// 	if err != nil {
// 		//Something must have gone wrong during runtime for this error to happen
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprintf(w, "Error during post-validation cookie extraction. Contact Admin")
// 		glog.Error("Error during post-validation cookie extraction, %w", err)
// 	}

// 	userA := UUID(cookie.Value)
// 	userB := blendReq.user
// 	category := blendCategory(blendReq.category) //artist
// 	timeDuration := blendTimeDuration(blendReq.timeDuration)

// 	blendNumber, err := h.svc.GetBlend(r.Context(), userA, userB, category, timeDuration)
// 	if err != nil {
// 		http.Error(w, "Error calculating blend", http.StatusInternalServerError)
// 		return
// 	}
// 	responseString := fmt.Sprintf(`{"blend_percentage": %d}`, blendNumber)
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Fprint(w, responseString)

// }

// This is where frontend consumes an invite link and expects a blend id in response.
// An auth association is made with user and blend id
func (h *BlendHandler) AddBlendFromInviteLink(w http.ResponseWriter, r *http.Request) {
	// /add
	glog.Info("Entered AddBlendFromInviteLink")

	blendResponse, err := utility.DecodeRequest[responseStruct](r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not decode Invite Link for new blend")
		return
	}

	blendLinkValue := blendLinkValue(blendResponse.Value)
	if blendResponse.Value == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Did not add blendlink value")
		return
	}

	// cookie, err := r.Cookie(h.sessionIdCookieName)
	// if err != nil {
	// 	//Something must have gone wrong during runtime for this error to happen
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintf(w, "Error during post-validation cookie extraction. Contact Admin")
	// 	glog.Error("Error during post-validation cookie extraction, %w", err)
	// }
	userA, err := h.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Could not interpret userid from request. Either try deleting all cookies and trying again or contact admin")
		glog.Errorf("Could not parse userid from context")
		return
	}

	//Validate link?

	blendId, err := h.svc.AddOrMakeBlendFromLink(r.Context(), userA, blendLinkValue)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, " Could not add/make blend from link: %s", err)
		glog.Errorf(" Could not add make/blend from link : %s from user :%s and error: %s", blendLinkValue, userA, err)
		return
	}

	if blendId == "0" { //Code for user trying to make blend with themselves
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, " Cannot make blend with yourself: %s", err)
		glog.Errorf(" User tried to make blend with themselves : %s :%s", blendLinkValue, userA)
		return
	} else if blendId == "-1" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, " Cannot add more than 2 users to blend: %s", err)
		glog.Errorf(" Not enough space on blend : %s :%s", blendLinkValue, userA)
		return
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

func (h *BlendHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	glog.Info("Entered GetUserInfo")

	userid, err := h.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, " Could not find user")
		glog.Errorf(" Could not find userinfo due to cannot find userid: %s and error: %s", userid, err)
		return
	}

	userInfo, err := h.svc.GetUserInfo(r.Context(), userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, " Could not get user info")
		glog.Errorf(" Could not get user info for userid: %s and error: %s", userid, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)

}

func (h *BlendHandler) DeleteBlend(w http.ResponseWriter, r *http.Request) {

	glog.Info("Entered DeleteBlend")

	// response := r.URL.Query()
	// blendId := blendId(response.Get("blendId"))
	userid, err := h.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, " Could not find user")
		glog.Errorf(" Could not delete blend due to cannot find user: %s and error: %s", userid, err)
		return
	}
	blendResponse, err := utility.DecodeRequest[deleteStruct](r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not decode Blend Id for deleting blend")
		glog.Errorf("Could not decode Blend Id for deleting blend from user", userid)
		return
	}

	blendId := blendId(blendResponse.BlendId)

	err = h.svc.DeleteBlend(r.Context(), userid, blendId)
	if err != nil || blendId == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, " Could not delete blend: %s", blendId)
		glog.Errorf(" Could not delete blend: %s and error: %s", blendId, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	resp := map[string]string{
		"blendId": string(blendId),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
