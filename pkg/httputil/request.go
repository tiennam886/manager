package httputil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func BindJSON(r *http.Request, obj interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, obj)
}
