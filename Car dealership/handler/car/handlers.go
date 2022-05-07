package car

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/models"
	"github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/service"

	"github.com/gorilla/mux"
)

type handler struct {
	service service.Cars
}

func New(c service.Cars) handler { //nolint
	return handler{service: c}
}

// GetCarByID handler layer function to get car record by giving car id
func (c handler) GetCarByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := c.service.GetCarByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)

		return
	}

	body, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetCarByBrand handler layer function to get car record by giving car brand
func (c handler) GetCarByBrand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	brand := r.URL.Query().Get("brand")
	isEngine := r.URL.Query().Get("isEngine")

	isEng, err := strconv.ParseBool(isEngine)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := c.service.GetCarByBrand(ctx, brand, isEng)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("DB error"))
	}

	body, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")

	_, _ = w.Write(body)
	w.WriteHeader(http.StatusOK)
}

// CreateCar handler layer function to create car record
func (c handler) CreateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	var car *models.Car

	err = json.Unmarshal(body, &car)
	if err != nil {
		//fmt.Println("debug")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := c.service.CreateCar(ctx, car)
	if err != nil {
		// to print error on postman

		//b, er := json.Marshal(err.Error())
		//if er != nil {
		//	w.WriteHeader(http.StatusBadRequest)
		//}
		w.WriteHeader(http.StatusBadRequest)
		//_, _ = w.Write(b)
		return
	}

	body, err = json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	w.Header().Set("Content-Type", "application/json")

	_, _ = w.Write(body)
}

// UpdateCar handler layer function to update car record
func (c handler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, _ := io.ReadAll(r.Body)

	var car models.Car

	err := json.Unmarshal(body, &car)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	param := mux.Vars(r)
	id := param["id"]

	res, err := c.service.UpdateCar(ctx, id, car)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	body, err = json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, _ = w.Write(body)
}

// DeleteCar handler layer function to delete car record
func (c handler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	param := mux.Vars(r)
	id := param["id"]

	_, err := c.service.DeleteCar(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		_, _ = w.Write([]byte("invalid id"))

		return
	}

	_, _ = w.Write([]byte("deleted"))
}
