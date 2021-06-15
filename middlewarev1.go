package cas

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/golang/glog"
	"net/http"
	"strconv"
)

var RedirectUrl = ""
var LoginUri = ""

const REDIRECT_CODE = 1088

func (c *Client) HandlerV1(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if glog.V(2) {
			glog.Infof("cas: handling %v request for %v", r.Method, r.URL)
		}

		setClient(r, c)

		if !IsAuthenticated(r) {
			genRedirectUrl := fmt.Sprintf("%s/login??service=%s%s", RedirectUrl, r.Host, LoginUri)
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
