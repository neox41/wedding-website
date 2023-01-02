package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"mrwebsites/config"
	"mrwebsites/db"
	"mrwebsites/middlewares"
	"mrwebsites/wedding"
	"net/http"
	"path/filepath"
)

func main() {
	flag.StringVar(&config.AdminUsername, "username", "admin", "Username for Admin portal")
	flag.StringVar(&config.AdminPassword, "password", "changeme_96xCf9Y33wUxBhcTywzF", "Password for Admin portal")
	flag.Parse()

	setupDatabase()

	mux := http.NewServeMux()

	Wedding(mux)

	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("data/certs"),
	}

	server := &http.Server{
		Addr:    config.Port,
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go func() {
		if err := http.ListenAndServe(":80", certManager.HTTPHandler(nil)); err != nil {
			panic(err)
		}
	}()
	if err := server.ListenAndServeTLS("", ""); err != nil {
		panic(err)
	}
}

func Wedding(mux *http.ServeMux) {
	// Redirect from WWW
	mux.HandleFunc(fmt.Sprintf("www.%s/", config.BaseURL), wedding.WWWRedirect)
	mux.HandleFunc(fmt.Sprintf("www.mattiaviviana.com"), wedding.WWWRedirect)
	mux.HandleFunc(fmt.Sprintf("mattiaviviana.com"), wedding.WWWRedirect)

	// Serving static files
	mux.Handle(fmt.Sprintf("%s/css/", config.BaseURL), middlewares.Static(http.FileServer(http.Dir(filepath.FromSlash(config.LocalSite)))))
	mux.Handle(fmt.Sprintf("%s/fonts/", config.BaseURL), middlewares.Static(http.FileServer(http.Dir(filepath.FromSlash(config.LocalSite)))))
	mux.Handle(fmt.Sprintf("%s/images/", config.BaseURL), middlewares.Static(http.FileServer(http.Dir(filepath.FromSlash(config.LocalSite)))))
	mux.Handle(fmt.Sprintf("%s/js/", config.BaseURL), middlewares.Static(http.FileServer(http.Dir(filepath.FromSlash(config.LocalSite)))))
	mux.Handle(fmt.Sprintf("%s/sass/", config.BaseURL), middlewares.Static(http.FileServer(http.Dir(filepath.FromSlash(config.LocalSite)))))
	mux.Handle(fmt.Sprintf("%s/admin/css", config.BaseURL), middlewares.Static(http.FileServer(http.Dir(filepath.FromSlash(config.LocalSite)))))
	mux.Handle(fmt.Sprintf("%s/admin/js", config.BaseURL), middlewares.Static(http.FileServer(http.Dir(filepath.FromSlash(config.LocalSite)))))
	mux.Handle(fmt.Sprintf("%s/admin/DataTables", config.BaseURL), middlewares.Static(http.FileServer(http.Dir(filepath.FromSlash(config.LocalSite)))))

	// Serving dynamic pages
	// Home and Default
	mux.HandleFunc(fmt.Sprintf("%s/", config.BaseURL), middlewares.Language(wedding.Home))
	mux.HandleFunc(fmt.Sprintf("%s/index.html", config.BaseURL), middlewares.Language(wedding.Home))
	mux.HandleFunc(fmt.Sprintf("%s/home.html", config.BaseURL), middlewares.Language(wedding.Home))
	mux.HandleFunc(fmt.Sprintf("%s/home", config.BaseURL), middlewares.Language(wedding.Home))

	// Gifting
	mux.HandleFunc(fmt.Sprintf("%s/gifting.html", config.BaseURL), middlewares.Language(wedding.Gifting))
	mux.HandleFunc(fmt.Sprintf("%s/gifting", config.BaseURL), middlewares.Language(wedding.Gifting))

	// Logistics
	mux.HandleFunc(fmt.Sprintf("%s/logistics.html", config.BaseURL), middlewares.Language(wedding.Logistics))
	mux.HandleFunc(fmt.Sprintf("%s/logistics", config.BaseURL), middlewares.Language(wedding.Logistics))

	// FAQ
	mux.HandleFunc(fmt.Sprintf("%s/faq.html", config.BaseURL), middlewares.Language(wedding.FAQ))
	mux.HandleFunc(fmt.Sprintf("%s/faq", config.BaseURL), middlewares.Language(wedding.FAQ))

	// Attendance
	mux.HandleFunc(fmt.Sprintf("%s/attendance.html", config.BaseURL), middlewares.Language(wedding.Attendance))
	mux.HandleFunc(fmt.Sprintf("%s/attendance", config.BaseURL), middlewares.Language(wedding.Attendance))
	mux.HandleFunc(fmt.Sprintf("%s/partecipazioni.html", config.BaseURL), middlewares.Language(wedding.Attendance))
	mux.HandleFunc(fmt.Sprintf("%s/partecipazioni", config.BaseURL), middlewares.Language(wedding.Attendance))

	// Admin
	adminPage := fmt.Sprintf("%s/admin", config.BaseURL)
	mux.HandleFunc(fmt.Sprintf("%s/css/", adminPage), middlewares.Static(http.FileServer(http.Dir(filepath.FromSlash(config.LocalSite)))))
	mux.HandleFunc(fmt.Sprintf("%s/DataTables/", adminPage), middlewares.Static(http.FileServer(http.Dir(filepath.FromSlash(config.LocalSite)))))
	mux.HandleFunc(fmt.Sprintf("%s/js/", adminPage), middlewares.Static(http.FileServer(http.Dir(filepath.FromSlash(config.LocalSite)))))
	mux.HandleFunc(fmt.Sprintf("%s/", adminPage), middlewares.Auth(wedding.Admin))
	mux.HandleFunc(adminPage, middlewares.Auth(AdminRedirect))
	mux.HandleFunc(fmt.Sprintf("%s/index.html", adminPage), middlewares.Auth(wedding.Admin))
	mux.HandleFunc(fmt.Sprintf("%s/attendees.html", adminPage), middlewares.Auth(wedding.Attendees))
	mux.HandleFunc(fmt.Sprintf("%s/families.html", adminPage), middlewares.Auth(wedding.Families))
	mux.HandleFunc(fmt.Sprintf("%s/import.html", adminPage), middlewares.Auth(wedding.Import))
	mux.HandleFunc(fmt.Sprintf("%s/export.html", adminPage), middlewares.Auth(wedding.Export))

}
func AdminRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s/", r.URL.Path), http.StatusMovedPermanently)
}
func setupDatabase() {
	if err := db.Setup(); err != nil {
		panic(err)
	}
}
