package wedding

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"io/ioutil"
	"mrwebsites/config"
	"mrwebsites/db"
	"mrwebsites/internal"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
)

func Admin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/admin/index.html", config.LocalSite))
	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, nil)
}

func Attendees(w http.ResponseWriter, request *http.Request) {
	var (
		tmpl *template.Template
		err  error
	)

	if request.Method == "GET" {
		action := request.URL.Query().Get("action")
		if len(action) < 1 {
			goto SERVEDEFAULT
		}
		if action == "getAttendees" {
			attendees := db.GetAttendees()
			// Setting the link
			for i, _ := range attendees {
				if !attendees[i].English {
					// Italian
					attendees[i].Link = fmt.Sprintf("%s%s%s/partecipazioni.html?lang=it&id=%s", config.Proto, config.BaseURL, config.Port, attendees[i].Link)
				} else {
					// English
					attendees[i].Link = fmt.Sprintf("%s%s%s/attendance.html?lang=en&id=%s", config.Proto, config.BaseURL, config.Port, attendees[i].Link)
				}
			}
			serveJSON(w, attendees)
			return
		}
	}

	if request.Method == "POST" {
		updateAttendeeRequest, err := ioutil.ReadAll(request.Body)
		if err != nil {
			serveResponse(w, err.Error())
			return
		}
		var updateAttendee internal.UpdateStatus
		if err := json.Unmarshal(updateAttendeeRequest, &updateAttendee); err != nil {
			serveResponse(w, err.Error())
			return
		}
		if updateAttendee.Status != "TO SEND" && updateAttendee.Status != "SENT" {
			serveResponse(w, "Invalid status")
			return
		}
		if err := db.UpdateStatus(updateAttendee.Status, updateAttendee.Id); err != nil {
			serveResponse(w, err.Error())
			return
		}
		// All good, updated!
		serveResponse(w, "Updated")
		return
	}

SERVEDEFAULT:
	tmpl, err = template.ParseFiles(fmt.Sprintf("%s/admin/attendees.html", config.LocalSite))
	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, nil)
}
func Families(w http.ResponseWriter, request *http.Request) {
	var (
		tmpl *template.Template
		err  error
	)

	if request.Method == "GET" {
		action := request.URL.Query().Get("action")
		if len(action) < 1 {
			goto SERVEDEFAULT
		}
		if action == "getFamilies" {
			families := db.GetFamilies()
			// Setting the link

			for i, _ := range families {
				if families[i].English == 0 {
					// Italian
					families[i].Link = fmt.Sprintf("%s%s%s/partecipazioni.html?lang=it&id=%s", config.Proto, config.BaseURL, config.Port, families[i].Link)
				} else {
					// English
					families[i].Link = fmt.Sprintf("%s%s%s/attendance.html?lang=en&id=%s", config.Proto, config.BaseURL, config.Port, families[i].Link)
				}
			}
			serveJSON(w, families)
			return
		}
	}

	if request.Method == "POST" {
		updateFamilyRequest, err := ioutil.ReadAll(request.Body)
		if err != nil {
			serveResponse(w, err.Error())
			return
		}
		var updateFamily internal.UpdateFamily
		if err := json.Unmarshal(updateFamilyRequest, &updateFamily); err != nil {
			serveResponse(w, err.Error())
			return
		}
		if _, err := mail.ParseAddress(updateFamily.Email); err != nil {
			serveResponse(w, "Invalid email")
			return
		}
		if err := db.UpdateFamily(updateFamily); err != nil {
			serveResponse(w, err.Error())
			return
		}
		// All good, updated!
		serveResponse(w, "Updated")
		return
	}

SERVEDEFAULT:
	tmpl, err = template.ParseFiles(fmt.Sprintf("%s/admin/families.html", config.LocalSite))
	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, nil)
}

func Import(w http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		serveResponse(w, "Method not supported")
		return
	}

	if err := request.ParseMultipartForm(10 << 20); err != nil {
		serveResponse(w, err.Error())
		return
	}

	action := request.FormValue("action")
	if action != "import" {
		serveResponse(w, "Invalid action")
		return
	}
	file, _, err := request.FormFile("file")
	if err != nil {
		serveResponse(w, "File empty")
		return
	}
	defer file.Close()
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		serveResponse(w, err.Error())
		return
	}
	var (
		attendees []internal.AttendeesOutput
		families  []internal.FamiliesOutput
	)
	for i, line := range lines {
		if i == 0 {
			continue
		}
		id, err := strconv.Atoi(line[0])
		if err != nil {
			serveResponse(w, err.Error())
			return
		}
		idNucleo, err := strconv.Atoi(line[1])
		if err != nil {
			serveResponse(w, err.Error())
			return
		}
		english := 0
		if line[4] != "IT" {
			english = 1
		}
		abroad, err := strconv.ParseBool(line[8])
		if err != nil {
			serveResponse(w, err.Error())
			return
		}
		bambino, err := strconv.ParseBool(line[9])
		if err != nil {
			serveResponse(w, err.Error())
			return
		}

		beve, err := strconv.ParseBool(line[11])
		if err != nil {
			serveResponse(w, err.Error())
			return
		}

		replied, err := strconv.ParseBool(line[13])
		if err != nil {
			serveResponse(w, err.Error())
			return
		}

		partecipazioneInviata, err := strconv.ParseBool(line[12])
		if err != nil {
			serveResponse(w, err.Error())
			return
		}
		status := "TO SEND"
		if partecipazioneInviata {
			status = "SENT"
		}
		attendee := internal.AttendeesOutput{
			ID:       id,
			IDNucleo: idNucleo,
			Name:     line[2],
			Surname:  line[3],

			// line [4] is invito sicuro
			Category:     line[6],
			Requirements: line[7],
			Abroad:       abroad,
			Bambino:      bambino,
			Beve:         beve,
			Attendance:   line[14],
			Table:        line[15],
		}
		attendees = append(attendees, attendee)
		newFamily := true
		for _, f := range families {
			if f.IDNucleo == idNucleo {
				newFamily = false
				break
			}
		}
		if newFamily {
			family := internal.FamiliesOutput{
				IDNucleo: idNucleo,
				Name:     line[5],
				Replied:  replied,
				Status:   status,
				English:  english,
				Link:     uuid.New().String(),
				Email:    "N/A",
			}
			families = append(families, family)
		}
	}

	if err := db.Import(attendees, families); err != nil {
		serveResponse(w, err.Error())
		return
	}
	// All good, imported!
	serveResponse(w, "ok")
}

func Export(w http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		serveResponse(w, "Method not supported")
		return
	}
	attendees := db.GetAttendees()
	var exportBuf bytes.Buffer
	exportCSV := csv.NewWriter(&exportBuf)
	headers := "ID,ID Nucleo,Nome,Cognome,Lingua Inglese,Nucleo Name,Categoria,Requirements,Abroad,Bambino,Beve,Vegan,Vegetarian,Requirements,Transport,Partecipazione inviata,Replied,Attendance,Tavolo"
	if err := exportCSV.Write(strings.Split(headers, ",")); err != nil {
		serveResponse(w, err.Error())
		return
	}

	for _, attendee := range attendees {
		line := []string{
			fmt.Sprintf("%d", attendee.ID),
			fmt.Sprintf("%d", attendee.IDNucleo),
			attendee.Name,
			attendee.Surname,
			strconv.FormatBool(attendee.English),
			attendee.NucleoName,
			attendee.Category,
			attendee.Requirements,
			strconv.FormatBool(attendee.Abroad),
			strconv.FormatBool(attendee.Bambino),
			strconv.FormatBool(attendee.Beve),
			strconv.FormatBool(attendee.Vegan),
			strconv.FormatBool(attendee.Vegetarian),
			attendee.Requirements,
			strconv.FormatBool(attendee.Transport),
			attendee.Status,
			strconv.FormatBool(attendee.Replied),
			attendee.Attendance,
			attendee.Table,
		}
		if err := exportCSV.Write(line); err != nil {
			serveResponse(w, err.Error())
			return
		}
	}
	exportCSV.Flush()
	if err := exportCSV.Error(); err != nil {
		serveResponse(w, err.Error())
		return
	}
	w.Header().Set("Content-Encoding", "UTF-8")
	w.Header().Set("Content-type", "text/csv charset=UTF-8")
	w.Header().Set("Content-Disposition", "attachment; filename=ExportListaInvitati.csv")
	w.WriteHeader(http.StatusOK)
	w.Write(exportBuf.Bytes())
}
func serveJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
func serveResponse(w http.ResponseWriter, response string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
