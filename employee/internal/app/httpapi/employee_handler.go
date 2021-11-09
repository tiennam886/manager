package httpapi

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/tiennam886/manager/employee/internal/service"
	"github.com/tiennam886/manager/pkg/httputil"

	"net/http"
)

func StaffAdd(w http.ResponseWriter, r *http.Request) {
	var payload service.AddEmployeeCommand

	if err := httputil.BindJSON(r, &payload); err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	staff, err := service.AddStaff(r.Context(), payload)
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Added staff uid=%s", staff.UID),
		Data:    staff,
	})
}

func StaffFindByUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	staff, err := service.FindStaffByUID(r.Context(), service.FindStaffByUIDCommand(uid))
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Found staff uid=%s", uid),
		Data:    staff,
	})
}

func StaffDeleteByUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	err := service.DeleteStaffByUID(r.Context(), service.DeleteStaffByUIDCommand(uid))
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Deleted staff uid=%s", uid),
	})
}
