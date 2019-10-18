package main

import (
	"log"
	"net/http"

	"./email"
	"./urlVerification"
)

func main() {
	http.HandleFunc("/login", email.SendMail)
	http.HandleFunc("/auth", urlVerification.VerifyToken)
	http.HandleFunc("/verified", urlVerification.Verified)

	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Println(err)
	}
}
