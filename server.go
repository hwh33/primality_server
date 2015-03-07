/* This is the HTTP server clients interact with. */

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/hwh33/primality_server/primality"
	"github.com/hwh33/primality_server/registrar"
	// "./primality"
	// "./registrar"
)

const address = ":80"
const debug = true

var serverRegistrar *registrar.Registrar

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if debug {
		fmt.Println("Root Handler; URL = " + r.URL.Path)
	}

	http.ServeFile(w, r, "html/entrypoint.html")
}

func entryPointHandler(w http.ResponseWriter, r *http.Request) {
	if debug {
		fmt.Println("Entry Point Handler; URL = " + r.URL.Path)
		fmt.Println("User type = " + r.FormValue("user type"))
	}

	userType := r.FormValue("user type")
	if userType == "new" {
		http.ServeFile(w, r, "html/new_user.html")
	} else if userType == "returning" {
		http.ServeFile(w, r, "html/login.html")
	} else {
		http.ServeFile(w, r, "html/entrypoint.html")
		// fmt.Println("Error: undefined user type: [" + userType + "]")
		// http.NotFound(w, r)
	}
}

func newUserHandler(w http.ResponseWriter, r *http.Request) {
	if debug {
		fmt.Println("New User Handler; URL = " + r.URL.Path)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	netID := r.FormValue("netID")

	if username == "" || password == "" || netID == "" {
		displayError(w, r, "Registration Error: you must enter values for "+
			"your username, net ID, and password")
		return
	}

	err := serverRegistrar.RegisterUser(username, netID, password)
	if err != nil {
		fmt.Println("Error while registering user: " + err.Error())
		displayError(w, r, "Registration error: "+err.Error())
		return
	}

	http.ServeFile(w, r, "html/primality_test.html")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if debug {
		fmt.Println("Login Handler; URL = " + r.URL.Path)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if serverRegistrar.IsPasswordAuthentic(username, password) {
		http.ServeFile(w, r, "html/primality_test.html")
	} else {
		http.ServeFile(w, r, "html/login_try_again.html")
	}
}

func primalityTestHandler(w http.ResponseWriter, r *http.Request) {
	if debug {
		fmt.Println("Primality Test Handler; URL = " + r.URL.Path)
	}

	inputString := r.FormValue("input")
	if inputString == "" {
		fmt.Println("No input detected in primality test.")
		displayError(w, r, "No input detected")
		return
	}
	base, bitsize := 10, 64
	input, err := strconv.ParseUint(inputString, base, bitsize)
	if err != nil {
		fmt.Println("Error processing user input: " + err.Error())
		displayError(w, r, "Input must be an integer")
		return
	}

	isPrime := primality.IsPrime(input)
	if isPrime {
		fmt.Fprintf(w, inputString+" is prime.")
	} else {
		fmt.Fprintf(w, inputString+" is not prime.")
	}
}

func displayError(w http.ResponseWriter, r *http.Request, errorMsg string) {
	templ, err := template.ParseFiles("html/error_template.html")
	if err != nil {
		fmt.Println("Error generating html template: " + err.Error())
		http.ServeFile(w, r, "html/backup_error.html")
	}
	templ.Execute(w, "Whoops! An error occurred. "+errorMsg+
		". Use your browser to return to the previous page.")
}

func main() {

	fmt.Println("Creating registrar file")
	var err error
	serverRegistrar, err = registrar.NewRegistrar("registrar_file.csv")
	if err != nil {
		panic("Error in creating registrar: " + err.Error())
	}

	fmt.Println("Starting up at " + address)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/entrypoint", entryPointHandler)
	http.HandleFunc("/new_user", newUserHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/primality_test", primalityTestHandler)

	http.ListenAndServe(address, nil)

}
