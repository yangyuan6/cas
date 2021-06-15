package cas

import (
	"fmt"
	"github.com/go-chi/render"
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
			//c := getClient(r)
			//if c == nil {
			//	err := "cas: redirect to cas failed as no client associated with request"
			//	http.Error(w, err, http.StatusInternalServerError)
			//	return
			//}
			//u, err := c.LoginUrlForRequest(r)
			//if err != nil {
			//	http.Error(w, err.Error(), http.StatusInternalServerError)
			//	return
			//}
			//
			//if glog.V(2) {
			//	glog.Infof("Redirecting client to %v with status %v", u, http.StatusFound)
			//}
			//RedirectToLogin(w, r)
			genRedirectUrl := fmt.Sprintf("%s/login?service=%s", RedirectUrl, url.QueryEscape("http://"+r.Host+LoginUri))
			render.JSON(w, r, map[string]string{"code": strconv.Itoa(REDIRECT_CODE), "data": genRedirectUrl, "msg": "redirect login."})

			return
		}

		if r.URL.Path == "/logout" {
			RedirectToLogout(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
