package cas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"net/http"
	"net/url"
	"strconv"
)

var RedirectUrl = ""
var LoginUri = ""

const REDIRECT_CODE = 401

func (c *Client) HandlerV1(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if glog.V(2) {
			glog.Infof("cas: handling %v request for %v", r.Method, r.URL)
		}

		setClient(r, c)

		if !IsAuthenticated(r) {
			genRedirectUrl := fmt.Sprintf("%s/login?service=%s", RedirectUrl, url.QueryEscape("http://"+r.Host+LoginUri))
			JSON(w, r, map[string]string{"code": strconv.Itoa(REDIRECT_CODE), "data": genRedirectUrl, "msg": "redirect login."})
			return
		}

		if r.URL.Path == "/logout" {
			RedirectToLogout(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

type contextKey struct {
	name string
}

var StatusCtxKey = &contextKey{"Status"}

func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if status, ok := r.Context().Value(StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write(buf.Bytes())
}
