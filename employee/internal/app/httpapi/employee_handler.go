package httpapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/tiennam886/manager/employee/internal/service"
	"github.com/tiennam886/manager/pkg/httputil"
	"github.com/tiennam886/manager/pkg/messaging/httpsub"
)

var sugarLogger *zap.SugaredLogger

func EmployeeAdd(w http.ResponseWriter, r *http.Request) {
	var payload service.AddEmployeeCommand

	if err := httputil.BindJSON(r, &payload); err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	employee, err := service.AddEmployee(r.Context(), payload)
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Added staff uid=%s", employee.UID),
		Data:    employee,
	})
}

func EmployeeGetAll(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	employees, err := service.GetAllEmployee(r.Context(), page, limit)
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("All Employees"),
		Data:    employees,
	})
}

func EmployeeFindByUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	employee, err := service.FindStaffByUID(r.Context(), service.FindEmployeeByUIDCommand(uid))
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Found staff uid=%s", uid),
		Data:    employee,
	})
}

func EmployeeDeleteByUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	err := service.DeleteEmployeeByUID(r.Context(), service.DeleteEmployeeByUIDCommand(uid))
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Deleted staff uid=%s", uid),
	})
}

func EmployeeUpdateByUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	var payload service.UpdateEmployeeCommand
	if err := httputil.BindJSON(r, &payload); err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	err := service.UpdateEmployeeByUid(r.Context(), service.UpdateEmployeeByUIDCommand(uid), payload)
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Updated staff uid=%s", uid),
		Data:    payload,
	})
}

func EventHandler(w http.ResponseWriter, r *http.Request) {
	sub := httpsub.NewSubscriber("add-team")
	httpsub.ConnectSub(*sub, "add-team")
	httpsub.HTTPHandler(w, r)
}
