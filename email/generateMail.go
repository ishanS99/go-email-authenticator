package email

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"../urlGeneration"
)

// SendMail creates a One Time Password Token for login verification
// of the user. For now sender and receiver are set with host gmail.com
func SendMail(w http.ResponseWriter, r *http.Request) {

	mail := createMail(r)

	// add authentication for the sender
	auth := smtp.PlainAuth(
		"",
		os.Getenv("FROM_USERNAME"),
		os.Getenv("FROM_PASSWORD"),
		"smtp.gmail.com",
	)

	// send the email to the logging user
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		os.Getenv("FROM_USERNAME")+"@gmail.com",
		[]string{os.Getenv("TO_USERNAME") + "@gmail.com"},
		[]byte(mail),
	)

	if err != nil {
		log.Println(err)
	} else {
		_, _ = fmt.Fprint(w, "Check your email for Verification")
	}
}

func createMail(r *http.Request) string {
	q := r.URL.Query()
	username := q.Get("username")

	url := urlGeneration.GenerateUrl(username)

	msg := "Subject: Go - Verification Mail\r\n" +
		"\r\nClick here for verification : "
	msg += url
	msg += " (Valid for 5 Minutes)"

	return msg
}
