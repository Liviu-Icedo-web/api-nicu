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
	"api-nicu/api/responses"
	"api-nicu/api/utils/formaterror"

	"api-nicu/api/models"

	"github.com/gorilla/mux"
)

func (server *Server) CreateRental(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	rental := models.Rental{}
	err = json.Unmarshal(body, &rental)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	fmt.Println("*** CreateRental body", body)
	fmt.Println("*** CreateRental body car", rental)

	rental.Prepare()
	err = rental.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	fmt.Println("*** CreateRental UID", uid)

	if uid != rental.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)+" Not the same user logged as the one passed to /rental-car/ "))
		return
	}

	rentalCreated, err := rental.SaveRental(server.DB)
	if err != nil {

		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, rentalCreated.ID))
	responses.JSON(w, http.StatusCreated, rentalCreated)
}

func (server *Server) GetRentalUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	rental := models.Rental{}

	uid, err := auth.ExtractTokenID(r)
	if uid != uint32(pid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)+" Not the same user logged as the one passed to /rental-user/ "))
		return
	}

	rentalList, err := rental.FindRentalUserID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, rentalList)
}

func (server *Server) GetRentalCar(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	rental := models.Rental{}
	rentalList, err := rental.FindRentalCarID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, rentalList)
}

func (server *Server) GetRentalOwner(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	rental := models.Rental{}
	rentalList, err := rental.FindUserOwnerRentals(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, rentalList)
}

func (server *Server) UpdateRental(w http.ResponseWriter, r *http.Request) {

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
	rental := models.Rental{}
	err = server.DB.Debug().Model(rental).Where("id = ?", pid).Take(&rental).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Rental not found"))
		return
	}

	// If a user attempt to update a car not belonging to him
	if uid != rental.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)+" User don't own this rental: "+string(rental.ID)))
		return
	}
	// Read the data cared
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {

		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	rentalUpdate := models.Rental{}
	err = json.Unmarshal(body, &rentalUpdate)
	if err != nil {
		log.Fatalln(err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token

	if uid != rentalUpdate.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)+" User don't own this car_id: "+string(rentalUpdate.UserID)))
		return
	}

	rentalUpdate.Prepare()
	err = rentalUpdate.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// carlocationUpdate, err := carlocationUpdate.UpdateACarLocation(server.DB, pid)
	rentUpdate, err := rentalUpdate.UpdateARental(server.DB, pid)

	if err != nil {
		log.Fatal(err)
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, rentUpdate)
}

func (server *Server) DeleteRentalUser(w http.ResponseWriter, r *http.Request) {

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

	// Check if the carlocation exist
	rental := models.Rental{}
	err = server.DB.Debug().Model(rental).Where("id = ?", pid).Take(&rental).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("No Location_id: "+string(pid)+" into database"))
		return
	}

	// Is the authenticated user, the owner of this car?

	if uid != rental.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized User dont owen this rental"))
		return
	}
	_, err = rental.DeleteARental(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusOK, "Rental: "+string(pid)+" delete")
}
