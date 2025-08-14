package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)
	actBefore := alice.New(commonHeaders, app.logRequest)
	authenticated := dynamic.Append(app.pathAuth)
	deckAuth := dynamic.Append(app.deckBuildAuth)

	mux.Handle("GET /home", authenticated.ThenFunc(app.viewDecks))
	mux.Handle("GET /cards/view/{value}", authenticated.ThenFunc(app.viewCards))
	mux.Handle("GET /cards/search/{name}", authenticated.ThenFunc(app.search))
	mux.Handle("PUT /cards/save", authenticated.ThenFunc(app.saveDeck))
	mux.Handle("GET /cards/builddeck/{deckID}", deckAuth.ThenFunc(app.buildDeck))
	mux.Handle("DELETE /cards/deleteDeck/{deckID}", deckAuth.ThenFunc(app.deleteDeck))

	//Routes for user auth
	mux.Handle("GET /users/signup", dynamic.ThenFunc(app.signup))
	mux.Handle("POST /users/signup", dynamic.ThenFunc(app.signupPost))
	mux.Handle("GET /users/login", dynamic.ThenFunc(app.login))
	mux.Handle("POST /users/login", dynamic.ThenFunc(app.loginPost))
	mux.Handle("POST /users/logout", dynamic.ThenFunc(app.logout))
	return actBefore.Then(mux)
}
