package wedding

import (
	"fmt"
	"html/template"
	"mrwebsites/config"
	"mrwebsites/internal"
	"net/http"
)

var WeddingWebsite internal.Pages

func init() {
	WeddingWebsite = internal.Pages{
		Error:   "",
		English: true,
		Pages: []internal.PageInfo{
			{
				TitleEN:  fmt.Sprintf("%s | Home", config.TitleEN),
				TitleIT:  fmt.Sprintf("%s | Home", config.TitleIT),
				NameEN:   "Home",
				NameIT:   "Home",
				Location: "index.html",
				Active:   false,
			},
			{
				TitleEN:  fmt.Sprintf("%s | Logistics", config.TitleEN),
				TitleIT:  fmt.Sprintf("%s | Logistica", config.TitleIT),
				NameEN:   "Getting There",
				NameIT:   "Come Arrivare",
				Location: "logistics.html",
				Active:   false,
			},
			{
				TitleEN:  fmt.Sprintf("%s | Gifting", config.TitleEN),
				TitleIT:  fmt.Sprintf("%s | Regali", config.TitleIT),
				NameEN:   "Gifting",
				NameIT:   "Lista Nozze",
				Location: "gifting.html",
				Active:   false,
			},

			{
				TitleEN:  fmt.Sprintf("%s | Attendance", config.TitleEN),
				TitleIT:  fmt.Sprintf("%s | Partecipazioni", config.TitleIT),
				NameEN:   "Attendance",
				NameIT:   "Conferma presenza",
				Location: "attendance.html",
				Active:   false,
			},

			{
				TitleEN:  fmt.Sprintf("%s | FAQ", config.TitleEN),
				TitleIT:  fmt.Sprintf("%s | FAQ", config.TitleIT),
				NameEN:   "FAQ",
				NameIT:   "FAQ",
				Location: "faq.html",
				Active:   false,
			},
		},
	}
}

func ServeStaticPage(w http.ResponseWriter, request *http.Request, name string) {
	var (
		location     string
		WeddingPages = WeddingWebsite
	)
	for page := range WeddingPages.Pages {
		if WeddingPages.Pages[page].NameIT == name || WeddingPages.Pages[page].NameEN == name {
			WeddingPages.Pages[page].Active = true
			location = WeddingPages.Pages[page].Location
		} else {
			WeddingPages.Pages[page].Active = false
		}
	}

	// Setting the language
	if request.Context().Value("english").(bool) {
		WeddingPages.English = true
	} else {
		WeddingPages.English = false
	}

	//tmpl := template.Must(template.ParseFiles(fmt.Sprintf("%s/index.html", LocalSite)))
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/%s", config.LocalSite, location), fmt.Sprintf("%s/nav.html", config.LocalSite), fmt.Sprintf("%s/title.html", config.LocalSite))
	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, WeddingPages)
}
func Home(w http.ResponseWriter, r *http.Request) {
	ServeStaticPage(w, r, "Home")
}
func Logistics(w http.ResponseWriter, r *http.Request) {
	ServeStaticPage(w, r, "Getting There")
}
func Gifting(w http.ResponseWriter, r *http.Request) {
	ServeStaticPage(w, r, "Gifting")
}
func FAQ(w http.ResponseWriter, r *http.Request) {
	ServeStaticPage(w, r, "FAQ")
}
