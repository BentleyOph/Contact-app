package utils

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"archive/zip"
	"os"
	"github.com/BentleyOph/htmx-go/pkg/models"
	"log"
)

type Archiver struct {
	mu              sync.Mutex
	archiveStatus   string
	archiveProgress float64
}

var (
	instance *Archiver
	once     sync.Once
)

func GetArchiver() *Archiver {
	once.Do(func() {
		instance = &Archiver{archiveStatus: "Waiting"}
	})
	return instance
}

func (a *Archiver) Status() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.archiveStatus
}

func (a *Archiver) Progress() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.archiveProgress
}

func (a *Archiver) ProgressPercent() int {
	return int(a.Progress() * 100)
}

func (a *Archiver) Run() {
	a.mu.Lock()
	if a.archiveStatus == "Waiting" {
		a.archiveStatus = "Running"
		a.archiveProgress = 0
		go a.runImpl()
	}
	a.mu.Unlock()
}

func (a *Archiver) runImpl() {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(rand.Float64() * float64(time.Second)))
		a.mu.Lock()
		if a.archiveStatus != "Running" {
			a.mu.Unlock()
			return
		}
		a.archiveProgress = float64(i+1) / 10
		fmt.Printf("Here... %.2f\n", a.archiveProgress)
		a.mu.Unlock()
	}
	time.Sleep(time.Second)
	a.mu.Lock()
	if a.archiveStatus == "Running" {
		a.archiveStatus = "Complete"
	}
	a.mu.Unlock()
}


func (a *Archiver) ArchiveFile() string {
	//use models.DownloadContants to generate the archive file
	contacts := models.DownloadContacts()

	//create a file with the contacts
	file, err := os.Create("contacts.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//create a new zip archive
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	//iterate over the contacts and add them to the zip file
	for _, contact := range contacts {
		contactFile, err := zipWriter.Create(contact.FirstName + ".txt")
		if err != nil {
			log.Fatal(err)
		}

		_, err = contactFile.Write([]byte(contact.LastName + "\n" + contact.Email + "\n" + contact.Phone))
		if err != nil {
			log.Fatal(err)
		}
	}

	return "contacts.zip"
}

func (a *Archiver) Reset() {
	a.mu.Lock()
	a.archiveStatus = "Waiting"
	a.mu.Unlock()
}
