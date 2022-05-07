package main

import (
	"github.com/gorilla/mux"
	store "github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/datastore/car"
	"github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/datastore/engine"
	"github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/driver"
	handler "github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/handler/car"
	"github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/middleware"
	service "github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/service/car"
	"log"
	"net/http"
)

func main() {

	db := driver.ConnectToSQL()
	st := store.New(db)

	engin := engine.New(db)
	svc := service.New(st, engin)
	list := handler.New(svc)

	r := mux.NewRouter()

	r.HandleFunc("/car/{id}", list.GetCarByID).Methods(http.MethodGet)
	r.HandleFunc("/cars", list.GetCarByBrand).Methods(http.MethodGet)
	r.HandleFunc("/car", list.CreateCar).Methods(http.MethodPost)
	r.HandleFunc("/car/del/{id}", list.DeleteCar).Methods(http.MethodDelete)
	r.HandleFunc("/car/upd/{id}", list.UpdateCar).Methods(http.MethodPut)
	r.Use(middleware.Auth)

	err := http.ListenAndServe("localhost:2000", r)
	if err != nil {
		log.Println("Cant Connect to port 2000", err)
	}
}

//TODO remove reflect.deepequal    DONE
//TODO pass car by value in update  DONE
//TODO print errros
//TODO REMOVE C.OUT     DONE
