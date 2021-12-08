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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	sugarLogger.Infow("POST /employee")
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	sugarLogger.Infow("GET /employee")
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
		sugarLogger.Errorf(err.Error())
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	sugarLogger.Infow("Get All Employees Successfully")
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: "All Employees",
		Data:    employees,
	})
}

func EmployeeFindByUID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	uid := chi.URLParam(r, "uid")
	sugarLogger.Infow(fmt.Sprintf("GET /employee/%s", uid))

	employee, err := service.FindStaffByUID(r.Context(), service.FindEmployeeByUIDCommand(uid))
	if err != nil {
		sugarLogger.Errorf(err.Error())
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	sugarLogger.Infow(fmt.Sprintf("Found staff uid= %s", uid))
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Found staff uid=%s", uid),
		Data:    employee,
	})
}

func EmployeeFindTeams(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	uid := chi.URLParam(r, "uid")
	sugarLogger.Infow(fmt.Sprintf("GET /employee/list/%s", uid))

	teamList, err := service.FindEmployeesTeams(r.Context(), service.FindEmployeeByUIDCommand(uid))
	if err != nil {
		sugarLogger.Errorf(err.Error())
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	sugarLogger.Infow(fmt.Sprintf("Found teams of staff uid= %s", uid))
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Found teams of staff uid=%s", uid),
		Data:    teamList,
	})
}

func EmployeeDeleteByUID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	uid := chi.URLParam(r, "uid")
	sugarLogger.Infow(fmt.Sprintf("DELETE /employee/%s", uid))

	err := service.DeleteEmployeeByUID(r.Context(), service.DeleteEmployeeByUIDCommand(uid))
	if err != nil {
		sugarLogger.Errorf(err.Error())
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	sugarLogger.Infow(fmt.Sprintf("Employee with ID %s was deleted\n", uid))
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Deleted staff uid=%s", uid),
	})
}

func EmployeeUpdateByUID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	uid := chi.URLParam(r, "uid")
	sugarLogger.Infow(fmt.Sprintf("UPDATE /employee/%s", uid))

	var payload service.UpdateEmployeeCommand
	if err := httputil.BindJSON(r, &payload); err != nil {
		sugarLogger.Errorf(err.Error())
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	err := service.UpdateEmployeeByUid(r.Context(), service.UpdateEmployeeByUIDCommand(uid), payload)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	sugarLogger.Infow(fmt.Sprintf("Updated staff uid=%s", uid))
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Updated staff uid=%s", uid),
		Data:    payload,
	})
}

func EmployeeAddToTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	payload := service.EmployeeToTeamCommand{
		EmployeeId: chi.URLParam(r, "eid"),
		TeamId:     chi.URLParam(r, "tid"),
	}

	err := service.AddEmployeeToTeam(r.Context(), payload)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	sugarLogger.Infow(fmt.Sprintf("staff uid=%s was added to team uid=%s", payload.EmployeeId, payload.TeamId))
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("staff uid=%s was added to team uid=%s", payload.EmployeeId, payload.TeamId),
		Data:    payload,
	})
}

func EmployeeRemoveFromTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	payload := service.EmployeeToTeamCommand{
		EmployeeId: chi.URLParam(r, "eid"),
		TeamId:     chi.URLParam(r, "tid"),
	}

	err := service.DeleteEmployeeToTeam(r.Context(), payload)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	sugarLogger.Infow(fmt.Sprintf("staff uid=%s was removed from team uid=%s", payload.EmployeeId, payload.TeamId))
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("staff uid=%s was removed from team uid=%s", payload.EmployeeId, payload.TeamId),
		Data:    payload,
	})
}

func EmployeeOption(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: "",
		Data:    nil,
	})
}

func EventHandler(w http.ResponseWriter, r *http.Request) {
	sub := httpsub.NewSubscriber("add-team")
	httpsub.ConnectSub(*sub, "add-team")
	httpsub.HTTPHandler(w, r)
}
