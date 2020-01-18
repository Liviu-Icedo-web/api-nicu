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

	//Cars routes
	s.Router.HandleFunc("/cars", middlewares.SetMiddlewareJSON(s.CreateCar)).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/cars", middlewares.SetMiddlewareJSON(s.GetCars)).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/cars/{id}", middlewares.SetMiddlewareJSON(s.GetCar)).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/cars-user/{id}", middlewares.SetMiddlewareJSON(s.GetCarbyUser)).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/cars/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateCar))).Methods("PUT", "OPTIONS")
	s.Router.HandleFunc("/cars/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteCar)).Methods("DELETE", "OPTIONS")

	//Cars location routes
	s.Router.HandleFunc("/car-location/", middlewares.SetMiddlewareJSON(s.CreateCarLocation)).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/cars-location/", middlewares.SetMiddlewareJSON(s.GetCarsLocation)).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/cars-location/{id}", middlewares.SetMiddlewareJSON(s.GetCarLocation)).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/car-location/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateCarLocation))).Methods("PUT", "OPTIONS")
	s.Router.HandleFunc("/car-location/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteCarLocation)).Methods("DELETE", "OPTIONS")

	//Rental routes
	s.Router.HandleFunc("/rental-car/", middlewares.SetMiddlewareJSON(s.CreateRental)).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/rental-car/{id}", middlewares.SetMiddlewareJSON(s.GetRentalCar)).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/rental-owners/{id}", middlewares.SetMiddlewareJSON(s.GetRentalOwner)).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/rental-user/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetRentalUser))).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/rental-user/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateRental))).Methods("PUT", "OPTIONS")
	s.Router.HandleFunc("/rental-user/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteRentalUser)).Methods("DELETE", "OPTIONS")
}
