package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"api-nicu/api/auth"
	"api-nicu/api/models"
	"api-nicu/api/responses"
	"api-nicu/api/utils/formaterror"

	"github.com/gorilla/mux"
)

func (server *Server) CreateCarLocation(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	car := models.CarLocation{}
	err = json.Unmarshal(body, &car)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	car.Prepare()
	err = car.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	car.CarID = uid

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != car.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	carCreated, err := car.SaveCarLocation(server.DB)
	if err != nil {

		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, carCreated.ID))
	responses.JSON(w, http.StatusCreated, carCreated)
}

func (server *Server) GetCars(w http.ResponseWriter, r *http.Request) {

	car := models.Car{}

	cars, err := car.FindAllCars(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, cars)
}

func (server *Server) GetCar(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	car := models.Car{}

	carReceived, err := car.FindCarByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, carReceived)
}

func (server *Server) UpdateCar(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the car id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the car exist
	car := models.Car{}
	err = server.DB.Debug().Model(models.Car{}).Where("id = ?", pid).Take(&car).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Car not found"))
		return
	}

	fmt.Printf("User id: %v", uid)
	fmt.Printf("\n Car id: %v", car.ID)
	fmt.Printf("\n First check: %v", car.User_id)

	// If a user attempt to update a car not belonging to him
	if uid != car.User_id {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data cared
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {

		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	carUpdate := models.Car{}
	err = json.Unmarshal(body, &carUpdate)
	if err != nil {
		log.Fatalln(err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token

	fmt.Printf("\n User id %v", uid)
	fmt.Printf("\n carUpdate.User_id %v", carUpdate.User_id)
	fmt.Printf("\n **** carUpdate %v", carUpdate)
	if uid != carUpdate.User_id {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	carUpdate.Prepare()
	err = carUpdate.Validate()
	fmt.Printf("Prepare  \n")
	if err != nil {
		fmt.Printf("Prepare  *****\n")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	carUpdated, err := carUpdate.UpdateACar(server.DB, pid)

	if err != nil {
		log.Fatal(err)
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, carUpdated)
}

func (server *Server) DeleteCar(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid car id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the car exist
	car := models.Car{}
	err = server.DB.Debug().Model(models.Car{}).Where("id = ?", pid).Take(&car).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this car?
	if uid != car.User_id {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = car.DeleteACar(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
