package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/sntruong12/josie-rides/internal/models"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Josie Rides" as the response body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// handles edge case when users nav to non existing routes
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
	id, err := strconv.Atoi(r.PathValue("id"))
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

func (app *application) rideCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = rideCreateForm{}

	app.render(w, http.StatusOK, "create.html", data)
}

type rideCreateForm struct {
	Title           string
	Description     string
	TrailName       string
	Distance        string
	DurationHours   string
	DurationMinutes string
	DurationSeconds string
	RodeAt          string
	Timezone        string
	Media           string

	FieldErrors map[string]string
}

// Add a rideCreate handler function.
func (app *application) rideCreatePost(w http.ResponseWriter, r *http.Request) {
	// Limit the request body size to 4096 bytes
	// r.Body = http.MaxBytesReader(w, r.Body, 4096)

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := rideCreateForm{
		Title:           strings.TrimSpace(r.PostForm.Get("title")),
		Description:     strings.TrimSpace(r.PostForm.Get("description")),
		Distance:        r.PostForm.Get("distance_miles"),
		DurationHours:   r.PostForm.Get("duration_hours"),
		DurationMinutes: r.PostForm.Get("duration_minutes"),
		DurationSeconds: r.PostForm.Get("duration_seconds"),
		RodeAt:          r.PostForm.Get("rode_at"),
		Timezone:        r.PostForm.Get("timezone"),
		TrailName:       strings.TrimSpace(r.PostForm.Get("trail_name")),

		FieldErrors: make(map[string]string),
	}

	// need to implement this later
	media := json.RawMessage("null")

	// validation
	if form.Title == "" {
		form.FieldErrors["title"] = "Title is required"
	}

	if utf8.RuneCountInString(form.Title) > 255 {
		form.FieldErrors["title"] = "This field must be less than 255 characters"
	}

	if form.Description == "" {
		form.FieldErrors["description"] = "Description is required"
	}

	if form.TrailName == "" {
		form.FieldErrors["trail_name"] = "Trail name is required"
	}

	if utf8.RuneCountInString(form.TrailName) > 255 {
		form.FieldErrors["trail_name"] = "This field must be less than 255 characters"
	}

	if form.Distance == "" {
		form.FieldErrors["distance"] = "Distance is required"
	}

	if utf8.RuneCountInString(form.Distance) > 8 {
		form.FieldErrors["distance"] = "Distance is too long"
	}

	if form.DurationHours == "" {
		form.FieldErrors["duration_hours"] = "Duration hours is required"
	}

	if form.DurationMinutes == "" {
		form.FieldErrors["duration_minutes"] = "Duration minutes is required"
	}

	if form.DurationSeconds == "" {
		form.FieldErrors["duration_seconds"] = "Duration seconds is required"
	}

	if form.RodeAt == "" {
		form.FieldErrors["rode_at"] = "Rode at is required"
	}

	if form.Timezone == "" {
		form.FieldErrors["timezone"] = "Timezone is required"

	}

	// convert distance to float64
	distanceFloat, err := strconv.ParseFloat(form.Distance, 64)
	if err != nil {
		form.FieldErrors["distance"] = "Distance is invalid"
	}

	if distanceFloat <= 0 || distanceFloat > 99999.99 {
		form.FieldErrors["distance"] = "Distance is invalid"
	}

	// convert rodeAt to time.Time utc
	loc, err := time.LoadLocation(form.Timezone)
	if err != nil {
		form.FieldErrors["timezone"] = "Timezone is invalid"
	}

	rodeAtTime, err := time.ParseInLocation("2006-01-02T15:04", form.RodeAt, loc)
	if err != nil {
		form.FieldErrors["rode_at"] = "Rode at is invalid"
	}
	rodeAtUTC := rodeAtTime.UTC()

	// convert duration into seconds
	h, err := strconv.Atoi(form.DurationHours)
	if err != nil {
		form.FieldErrors["duration_hours"] = "Duration hours is invalid"
	}

	if h < 0 || h > 99 {
		form.FieldErrors["duration_hours"] = "Duration hours is invalid"
	}

	m, err := strconv.Atoi(form.DurationMinutes)
	if err != nil {
		form.FieldErrors["duration_minutes"] = "Duration minutes is invalid"
	}

	if m < 0 || m > 59 {
		form.FieldErrors["duration_minutes"] = "Duration minutes is invalid"
	}

	s, err := strconv.Atoi(form.DurationSeconds)
	if err != nil {
		form.FieldErrors["duration_seconds"] = "Duration seconds is invalid"
	}

	if s < 0 || s > 59 {
		form.FieldErrors["duration_seconds"] = "Duration seconds is invalid"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	convertedDuration := (h * 3600) + (m * 60) + s

	id, err := app.rides.Create(form.Title, form.Description, form.TrailName, distanceFloat, convertedDuration, rodeAtUTC, media)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/ride/view/%d", id), http.StatusSeeOther)
}
