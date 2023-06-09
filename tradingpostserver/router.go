
package tradingpostserver

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/gorilla/mux"
	"github.com/royvandewater/trading-post/auth0creds"
	"github.com/royvandewater/trading-post/buyorderscontroller"
	"github.com/royvandewater/trading-post/ordersservice"
	"github.com/royvandewater/trading-post/profilescontroller"
	"github.com/royvandewater/trading-post/sellorderscontroller"
	"github.com/royvandewater/trading-post/userscontroller"
	"github.com/royvandewater/trading-post/usersservice"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

func healthcheck(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("{\"online\": true}"))
}

func newRouter(auth0Creds auth0creds.Auth0Creds, mongoDB *mgo.Session) http.Handler {
	usersService := usersservice.New(auth0Creds, mongoDB)
	ordersService := ordersservice.New(mongoDB, usersService)

	buyOrdersController := buyorderscontroller.New(ordersService)
	profilesController := profilescontroller.New(usersService)
	sellOrdersController := sellorderscontroller.New(ordersService)
	usersController := userscontroller.New(usersService)

	profileRouter := mux.NewRouter().PathPrefix("/profile").Subrouter()
	profileRouter.Methods("GET").Path("/").HandlerFunc(profilesController.Get)

	profileRouter.Methods("POST").Path("/buy-orders").HandlerFunc(buyOrdersController.Create)
	profileRouter.Methods("GET").Path("/buy-orders").HandlerFunc(buyOrdersController.List)
	profileRouter.Methods("GET").Path("/buy-orders/{id}").HandlerFunc(buyOrdersController.Get)

	profileRouter.Methods("POST").Path("/sell-orders").HandlerFunc(sellOrdersController.Create)
	profileRouter.Methods("GET").Path("/sell-orders").HandlerFunc(sellOrdersController.List)
	profileRouter.Methods("GET").Path("/sell-orders/{id}").HandlerFunc(sellOrdersController.Get)

	router := mux.NewRouter()
	router.Methods("GET").Path("/profiles").HandlerFunc(usersController.List)
	router.Methods("GET").Path("/callback").HandlerFunc(usersController.Login)
	router.Methods("POST").Path("/tokens").HandlerFunc(usersController.Token)
	router.PathPrefix("/profile").Handler(negroni.New(
		negroni.HandlerFunc(usersController.Authenticate),
		negroni.Wrap(profileRouter),
	))
	router.Methods("GET").Path("/healthcheck").HandlerFunc(healthcheck)
	router.Methods("GET").Handler(http.FileServer(http.Dir("html/")))

	allowCORS := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"DELETE", "GET", "PATCH", "POST", "PUT"},
		AllowCredentials: true,
	})
	return negroni.New(allowCORS, negroni.Wrap(router))
}