package controllers

import "api-nicu/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	//s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Posts routes
	s.Router.HandleFunc("/cars", middlewares.SetMiddlewareJSON(s.CreateCar)).Methods("POST")
	s.Router.HandleFunc("/cars", middlewares.SetMiddlewareJSON(s.GetCars)).Methods("GET")
	s.Router.HandleFunc("/cars/{id}", middlewares.SetMiddlewareJSON(s.GetCar)).Methods("GET")
	s.Router.HandleFunc("/cars/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateCar))).Methods("PUT")
	s.Router.HandleFunc("/cars/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteCar)).Methods("DELETE")
}
