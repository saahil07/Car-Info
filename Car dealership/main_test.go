package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/models"
)

// Test_Main this is a test function for main
func Test_Main(t *testing.T) {
	c := http.Client{}

	go main()
	time.Sleep(time.Second * 3)

	id1 := "f96f3281-fca7-4610-a48c-4ae197d6fd2a"
	id := uuid.New()
	c1 := models.Car{ID: id, Name: "AQ 3", Year: 2000, Brand: "Mercedes", FuelType: "diesel",
		Engine: models.Engine{EngineID: id, Displacement: 400, NoOfCylinder: 6}}

	testcases := []struct {
		desc       string
		method     string
		pathparam  string
		queryparam string
		body       models.Car
		status     int
	}{
		{
			desc:      "create car",
			method:    http.MethodPost,
			pathparam: "car",
			body:      c1,
			status:    http.StatusCreated,
		},
		{desc: "get by id",
			method:    http.MethodGet,
			pathparam: "car/" + id1,
			status:    http.StatusOK,
		},
		{desc: "get by brand",
			method:     http.MethodGet,
			queryparam: "cars?brand=Ferrari&isEngine=true",
			status:     http.StatusOK,
		},
		{desc: "update car",
			method:    http.MethodPut,
			pathparam: "car/upd/" + id.String(),
			body:      c1,
			status:    http.StatusOK,
		},
		{desc: "delete car",
			method:    http.MethodDelete,
			pathparam: "car/del/" + id.String(),
			status:    http.StatusOK,
		},
	}

	for i, tc := range testcases {
		body, _ := json.Marshal(tc.body)
		req, err := http.NewRequest(tc.method, "http://localhost:2000/"+tc.pathparam+tc.queryparam, bytes.NewReader(body))

		if err != nil {
			log.Println(err)
			return
		}

		req.Header.Set("authorize", "0000")

		res, err := c.Do(req)
		if err != nil {
			log.Println(err)
			return
		}

		if tc.status != res.StatusCode {
			t.Errorf("testcase %v failed\n desc: %v\tExpected : %v\tGot: %v", i, tc.desc, tc.status, res.StatusCode)

			res.Body.Close()
		}
	}
}
