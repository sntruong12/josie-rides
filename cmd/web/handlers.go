package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

	data := app.newTemplateData(r)
	data.Rides = rides

	// Use the new render helper.
	app.render(w, http.StatusOK, "home.html", data)
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

	data := app.newTemplateData(r)
	data.Ride = ride

	app.render(w, http.StatusOK, "view.html", data)
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
