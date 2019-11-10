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

func (server *Server) CreateCarLocation(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	fmt.Println("*** Car Controller body", body)
	carlocation := models.CarLocation{}
	err = json.Unmarshal(body, &carlocation)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	fmt.Println("*** Car Controller body", body)
	fmt.Println("*** Car Controller body car", carlocation)

	carlocation.Prepare()
	err = carlocation.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	car := models.Car{}
	err = server.DB.Debug().Model(models.Car{}).Where("id = ?", carlocation.CarID).Take(&car).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("User don't own this car_id: "+string(carlocation.CarID)))
		return
	}

	if uid != car.User_id {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)+" User don't own this car_id: "+string(carlocation.CarID)))
		return
	}
	fmt.Println("*** Car Controller", car)
	carCreated, err := carlocation.SaveCarLocation(server.DB)
	if err != nil {

		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, carCreated.ID))
	responses.JSON(w, http.StatusCreated, carCreated)
}

func (server *Server) GetCarsLocation(w http.ResponseWriter, r *http.Request) {

	car := models.CarLocation{}

	cars, err := car.FindAllCarsLocation(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, cars)
}

func (server *Server) GetCarLocation(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	car := models.CarLocation{}

	carReceived, err := car.FindCarLocationByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, carReceived)
}

func (server *Server) UpdateCarLocation(w http.ResponseWriter, r *http.Request) {

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
	carlocation := models.CarLocation{}
	err = server.DB.Debug().Model(models.CarLocation{}).Where("id = ?", pid).Take(&carlocation).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Car not found"))
		return
	}

	car := models.Car{}
	err = server.DB.Debug().Model(models.Car{}).Where("id = ?", carlocation.CarID).Take(&car).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Database error: User don't own this car_id: "+string(carlocation.CarID)))
		return
	}
	// If a user attempt to update a car not belonging to him
	if uid != car.User_id {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)+" User don't own this car_id: "+string(carlocation.CarID)))
		return
	}
	// Read the data cared
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {

		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	carlocationUpdate := models.CarLocation{}
	err = json.Unmarshal(body, &carlocationUpdate)
	if err != nil {
		log.Fatalln(err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	fmt.Printf("\n **** carUpdate %v", carlocationUpdate)

	carUpdate := models.Car{}
	err = server.DB.Debug().Model(models.Car{}).Where("id = ?", carlocation.CarID).Take(&carUpdate).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("User don't own this car_id: "+string(carlocation.CarID)))
		return
	}

	fmt.Println("*** Updatecar UID ", uid)
	fmt.Println("*** Updatecar car to update ", carUpdate)
	if uid != carUpdate.User_id {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)+" User don't own this car_id: "+string(carlocationUpdate.CarID)))
		return
	}

	carlocationUpdate.Prepare()
	err = carlocationUpdate.Validate()
	fmt.Printf("Prepare  \n")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// carlocationUpdate, err := carlocationUpdate.UpdateACarLocation(server.DB, pid)
	carlocUpdate, err := carlocationUpdate.UpdateACarLocation(server.DB, pid)

	if err != nil {
		log.Fatal(err)
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, carlocUpdate)
}

func (server *Server) DeleteCarLocation(w http.ResponseWriter, r *http.Request) {

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
	carLocation := models.CarLocation{}
	err = server.DB.Debug().Model(models.CarLocation{}).Where("id = ?", pid).Take(&carLocation).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("No Location_id: "+string(pid)+" into database"))
		return
	}

	// Is the authenticated user, the owner of this car?
	car := models.Car{}
	err = server.DB.Debug().Model(models.Car{}).Where("id = ?", carLocation.CarID).Take(&car).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("We dont find a car_location for this car_id: "+string(carLocation.CarID)))
		return
	}

	if uid != car.User_id {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized User dont owen this car"))
		return
	}
	_, err = carLocation.DeleteACarLocation(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusOK, "Location: "+string(pid)+" delete")
}
