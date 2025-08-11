package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"muligansim.ryanharris.net/internal/models"
)

func (app *application) viewDecks(w http.ResponseWriter, r *http.Request) {
	// data := templateData{IsAuthenticated: app.isAuthenticated(r)}
	// data.Flash = app.sessionManager.PopString(r.Context(), "flash")
	// userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	// decks, err := app.decks.GetUserDecks(userID)

	// if errors.Is(err, models.ErrNoRecord) {
	// 	app.render(w, r, http.StatusOK, "home.html", data)
	// } else if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	// data.Decks = decks

	// app.render(w, r, http.StatusOK, "home.html", data)
	data := templateData{IsAuthenticated: app.isAuthenticated(r)}
	data.Flash = app.sessionManager.PopString(r.Context(), "flash")
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	decks, err := app.decks.GetUserDecks(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data.Decks = decks

	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *application) viewCards(w http.ResponseWriter, r *http.Request) {

	value := r.PathValue("value")

	cards, err := app.cards.GetByName(value)
	if err != nil {
		app.serverError(w, r, err)
	}

	data := templateData{
		Search:          value,
		Cards:           cards,
		IsAuthenticated: app.isAuthenticated(r),
	}

	app.render(w, r, http.StatusOK, "builddeck.html", data)
}

func (app *application) buildDeck(w http.ResponseWriter, r *http.Request) {
	data := templateData{IsAuthenticated: app.isAuthenticated(r)}
	deckID, err := strconv.Atoi(r.PathValue("deckID"))
	if err != nil {
		app.serverError(w, r, err)
	}

	if deckID != -1 {
		deck, err := app.decks.GetDeckWithCards(deckID)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.serverError(w, r, err)
			}
		}
		data.Deck = deck
	} else {
		data.Deck.DeckID = -1
	}

	app.render(w, r, http.StatusOK, "builddeck.html", data)
}

func (app *application) search(w http.ResponseWriter, r *http.Request) {

	name := r.PathValue("name")

	fmt.Print(name)

	cards, err := app.cards.GetByName(name)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(cards)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	/*http.Redirect(w, r, fmt.Sprintf("/cards/view/%s", name), http.StatusSeeOther)*/

}

type saveID struct {
	Id int `json:"id"`
}

var deckNameReg = regexp.MustCompile(`^[A-Za-z0-9'\s]+$`)

func (app *application) saveDeck(w http.ResponseWriter, r *http.Request) {
	var save models.SaveDeck
	err := decodeJSONBody(w, r, &save)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if !deckNameReg.MatchString(save.Name) && save.Name != "" {
		log.Print("Bad Input")
		http.Error(w, "Invalid Deck Name", http.StatusBadRequest)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	deck := saveID{}
	deck.Id, err = app.decks.SaveDeckChanges(r.Context(), save, id)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}

	data, err := json.Marshal(deck)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

type signupForm struct {
	Name   string
	Email  string
	Pass   string
	Errors map[string]string
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	data := templateData{IsAuthenticated: app.isAuthenticated(r)}
	data.Form = signupForm{}
	app.render(w, r, http.StatusOK, "signup.html", data)

}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, r, err)
	}

	name := r.PostForm.Get("uname")
	email := r.PostForm.Get("email")
	pass := r.PostForm.Get("pass")

	form := signupForm{
		Name:   name,
		Email:  email,
		Pass:   pass,
		Errors: map[string]string{},
	}

	var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !EmailRX.MatchString(email) {
		form.Errors["email"] = "Invalid Email"
	}

	if len(name) == 0 {
		form.Errors["name"] = "Invalid Name"
	}

	if len(pass) < 8 {
		form.Errors["pass"] = "Invalid Password"
	}

	if len(form.Errors) == 0 {
		err = app.users.Insert(name, email, pass)
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors["email"] = "Email already in use"
		}
	}

	if len(form.Errors) > 0 {
		data := templateData{}
		data.Form = form
		app.render(w, r, http.StatusSeeOther, "signup.html", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "You successfully signed up!")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

type loginForm struct {
	Errors map[string]string
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	form := loginForm{
		Errors: map[string]string{},
	}
	data := templateData{Form: form, IsAuthenticated: app.isAuthenticated(r)}
	data.Flash = app.sessionManager.PopString(r.Context(), "flash")
	app.render(w, r, http.StatusOK, "login.html", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		app.serverError(w, r, err)
	}

	email := r.PostForm.Get("email")
	pass := r.PostForm.Get("pass")

	if len(email) == 0 || len(pass) == 0 {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	id, err := app.users.Authenticate(email, pass)

	if err != nil {

		if errors.Is(err, models.ErrInvalidCredentials) {
			form := loginForm{
				Errors: map[string]string{},
			}
			data := templateData{}
			form.Errors["credentials"] = "Invalid username or password"
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.sessionManager.Put(r.Context(), "flash", "You successfully logged in!")

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}
