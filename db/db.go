package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"mrwebsites/config"
	mrerrors "mrwebsites/errors"
	"mrwebsites/internal"
	"path/filepath"
	"sync"
)

const ()

var (
	db     *sql.DB
	dbLock sync.RWMutex
)

const (
	CONFIRMED = "CONFIRMED"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func connect() *sql.DB {
	defer mrerrors.RecoverError()
	db, err := sql.Open("sqlite3", filepath.FromSlash(config.DB))
	mrerrors.CheckErrorPanic(err)
	err = db.Ping()
	mrerrors.CheckErrorPanic(err)
	return db
}
func close(database *sql.DB) {
	database.Close()
}

func NewRequest(name, email string) error {
	status := "new"
	dbLock.Lock()
	defer dbLock.Unlock()
	db = connect()
	defer close(db)
	stmt, err := db.Prepare("INSERT INTO requests (name, email, status) VALUES (?, ?, ?);")
	defer stmt.Close()
	mrerrors.CheckErrorPanic(err)
	_, err = stmt.Exec(name, email, status)
	mrerrors.CheckErrorPanic(err)

	return err
}
func Delete(table string) error {
	dbLock.Lock()
	defer dbLock.Unlock()
	db = connect()
	defer close(db)
	_, err := db.Exec(fmt.Sprintf("DELETE from %s;", table))
	mrerrors.CheckErrorPanic(err)
	return err
}
func Import(attendees []internal.AttendeesOutput, families []internal.FamiliesOutput) error {
	// Delete
	if err := Delete("attendees"); err != nil {
		return err
	}
	if err := Delete("families"); err != nil {
		return err
	}
	// Upload
	dbLock.Lock()
	defer dbLock.Unlock()
	db = connect()
	defer close(db)
	stmt, err := db.Prepare("INSERT INTO attendees (id, name, IDNucleo, surname, vegan, vegetarian, transport, attendance, category, abroad, bambino, beve, tavolo, requirements) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);")
	defer stmt.Close()
	mrerrors.CheckErrorPanic(err)
	for _, attendee := range attendees {
		_, err = stmt.Exec(attendee.ID, attendee.Name, attendee.IDNucleo, attendee.Surname, attendee.Vegan, attendee.Vegetarian,
			attendee.Transport, attendee.Attendance, attendee.Category, attendee.Abroad, attendee.Bambino,
			attendee.Beve, attendee.Table, attendee.Requirements)
	}

	// Families
	stmt, err = db.Prepare("INSERT INTO families (IDNucleo, Email, link, english, confirmed, name, status) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?);")
	mrerrors.CheckErrorPanic(err)
	for _, family := range families {
		_, err = stmt.Exec(family.IDNucleo, family.Email, family.Link, family.English, family.Replied, family.Name, family.Status)
	}
	return nil
}

func GetAttendees() []internal.AttendeesOutput {
	var attendees []internal.AttendeesOutput
	dbLock.Lock()
	defer dbLock.Unlock()
	db = connect()
	defer close(db)
	rows, err := db.Query("SELECT families.name, attendees.id, families.IDNucleo, attendees.name, attendees.surname, families.Email, families.status, families.link, families.english, attendees.vegan, attendees.vegetarian, attendees.requirements, attendees.transport, attendees.attendance, attendees.category, attendees.abroad, attendees.bambino, attendees.beve, attendees.tavolo " +
		"FROM attendees INNER JOIN families ON families.IDNucleo = attendees.IDNucleo " +
		"ORDER BY attendees.IDNucleo ASC;")
	check(err)
	var item internal.AttendeesOutput
	for rows.Next() {
		err = rows.Scan(&item.NucleoName, &item.ID, &item.IDNucleo, &item.Name, &item.Surname, &item.Email, &item.Status, &item.Link, &item.English,
			&item.Vegan, &item.Vegetarian, &item.Requirements, &item.Transport, &item.Attendance, &item.Category, &item.Abroad, &item.Bambino, &item.Beve, &item.Table)
		check(err)
		attendees = append(attendees, item)
	}
	return attendees
}
func GetFamilyAttendees(link string) []internal.AttendeesOutput {
	var attendees []internal.AttendeesOutput
	dbLock.Lock()
	defer dbLock.Unlock()
	db = connect()
	defer close(db)
	rows, err := db.Query("SELECT attendees.id, families.IDNucleo, attendees.name, attendees.surname, families.Email, families.status, families.link, families.english "+
		"FROM attendees INNER JOIN families ON families.IDNucleo = attendees.IDNucleo "+
		"WHERE families.link = ?;", link)
	check(err)
	var item internal.AttendeesOutput
	for rows.Next() {
		err = rows.Scan(&item.ID, &item.IDNucleo, &item.Name, &item.Surname, &item.Email, &item.Status, &item.Link, &item.English)
		check(err)
		attendees = append(attendees, item)
	}
	return attendees
}
func GetFamilies() []internal.FamiliesOutput {
	var families []internal.FamiliesOutput
	dbLock.Lock()
	defer dbLock.Unlock()
	db = connect()
	defer close(db)
	rows, err := db.Query("SELECT name, IDNucleo, Email, english, link, confirmed, status FROM families;")
	check(err)
	var item internal.FamiliesOutput
	for rows.Next() {
		err = rows.Scan(&item.Name, &item.IDNucleo, &item.Email, &item.English, &item.Link, &item.Replied, &item.Status)
		check(err)
		families = append(families, item)
	}
	return families
}
func UpdateStatus(status string, id int) error {
	dbLock.Lock()
	defer dbLock.Unlock()
	db = connect()
	defer close(db)
	stmt, err := db.Prepare("UPDATE families SET status = ? WHERE IDNucleo = ?;")
	defer stmt.Close()
	check(err)
	_, err = stmt.Exec(status, id)
	check(err)
	return err
}
func UpdateConfirmation(attendees []internal.Attendee) error {
	dbLock.Lock()
	defer dbLock.Unlock()
	db = connect()
	defer close(db)
	var (
		err  error
		stmt *sql.Stmt
	)
	for _, attendee := range attendees {
		stmt, err = db.Prepare("UPDATE attendees SET vegan = ?, vegetarian = ?, requirements = ?, transport = ?, attendance = ? WHERE id = ?;")
		defer stmt.Close()
		check(err)
		_, err = stmt.Exec(attendee.Vegan, attendee.Vegetarian, attendee.Requirements, attendee.Transport, attendee.Attendance, attendee.ID)
		check(err)
	}
	return err
}
func UpdateLink(link string) error {
	dbLock.Lock()
	defer dbLock.Unlock()
	db = connect()
	defer close(db)
	confirmed := true
	stmt, err := db.Prepare("UPDATE families SET confirmed = ? WHERE link = ?;")
	defer stmt.Close()
	check(err)
	_, err = stmt.Exec(confirmed, link)
	check(err)
	return err
}
func UpdateFamily(family internal.UpdateFamily) error {
	dbLock.Lock()
	defer dbLock.Unlock()
	db = connect()
	defer close(db)
	stmt, err := db.Prepare("UPDATE families SET Email = ? WHERE IDNucleo = ?;")
	defer stmt.Close()
	check(err)
	_, err = stmt.Exec(family.Email, family.Id)
	check(err)
	return err
}
func IsConfirmed(link string) bool {
	dbLock.Lock()
	defer dbLock.Unlock()
	db = connect()
	defer close(db)
	var confirmed bool

	row := db.QueryRow("SELECT confirmed FROM families WHERE link = ? LIMIT 1;", link)
	if err := row.Scan(&confirmed); err == nil {
		return confirmed
	}
	return false
}
