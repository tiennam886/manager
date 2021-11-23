package httpapi

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/tiennam886/manager/pkg/httputil"
	"github.com/tiennam886/manager/pkg/messaging"
	"github.com/tiennam886/manager/pkg/messaging/httppub"
	"github.com/tiennam886/manager/team/internal/service"
)

func TeamAdd(w http.ResponseWriter, r *http.Request) {
	var payload service.AddTeamCommand

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
		Message: fmt.Sprintf("Added staff uid=%s", employee.UID),
		Data:    employee,
	})
}

func TeamFindByUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	staff, err := service.FindTeamByUID(r.Context(), service.FindTeamByUIDCommand(uid))
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Found staff uid=%s", uid),
		Data:    staff,
	})
}

func TeamDeleteByUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	err := service.DeleteTeamByUID(r.Context(), service.DeleteTeamByUIDCommand(uid))
	if err != nil {
		httputil.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Deleted staff uid=%s", uid),
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
