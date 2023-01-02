package wedding

import (
	"fmt"
	"mrwebsites/config"
	"net/http"
)

func WWWRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s%s%s", config.Proto, config.BaseURL, r.URL.Path), http.StatusMovedPermanently)
}
