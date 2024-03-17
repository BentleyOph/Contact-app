package controllers

import (
	// "encoding/json"
	"fmt"
	"github.com/BentleyOph/htmx-go/pkg/models"
	// "github.com/BentleyOph/htmx-go/pkg/utils"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	// "time"
)

var errorCheck = func(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var NewContact models.Contact

func RedirectRootToContacts(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func CreateContact(w http.ResponseWriter, r *http.Request) {
	contact := &models.Contact{}
	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	if firstName != "" && lastName != "" && email != "" && phone != "" {
		contact.FirstName = firstName
		contact.LastName = lastName
		contact.Email = email
		contact.Phone = phone
		contact.CreateContact()
		http.Redirect(w, r, "/contacts", http.StatusCreated)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		tmpl := template.Must(template.ParseFiles(wd + "/../../static/newcontact.html"))
		err = tmpl.Execute(w, nil)
		errorCheck(err)
	}
}

func GetContacts(w http.ResponseWriter, r *http.Request) {
	// Get the search term from the request
	searchTerm := r.URL.Query().Get("search")
	fmt.Println("Search term: ", searchTerm)
	// Pass the search term to the GetAllContacts function
	contacts := models.GetAllContacts(searchTerm)
	fmt.Println("Contacts: ", contacts)
	data := struct {
		SearchTerm string
		Contacts   []models.Contact
	}{
		SearchTerm: searchTerm,
		Contacts:   contacts,
	}

	// If the request header HX-Trigger is equal to “search” we want to do something different
	if r.Header.Get("HX-Trigger") == "search" {
		// Return as a response the rows.html template with the contacts only passed in
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		tmpl := template.Must(template.ParseFiles(wd + "/../../static/rows2.html"))
		// fmt.Println("Search triggered")
		err = tmpl.Execute(w, data)
		errorCheck(err)
		return

	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("wd: ", wd)
	// Parse the index.html file and the rows.html file
	tmpl := template.Must(template.ParseFiles(wd+"/../../static/index.html", wd+"/../../static/rows.html", wd+"/../../static/archive-ui.html"))
	errorCheck(err)

	// Pass the contacts and the search term to the template

	err = tmpl.Execute(w, data)
	errorCheck(err)
}

func GetContactsCount(w http.ResponseWriter, r *http.Request) {
	//just return a string indicating the total number of contacts in the db
	contactsCount := models.GetContactsCount()
	fmt.Fprintf(w, "(Total contacts: %d)", contactsCount)
}

func GetContactById(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	vars := mux.Vars(r)
	contactId := vars["id"]
	ID, err := strconv.ParseInt(contactId, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing red")
	}
	contactDetails, _ := models.GetContactById(ID)
	tmpl := template.Must(template.ParseFiles(wd + "/../../static/contactdetails.html"))
	err = tmpl.Execute(w, contactDetails)
	errorCheck(err)
}
func UpdateForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contactId := vars["id"]
	ID, err := strconv.ParseInt(contactId, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing blue")
		return
	}

	// Get the contact from the database
	contact, _ := models.GetContactById(ID)

	if contact == nil {
		// Handle case where contact doesn't exist
		fmt.Println("Contact not found")
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	tmpl := template.Must(template.ParseFiles(wd + "/../../static/editcontact.html"))
	err = tmpl.Execute(w, contact)
	errorCheck(err)
}

func UpdateContactByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contactId := vars["id"]
	ID, err := strconv.ParseInt(contactId, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing yellow")
		return
	}

	// Get the updated contact details from the form data
	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	email := r.FormValue("email")
	phone := r.FormValue("phone")

	// Get the contact from the database
	contact, _ := models.GetContactById(ID)

	if contact == nil {
		// Handle case where contact doesn't exist
		fmt.Println("Contact not found")
		return
	}

	// Update the contact details
	contact.FirstName = firstName
	contact.LastName = lastName
	contact.Email = email
	contact.Phone = phone

	// Save the updated contact details in the database
	models.UpdateContact(ID, *contact)
}

func DeleteContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contactId := vars["id"]
	ID, err := strconv.ParseInt(contactId, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing")
	}

	// Delete the contact from the database
	models.DeleteContact(ID)

	// Redirect to the contacts page
	// http.Redirect(w, r, "/contacts", http.StatusOK)
}

// DeleteSelectedContacts deletes the selected contacts from the database and redirects to the contacts page.
func DeleteSelectedContacts(w http.ResponseWriter, r *http.Request) {
	// Get the contact IDs from the request
	e := r.ParseForm()
	if e != nil {
		fmt.Println("Error parsing form: ", e)
		return
	}
	contactIds := r.Form["selectedContacts"]
	fmt.Println("Contact IDs: ", contactIds)

	// Convert the contact IDs to an array of integers
	var ids []int
	for _, id := range contactIds {
		i, err := strconv.Atoi(string(id))
		if err != nil {
			fmt.Println("Error while converting to int: ", err)
		}
		ids = append(ids, i)
	}
	fmt.Println("IDs: ", ids)

	// Delete the contacts from the database
	for _, id := range ids {
		models.DeleteContact(int64(id))
	}

	// Redirect to the contacts page
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

// func StartArchiveContacts(w http.ResponseWriter, r *http.Request) {
// 	archiver := GetArchiver()
// 	archiver.Run()
// 	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
// }

// func ResetArchiveContacts(w http.ResponseWriter, r *http.Request) {
// 	archiver := GetArchiver()
// 	archiver.Reset()
// 	http.Redirect(w, r, "/contacts", http.StatusSeeOther)

// }

// func GetArchiveContent(w http.ResponseWriter, r *http.Request) {
// 	manager := GetArchiver()
// 	archiveFile, e := os.Open(manager.ArchiveFile())
// 	if e != nil {
// 		http.Error(w, "Archive Not Found", http.StatusNotFound)
// 		return
// 	}
// 	defer archiveFile.Close()
// 	w.Header().Set("Content-Disposition", "attachment; filename=\"archive.json\"")
// 	w.Header().Set("Content-Type", "application/json")
// 	http.ServeContent(w, r, "archive.json", time.Now(), archiveFile)

// }
