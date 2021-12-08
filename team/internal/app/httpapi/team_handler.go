package httpapi

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/tiennam886/manager/pkg/httputil"
	"github.com/tiennam886/manager/pkg/messaging"
	"github.com/tiennam886/manager/pkg/messaging/httppub"
	"github.com/tiennam886/manager/team/internal/service"
)

var sugarLogger *zap.SugaredLogger

func TeamGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	sugarLogger.Infow("GET /team")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	teams, err := service.GetAllTeam(r.Context(), page, limit)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	sugarLogger.Infow("Get All Teams Successfully")
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: "All Employees",
		Data:    teams,
	})
}

func TeamAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	var payload service.AddTeamCommand

	sugarLogger.Infow("POST /team")

	if err := httputil.BindJSON(r, &payload); err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	employee, err := service.AddTeam(r.Context(), payload)
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Added team uid=%s", employee.UID),
		Data:    employee,
	})
}

func TeamFindByUID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	uid := chi.URLParam(r, "uid")

	sugarLogger.Infow(fmt.Sprintf("GET /team/%s", uid))

	staff, err := service.FindTeamByUID(r.Context(), service.FindTeamByUIDCommand(uid))
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	sugarLogger.Infow(fmt.Sprintf("Found team uid= %s", uid))
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Found staff uid=%s", uid),
		Data:    staff,
	})
}

func TeamDeleteByUID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	uid := chi.URLParam(r, "uid")

	sugarLogger.Infow(fmt.Sprintf("DELETE /team/%s", uid))

	err := service.DeleteTeamByUID(r.Context(), service.DeleteTeamByUIDCommand(uid))
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	sugarLogger.Infow(fmt.Sprintf("Team with ID %s was deleted\n", uid))
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Deleted team uid=%s", uid),
	})
}

func TeamUpdateByUID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	uid := chi.URLParam(r, "uid")
	sugarLogger.Infow(fmt.Sprintf("UPDATE /team/%s", uid))

	var payload service.UpdateTeamCommand
	if err := httputil.BindJSON(r, &payload); err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	err := service.UpdateTeamByUid(r.Context(), service.UpdateTeamByUIDCommand(uid), payload)
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	sugarLogger.Infow(fmt.Sprintf("Updated team uid=%s", uid))
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Updated team uid=%s", uid),
		Data:    payload,
	})
}

type Command struct {
	Event string
	Json  []byte
}

func (c Command) Name() string {
	return c.Event
}

func (c Command) JSON() []byte {
	return c.Json
}

func TeamNotice(w http.ResponseWriter, r *http.Request) {
	pub := httppub.NewPublisher(
		"add-team",
		url.URL{
			Host: "localhost:8080",
			Path: "/api/v1/employee/event",
		},
		r.Header)
	err := httppub.ConnectPub(*pub, "add-team")
	if err != nil {
		fmt.Println(err.Error())
	}
	var event messaging.Event
	event = &Command{
		Event: "add-team",
		Json:  []byte{},
	}
	httppub.Publish(event)
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Data: event,
	})
}

func TeamOption(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: "",
		Data:    nil,
	})
}
