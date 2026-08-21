package main

import (
	"encoding/json"
	"errors"
	"fmt"

	// "html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/sntruong12/josie-rides/internal/models"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Josie Rides" as the response body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't, use
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the handler
	// would keep executing.
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	rides, err := app.rides.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, ride := range rides {
		fmt.Fprintf(w, "%+v\n", ride)
	}

	// Initialize a slice containing the paths to html files. It's important
	// to note that the file containing our base template must be the *first*
	// file in the slice.
	// files := []string{
	// 	"./ui/html/base.html",
	// 	"./ui/html/pages/home.html",
	// 	"./ui/html/partials/nav.html",
	// }

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and use
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// Use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body.
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, err)
	// }
}

// Add a rideView handler function.
func (app *application) rideView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	ride, err := app.rides.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", ride)
}

// Add a rideCreate handler function.
func (app *application) rideCreate(w http.ResponseWriter, r *http.Request) {
	// don't allow any method other than POST
	if r.Method != http.MethodPost {
		// Use the Header().Set() method to add an 'Allow: POST' header to the
		// response header map. The first parameter is the header name, and
		// the second parameter is the header value.
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "dummy title"
	description := "dummy description"
	trailName := "dummy trail name"
	distance := 10.0
	duration := 9000
	rodeAt := time.Now()
	media := json.RawMessage("null")

	id, err := app.rides.Create(title, description, trailName, distance, duration, rodeAt, media)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/ride/view?id=%d", id), http.StatusSeeOther)
}
