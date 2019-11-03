package controllers

import "api-nicu/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	//s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/user/{id}", middlewares.SetMiddlewareJSON(s.GetUserId)).Methods("GET")
	s.Router.HandleFunc("/user", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUser))).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT", "OPTIONS")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE", "OPTIONS")

	//Posts routes
	s.Router.HandleFunc("/cars", middlewares.SetMiddlewareJSON(s.CreateCar)).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/cars", middlewares.SetMiddlewareJSON(s.GetCars)).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/cars/{id}", middlewares.SetMiddlewareJSON(s.GetCar)).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/cars/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateCar))).Methods("PUT", "OPTIONS")
	s.Router.HandleFunc("/cars/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteCar)).Methods("DELETE", "OPTIONS")
}
