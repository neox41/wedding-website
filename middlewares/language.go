package middlewares

import (
	"context"
	"mrwebsites/security"
	"net/http"
)

func Language(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		english := false
		lang := request.URL.Query().Get("lang")
		if len(lang) > 0 {
			// Get lang setting via GET param
			if lang == "en"{
				english = true
			}
		}else{
			// No GET param, get default setting from IP address
			if !security.IsForItalian(request.RemoteAddr) {
				english = true
			}
		}

		ctx := context.WithValue(request.Context(), "english", english)
		request = request.WithContext(ctx)
		next(w, request)
	}
}
