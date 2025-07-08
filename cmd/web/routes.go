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

	mux.Handle("GET /decks/home", dynamic.ThenFunc(app.viewDecks))
	mux.Handle("GET /cards/view/{value}", dynamic.ThenFunc(app.viewCards))
	mux.Handle("POST /cards/search", dynamic.ThenFunc(app.searchCardsPost))
	mux.Handle("GET /cards/builddeck", dynamic.ThenFunc(app.buildDeck))

	//Routes for user auth
	mux.Handle("GET /users/signup", dynamic.ThenFunc(app.signup))
	mux.Handle("POST /users/signup", dynamic.ThenFunc(app.signupPost))
	mux.Handle("GET /users/login", dynamic.ThenFunc(app.login))
	mux.Handle("POST /users/login", dynamic.ThenFunc(app.loginPost))
	return actBefore.Then(mux)
}
