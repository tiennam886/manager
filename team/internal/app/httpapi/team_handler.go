package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis"
	"github.com/tiennam886/manager/pkg/httputil"
	"github.com/tiennam886/manager/team/internal/service"
	"net/http"
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

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func TeamNotice(w http.ResponseWriter, r *http.Request) {
	subscriber := redisClient.Subscribe("send-user-data")

	var user string
	msg, err := subscriber.ReceiveMessage()
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(msg.Payload), &user); err != nil {
		panic(err)
	}
	_ = httputil.WriteJsonOK(w, httputil.ResponseBody{
		Message: fmt.Sprintf("Deleted staff uid=%s", user),
	})

	fmt.Println("Received message from " + msg.Channel + " channel.")
	fmt.Printf("%+v\n", user)
}
