package car

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/models"
	"github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/service"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

//type mockBody httptest.ResponseRecorder
//
//func (m mockBody) Write([]byte) (n int, e error) {
//	m.Code = http.StatusBadRequest
//	return 0, errors.New("write error")
//}

type ResponseR struct {
	Code int

	HeaderMap http.Header

	Body *bytes.Buffer

	Flushed bool

	result      *http.Response // cache of Result's return value
	snapHeader  http.Header    // snapshot of HeaderMap at first Write
	wroteHeader bool
}

func (r ResponseR) Header() http.Header {
	return r.HeaderMap
}

func (r ResponseR) WriteHeader(statusCode int) {
	//r.Code = statusCode

	r.WriteHeader(statusCode)
}

func (r ResponseR) Write(buf []byte) (n int, e error) {
	return len(buf), errors.New("WRITE ERR")
}

func NewRr() *ResponseR {
	return &ResponseR{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
		Code:      200,
	}
}

func TestWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mocks := service.NewMockCars(ctrl)
	mockhandler := New(mocks)
	id := uuid.New()

	testCar := models.Car{
		ID: id, Engine: models.Engine{EngineID: id, Displacement: 500, NoOfCylinder: 2, CarRange: 200},
		Name: "Model 3", Year: 2018, Brand: "Tesla", FuelType: "petrol"}

	testcases := []struct {
		desc       string
		statuscode int
		err        error
	}{
		{"write", http.StatusBadRequest,
			errors.New("write error")},
	}

	mocks.EXPECT().GetCarByID(gomock.Any(), id.String()).Return(testCar, nil)
	for _, tc := range testcases {

		r := httptest.NewRequest(http.MethodGet, "/car/{id}"+id.String(), nil)
		w := NewRr()

		r = mux.SetURLVars(r, map[string]string{"id": id.String()})

		mockhandler.GetCarByID(w, r)

		fmt.Println(w.Code)

		assert.Equal(t, tc.statuscode, w.Code)
	}

}

// TestGetByID handler layer test function to test handler layer GetbyId function
func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := service.NewMockCars(ctrl)
	s := New(mockService)

	id1 := uuid.New()
	id2 := uuid.New()

	testCar := models.Car{
		ID: id1, Engine: models.Engine{EngineID: id1, Displacement: 500, NoOfCylinder: 2, CarRange: 200},
		Name: "Model 3", Year: 2018, Brand: "Tesla", FuelType: "petrol"}

	testCases := []struct {
		desc       string
		id         uuid.UUID
		statusCode int
		mock       []*gomock.Call
	}{
		{
			desc:       "success case",
			id:         id1,
			statusCode: http.StatusOK,
			mock: []*gomock.Call{mockService.EXPECT().GetCarByID(gomock.Any(), id1.String()).
				Return(testCar, nil)}},
		{
			desc:       "not found",
			id:         id2,
			statusCode: http.StatusInternalServerError,
			mock: []*gomock.Call{
				mockService.EXPECT().GetCarByID(gomock.Any(), id2.String()).
					Return(models.Car{}, errors.New("invalid Id")),
			},
		},
		{
			desc:       "not found",
			id:         uuid.Nil,
			statusCode: http.StatusInternalServerError,
			mock: []*gomock.Call{
				mockService.EXPECT().GetCarByID(gomock.Any(), uuid.Nil.String()).
					Return(models.Car{}, errors.New("invalid Id")),
			},
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest("GET", "/car/{id}"+tc.id.String(), nil)

		res := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{
			"id": tc.id.String(),
		})

		s.GetCarByID(res, req)

		if res.Code != tc.statusCode {
			t.Errorf("Expected Status Code: %v, Got: %v", tc.statusCode, res.Code)
		}
	}
}

// TestCarGetbyBrand handler layer test function to test handler layer GetbyBrand function
func TestCarGetbyBrand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := service.NewMockCars(ctrl)
	s := New(mockService)

	var (
		id  = uuid.New()
		id1 = uuid.New()
	)

	cars := []models.Car{
		{ID: id, Name: "model x",
			Engine: models.Engine{EngineID: uuid.MustParse("38ec1d7a-834f-11ec-a8a3-0242ac120002"), Displacement: 0,
				NoOfCylinder: 0, CarRange: 450}, Year: 2014, Brand: "Tesla", FuelType: "electric"},
		{ID: id1, Name: "model 3",
			Engine: models.Engine{EngineID: uuid.MustParse("5c62818a-834b-11ec-a8a3-0242ac120002"),
				Displacement: 0, NoOfCylinder: 0, CarRange: 498}, Year: 2017, Brand: "Tesla", FuelType: "electric"},
	}

	testCases := []struct {
		desc       string
		brand      string
		isEngine   string
		statusCode int
		mock       []*gomock.Call
	}{
		{desc: "success case", brand: "Tesla", isEngine: "true", statusCode: http.StatusOK,
			mock: []*gomock.Call{mockService.EXPECT().GetCarByBrand(gomock.Any(), "Tesla", true).
				Return(cars, nil)},
		},
		{
			desc: "error", brand: "Maruti", isEngine: "false", statusCode: http.StatusInternalServerError,
			mock: []*gomock.Call{mockService.EXPECT().GetCarByBrand(gomock.Any(), "Maruti", false).
				Return([]models.Car{}, errors.New("error"))},
		},
		{
			desc: "error", brand: "Maruti", isEngine: "hello", statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest("GET", "/cars"+"?brand="+tc.brand+"&isEngine="+tc.isEngine, nil)
		res := httptest.NewRecorder()

		s.GetCarByBrand(res, req)

		if res.Code != tc.statusCode {
			t.Errorf("Expected Status Code: %v, Got: %v", tc.statusCode, res.Code)
		}
	}
}

// TestCreateCar handler layer test function to test handler layer Create function
func TestCreateCar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := service.NewMockCars(ctrl)
	s := New(mockService)

	id := uuid.New()
	car := models.Car{ID: id,
		Name: "model x", Engine: models.Engine{EngineID: id,
			Displacement: 0, NoOfCylinder: 0, CarRange: 450},
		Year: 2014, Brand: "Tesla", FuelType: "electric"}

	testCases := []struct {
		desc       string
		car        models.Car
		statusCode int
		err        error
	}{
		{desc: "empty body", car: car, statusCode: http.StatusInternalServerError, err: errors.New("error")},
		{desc: "success", car: car, statusCode: http.StatusCreated, err: nil},
		{desc: "fail", car: car, statusCode: http.StatusBadRequest, err: errors.New("error")},
		{desc: "fail", car: car, statusCode: http.StatusBadRequest, err: errors.New("error")},
	}

	gomock.InOrder(
		mockService.EXPECT().CreateCar(gomock.Any(), &car).Return(car, nil),
		mockService.EXPECT().CreateCar(gomock.Any(), &car).Return(models.Car{}, errors.New("error")),
		mockService.EXPECT().CreateCar(gomock.Any(), &car).Return(models.Car{}, errors.New("body error")),
	)

	for i, tc := range testCases {
		body, _ := json.Marshal(tc.car)
		if i == 0 {
			body = nil
		}

		req := httptest.NewRequest("POST", "/car/", bytes.NewBuffer(body))
		res := httptest.NewRecorder()

		s.CreateCar(res, req)

		if res.Code != tc.statusCode {
			t.Errorf("Expected Status Code: %v, Got: %v", tc.statusCode, res.Code)
		}
	}
}

// TestUpdateCar handler layer test function to test handler layer Update function
func TestUpdateCar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := service.NewMockCars(ctrl)
	s := New(mockService)

	id := uuid.New()
	car := models.Car{
		ID: id,
		Engine: models.Engine{EngineID: id, Displacement: 200,
			NoOfCylinder: 1, CarRange: 222}, Name: "BMW", Year: 2018, Brand: "BMW", FuelType: "petrol"}

	id1 := uuid.New()

	testCases := []struct {
		desc       string
		id         uuid.UUID
		car        models.Car
		statusCode int
		err        error
		mock       *gomock.Call
	}{
		{desc: "empty body", id: id1, car: car, statusCode: http.StatusInternalServerError,
			err: errors.New("error"), mock: nil},
		{desc: "success", id: id1, car: car, statusCode: http.StatusOK, err: nil,
			mock: mockService.EXPECT().UpdateCar(gomock.Any(), id1.String(), car).Return(car, nil)},
		{desc: "fail", id: id1, car: car, statusCode: http.StatusBadRequest, err: errors.New("error"),
			mock: mockService.EXPECT().UpdateCar(gomock.Any(), id1.String(), car).
				Return(car, errors.New("error"))},
	}

	for i, tc := range testCases {
		body, _ := json.Marshal(tc.car)
		if i == 0 {
			body = []byte("")
		}

		req := httptest.NewRequest("PUT", "/car/update/{id}"+tc.id.String(), bytes.NewBuffer(body))
		res := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{
			"id": tc.id.String(),
		})
		s.UpdateCar(res, req)

		if res.Code != tc.statusCode {
			t.Errorf("Expected Status Code: %v, Got: %v", tc.statusCode, res.Code)
		}
	}
}

// TestDeleteCar handler layer test function to test handler layer Delete function
func TestDeleteCar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := service.NewMockCars(ctrl)
	s := New(mockService)

	id1 := uuid.New()
	id2 := uuid.New()

	testCases := []struct {
		id         uuid.UUID
		statusCode int
		err        error
		mock       *gomock.Call
	}{
		{id: id1, statusCode: http.StatusOK, err: nil, mock: mockService.EXPECT().
			DeleteCar(gomock.Any(), id1.String()).Return(models.Car{}, nil)},
		{id: id2, statusCode: http.StatusBadRequest, err: errors.New("invalid id"),
			mock: mockService.EXPECT().DeleteCar(gomock.Any(), id2.String()).
				Return(models.Car{}, errors.New("invalid id"))},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest("DELETE", "/car/delete/{id}"+tc.id.String(), nil)
		res := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{
			"id": tc.id.String(),
		})
		s.DeleteCar(res, req)

		if res.Code != tc.statusCode {
			t.Errorf("Expected Status Code: %v, Got: %v", tc.statusCode, res.Code)
		}
	}
}
