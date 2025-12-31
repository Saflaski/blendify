package shared

import (
	"backend-lastfm/internal/auth"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

type SharedHandler struct {
	svc SharedService
}

func NewSharedHandler(s SharedService) *SharedHandler {
	return &SharedHandler{svc: s}
}

func (h *SharedHandler) DeleteAllData(w http.ResponseWriter, r *http.Request) {
	glog.Info("Entered Delete All Data")
	userid, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not interpret userid from request. Either try deleting all cookies and trying again or contact admin")
		glog.Errorf("Could not parse userid from context, %s", err)
		return
	}

	err = h.svc.DeleteAllUserData(r.Context(), userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not delete user information. Try again or contact admin")
		glog.Errorf("Could not delete all user data:, %s", err)
		return
	}

	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted all user data")
	glog.Infof("Deleted all user data:, %s", err)

}
