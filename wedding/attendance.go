package wedding

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"html/template"
	"mrwebsites/config"
	"mrwebsites/db"
	"mrwebsites/internal"
	"mrwebsites/security"
	"net/http"
	"net/mail"
	"time"
)

var cacheAttendance *cache.Cache

func init() {
	cacheAttendance = cache.New(5*time.Minute, 10*time.Minute)
}
func Attendance(w http.ResponseWriter, request *http.Request) {
	var (
		tmpl         *template.Template
		WeddingPages = WeddingWebsite
		err          error
		name         = "Attendance"
		location     string
		idLink       string
		resultPage   = fmt.Sprintf("%s/attendance_default.html", config.LocalSite)
		errorPage    = fmt.Sprintf("%s/attendance_noerror.html", config.LocalSite)
	)
	WeddingPages.Error = ""
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

	if request.Method == "POST" {
		if err := request.ParseForm(); err != nil {
			WeddingPages.Error = "Parsing form"
			goto SERVETEMPLATE
		}
		action := request.FormValue("action")
		if action == "" || (action != "register" && action != "confirm") {
			WeddingPages.Error = "Invalid Action"
			goto SERVETEMPLATE
		}

		// Check if this request comes from a default page
		if action == "register" {
			namePartecipant := request.FormValue("name")
			if namePartecipant == "" {
				WeddingPages.Error = "Empty Name"
				goto SERVETEMPLATE
			}
			email := request.FormValue("email")
			if email == "" {
				WeddingPages.Error = "Empty Email"
				goto SERVETEMPLATE
			} else {
				_, err := mail.ParseAddress(email)
				if err != nil {
					WeddingPages.Error = "Invalid Email"
					goto SERVETEMPLATE
				}
			}

			// Check limit
			if result, message := security.Attendance(request.RemoteAddr); !result {
				WeddingPages.Error = message
				goto SERVETEMPLATE
			}

			if err := db.NewRequest(namePartecipant, email); err != nil {
				if WeddingPages.English {
					WeddingPages.Error = "Unable to register your request. Contact Mattia or Viviana!"
				} else {
					WeddingPages.Error = "Servizio non disponibile. Contatta Mattia or Viviana!"
				}
				goto SERVETEMPLATE
			}

			// All good, registered!
			resultPage = fmt.Sprintf("%s/attendance_requested.html", config.LocalSite)
		}
		// Check if this request comes from the confirmation
		if action == "confirm" {
			// Get ID
			idLink = request.FormValue("idLink")
			if idLink == "" {
				WeddingPages.Error = "Empty ID"
				goto SERVECONFIRM
			}
			// Check if it's valid
			familyAttendees := db.GetFamilyAttendees(idLink)
			if len(familyAttendees) < 1 {
				WeddingPages.Error = "Invalid ID"
				goto SERVECONFIRM
			}
			// All good, valid ID, get all values
			var (
				attendees []internal.Attendee
			)
			for _, attendee := range familyAttendees {
				var (
					vegan, vegetarian, transport bool
					requirements, attendance     string
				)
				// Check attendance (attendance4=no)
				attendance = request.FormValue(fmt.Sprintf("attendance%d", attendee.ID))
				if attendance == "" {
					if WeddingPages.English {
						WeddingPages.Error = fmt.Sprintf("Missing attendance for %s %s", attendee.Name, attendee.Surname)
					} else {
						WeddingPages.Error = fmt.Sprintf("Partecipazione mancante per %s %s", attendee.Name, attendee.Surname)
					}
					goto SERVECONFIRM
				}
				// Check Vegan
				if r := request.FormValue(fmt.Sprintf("vegan%d", attendee.ID)); len(r) > 0 {
					vegan = true
				}
				// Check Vegetarian
				if r := request.FormValue(fmt.Sprintf("vegetarian%d", attendee.ID)); len(r) > 0 {
					vegetarian = true
				}
				// Check Transport
				if r := request.FormValue(fmt.Sprintf("transport%d", attendee.ID)); len(r) > 0 {
					transport = true
				}
				// Check Requirements
				if r := request.FormValue(fmt.Sprintf("requirements%d", attendee.ID)); len(r) > 0 {
					if len(r) > 1000000 {
						WeddingPages.Error = "'Requirements' text box too long"
						goto SERVECONFIRM
					}
					requirements = r
				}
				attendees = append(attendees, internal.Attendee{
					ID:           attendee.ID,
					Name:         attendee.Name,
					Surname:      attendee.Surname,
					Vegan:        vegan,
					Vegetarian:   vegetarian,
					Requirements: requirements,
					Transport:    transport,
					Attendance:   attendance,
				})
			}

			// All good, all confirmed. Let's insert in the database
			if err := db.UpdateConfirmation(attendees); err != nil {
				if WeddingPages.English {
					WeddingPages.Error = "Unable to confirm your attendance. Contact Mattia or Viviana!"
				} else {
					WeddingPages.Error = "Servizio non disponibile. Contatta Mattia or Viviana!"
				}
				goto SERVECONFIRM
			}

			if err := db.UpdateLink(idLink); err != nil {
				if WeddingPages.English {
					WeddingPages.Error = "Unable to confirm your attendance. Contact Mattia or Viviana!"
				} else {
					WeddingPages.Error = "Servizio non disponibile. Contatta Mattia or Viviana!"
				}
				goto SERVECONFIRM
			}

			// All good, confirmed!
			resultPage = fmt.Sprintf("%s/attendance_confirmed.html", config.LocalSite)
			goto SERVETEMPLATE
		}

	}
	if request.Method == "GET" {
		id := request.URL.Query().Get("id")
		if len(id) < 1 {
			goto SERVETEMPLATE
		}
		// Check limit
		if result, message := security.Link(request.RemoteAddr); !result {
			WeddingPages.Error = message
			goto SERVETEMPLATE
		}
		if isConfirmed := db.IsConfirmed(id); isConfirmed {
			resultPage = fmt.Sprintf("%s/attendance_confirmed.html", config.LocalSite)
			goto SERVETEMPLATE
		}
		familyAttendees := db.GetFamilyAttendees(id)
		if len(familyAttendees) < 1 {
			goto SERVETEMPLATE
		}
		// All good, provide confirmation details
		security.LogToFile(fmt.Sprintf("Requested details for %s Link ID (%s %s)", id, request.RemoteAddr, request.UserAgent()))
		WeddingPages.IdInvite = id
		if WeddingPages.English {
			WeddingPages.LinkLang = fmt.Sprintf("%s?lang=it&id=%s", location, id)
		} else {
			WeddingPages.LinkLang = fmt.Sprintf("%s?lang=en&id=%s", location, id)
		}

		var attendees []internal.Attendee
		for _, familyAttendee := range familyAttendees {
			attendee := internal.Attendee{
				ID:      familyAttendee.ID,
				Name:    familyAttendee.Name,
				Surname: familyAttendee.Surname,
			}
			attendees = append(attendees, attendee)
		}
		WeddingPages.Attendees = attendees
		resultPage = fmt.Sprintf("%s/attendance_confirm.html", config.LocalSite)
		tmpl, err = template.ParseFiles(fmt.Sprintf("%s/%s", config.LocalSite, location), fmt.Sprintf("%s/nav.html", config.LocalSite), fmt.Sprintf("%s/title.html", config.LocalSite), resultPage, errorPage)
		if err != nil {
			panic(err)
		}
		tmpl.Execute(w, WeddingPages)
		return
	}

SERVETEMPLATE:
	if WeddingPages.Error != "" {
		errorPage = fmt.Sprintf("%s/attendance_error.html", config.LocalSite)
	}
	tmpl, err = template.ParseFiles(fmt.Sprintf("%s/%s", config.LocalSite, location), fmt.Sprintf("%s/nav.html", config.LocalSite), fmt.Sprintf("%s/title.html", config.LocalSite), resultPage, errorPage)
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, WeddingPages)
	return

SERVECONFIRM:
	if WeddingPages.Error != "" {
		errorPage = fmt.Sprintf("%s/attendance_error.html", config.LocalSite)
	}
	resultPage = fmt.Sprintf("%s/attendance_confirm.html", config.LocalSite)
	tmpl, err = template.ParseFiles(fmt.Sprintf("%s/%s", config.LocalSite, location), fmt.Sprintf("%s/nav.html", config.LocalSite), fmt.Sprintf("%s/title.html", config.LocalSite), resultPage, errorPage)
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, WeddingPages)
	return
}
