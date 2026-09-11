package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sntruong12/josie-rides/internal/models"
	"github.com/sntruong12/josie-rides/internal/validator"
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
	Title           string `form:"title"`
	Description     string `form:"description"`
	TrailName       string `form:"trail_name"`
	Distance        string `form:"distance"`
	DurationHours   string `form:"duration_hours"`
	DurationMinutes string `form:"duration_minutes"`
	DurationSeconds string `form:"duration_seconds"`
	RodeAt          string `form:"rode_at"`
	Timezone        string `form:"timezone"`
	Media           string `form:"media"`

	validator.Validator `form:"-"`
}

// Add a rideCreate handler function.
func (app *application) rideCreatePost(w http.ResponseWriter, r *http.Request) {
	// Limit the request body size to 4096 bytes
	// r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var form rideCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// form.Title = strings.TrimSpace(form.Title)
	// form.Description = strings.TrimSpace(form.Description)
	// form.TrailName = strings.TrimSpace(form.TrailName)

	// need to implement this later
	media := json.RawMessage("null")

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 255), "title", "This field must be less than 255 characters")

	form.CheckField(validator.NotBlank(form.Description), "description", "This field cannot be blank")

	form.CheckField(validator.NotBlank(form.TrailName), "trail_name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.TrailName, 255), "trail_name", "This field must be less than 255 characters")

	form.CheckField(validator.NotBlank(form.Distance), "distance", "Distance is required")
	form.CheckField(validator.MaxChars(form.Distance, 8), "distance", "Distance is too long")
	form.CheckField(validator.ValidFloat(form.Distance, 0.01, 99999.99), "distance", "Distance is invalid")

	form.CheckField(validator.NotBlank(form.DurationHours), "duration_hours", "Duration hours is required")
	form.CheckField(validator.MaxChars(form.DurationHours, 2), "duration_hours", "Duration hours is too long")
	form.CheckField(validator.ValidInt(form.DurationHours, 0, 99), "duration_hours", "Duration hours is invalid")

	form.CheckField(validator.NotBlank(form.DurationMinutes), "duration_minutes", "Duration minutes is required")
	form.CheckField(validator.MaxChars(form.DurationMinutes, 2), "duration_minutes", "Duration minutes is too long")
	form.CheckField(validator.ValidInt(form.DurationMinutes, 0, 59), "duration_minutes", "Duration minutes is invalid")

	form.CheckField(validator.NotBlank(form.DurationSeconds), "duration_seconds", "Duration seconds is required")
	form.CheckField(validator.MaxChars(form.DurationSeconds, 2), "duration_seconds", "Duration seconds is too long")
	form.CheckField(validator.ValidInt(form.DurationSeconds, 0, 59), "duration_seconds", "Duration seconds is invalid")

	form.CheckField(validator.NotBlank(form.RodeAt), "rode_at", "Rode at is required")
	form.CheckField(validator.ValidTimeAndTimezone(form.RodeAt, form.Timezone), "rode_at", "Rode at is invalid")

	form.CheckField(validator.NotBlank(form.Timezone), "timezone", "Timezone is required")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	// convert rodeAt to time.Time utc
	loc, _ := time.LoadLocation(form.Timezone)
	rodeAtTime, _ := time.ParseInLocation("2006-01-02T15:04", form.RodeAt, loc)
	rodeAtUTC := rodeAtTime.UTC()
	distanceFloat, _ := strconv.ParseFloat(form.Distance, 64)
	h, _ := strconv.Atoi(form.DurationHours)
	m, _ := strconv.Atoi(form.DurationMinutes)
	s, _ := strconv.Atoi(form.DurationSeconds)
	convertedDuration := (h * 3600) + (m * 60) + s

	id, err := app.rides.Create(form.Title, form.Description, form.TrailName, distanceFloat, convertedDuration, rodeAtUTC, media)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the Put() method to add a string value ("Ride successfully
	// created!") and the corresponding key ("flash") to the session data.
	app.sessionManager.Put(r.Context(), "flash", "Ride successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/ride/view/%d", id), http.StatusSeeOther)
}

// Create a new userSignupForm struct.
type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}

	app.render(w, http.StatusOK, "signup.html", data)
}
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// form.Email = strings.TrimSpace(form.Email)
	// form.Name = strings.TrimSpace(form.Name)
	// form.Password = strings.TrimSpace(form.Password)

	// Validate the form contents using our helper functions.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = app.users.Create(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked.
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// Create a new userLoginForm struct.
type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	// Decode the form data into the userLoginForm struct.
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Do some validation checks on the form. We check that both email and
	// password are provided, and also check the format of the email address as
	// a UX-nicety (in case the user makes a typo).
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Use the RenewToken() method on the current session to change the session
	// ID. It's good practice to generate a new session ID when the
	// authentication state or privilege levels changes for the user (e.g. login
	// and logout operations).
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Add the ID of the current user to the session, so that they are now
	// 'logged in'.
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	// Redirect the user to the create ride page.
	http.Redirect(w, r, "/ride/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Use the RenewToken() method on the current session to change the session
	// ID again.
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	// Add a flash message to the session to confirm to the user that they've been
	// logged out.
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	// Redirect the user to the application home page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
