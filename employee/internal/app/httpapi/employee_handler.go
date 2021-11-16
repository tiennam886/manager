package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis"
	"github.com/tiennam886/manager/employee/internal/service"
	"github.com/tiennam886/manager/pkg/httputil"

	"net/http"
)

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

func EmployeeFindByUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	staff, err := service.FindStaffByUID(r.Context(), service.FindEmployeeByUIDCommand(uid))
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Found staff uid=%s", uid),
		Data:    staff,
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

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func EmployeeAddToTeam(w http.ResponseWriter, r *http.Request) {
	eid := chi.URLParam(r, "eid")
	tid := chi.URLParam(r, "tid")

	payload, err := json.Marshal(eid + tid)
	if err != nil {
		panic(err)
	}

	if err := redisClient.Publish("send-user-data", payload).Err(); err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Add employee eid=%s to team tid=%s", eid, tid),
	})
}
