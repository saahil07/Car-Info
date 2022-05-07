package car

import (
	"context"
	"errors"
	"log"
	"reflect"
	"testing"

	"github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

// TestGetByID function to test store layer GetbyId function
func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	a := New(db)

	var (
		id  = uuid.New()
		id1 = uuid.New()

		car1 = models.Car{ID: id, Name: "Q2", Year: 2009, Brand: "BMW", FuelType: "petrol", Engine: models.Engine{
			EngineID: id,
		}}
		er = errors.New("all expectations were already fulfilled")
	)

	testCases := []struct {
		desc      string
		id        uuid.UUID
		outputCar models.Car
		err       error
	}{
		{desc: "success", id: id, outputCar: car1, err: nil},
		{"car ID invalid", id1, models.Car{}, er},
	}

	rows := sqlmock.NewRows([]string{"id", "engine_id", "name", "year", "brand", "fuelType"}).
		AddRow(car1.ID.String(), car1.Engine.EngineID.String(), car1.Name, car1.Year, car1.Brand, car1.FuelType)

	mock.ExpectQuery("SELECT * FROM Car WHERE ID=?;").WithArgs(id).WillReturnRows(rows)
	mock.ExpectQuery("SELECT * FROM Car WHERE ID=?;").WithArgs(id1).WillReturnError(er)

	for i, tc := range testCases {
		resp, err := a.GetCarByID(context.TODO(), tc.id.String())

		if resp != tc.outputCar {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.outputCar)
		}

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

// TestGetbybrand function to test store layer GetbyBrand function
func TestGetbybrand(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	a := New(db)

	defer db.Close()

	var (
		id         = uuid.New()
		id1        = uuid.New()
		id2        = uuid.New()
		queryError = errors.New("query error")
		er         = errors.New("sql: expected 5 destination arguments in Scan, not 6")

		car = models.Car{ID: id, Name: "GenX", Year: 2015, Brand: "Ferrari",
			FuelType: "electric", Engine: models.Engine{EngineID: id}}

		car1 = models.Car{ID: id1, Name: "Ferrari AQ", Year: 2020, Brand: "Ferrari",
			FuelType: "electric", Engine: models.Engine{EngineID: id1}}

		car2 = models.Car{ID: id2, Name: "X4", Brand: "Porsche",
			FuelType: "electric", Engine: models.Engine{EngineID: id2}}
	)

	testCases := []struct {
		desc   string
		brand  string
		output []models.Car
		eng    bool
		err    error
	}{
		{desc: "show all Car", brand: "Ferrari", output: []models.Car{car, car1}, eng: true, err: nil},
		{"no brand name", "", nil, false, queryError},
		{"arguments missing", "Porsche", nil, true, er},
		{"row error", "BMW", []models.Car{}, true, errors.New("err")},
	}

	rows := sqlmock.NewRows([]string{"id", "engine_id", "name", "year", "brand", "fuel_type"}).
		AddRow(id.String(), id.String(), car.Name, car.Year, car.Brand, car.FuelType).
		AddRow(id1.String(), id1.String(), car1.Name, car1.Year, car1.Brand, car1.FuelType)

	rows2 := sqlmock.NewRows([]string{"id", "engine_id", "name", "brand", "fuel_type"}).
		AddRow(id2.String(), id2.String(), car2.Name, car2.Brand, car2.FuelType)

	rows3 := sqlmock.NewRows([]string{"id", "engine_id", "name", "brand", "fuel_type"}).
		AddRow(id2.String(), id2.String(), car2.Name, car2.Year, car2.Brand).RowError(0, errors.New("err"))

	mock.ExpectQuery("select * from Car where brand=?;").WithArgs("Ferrari").WillReturnRows(rows)
	mock.ExpectQuery("select * from Car where brand=?;").WithArgs("").WillReturnError(queryError)
	mock.ExpectQuery("select * from Car where brand=?;").WithArgs("Porsche").WillReturnRows(rows2)
	mock.ExpectQuery("select * from Car where brand=?;").WithArgs("BMW").WillReturnRows(rows3)

	for i, tc := range testCases {
		car, err := a.GetCarsByBrand(context.TODO(), tc.brand, tc.eng)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if !reflect.DeepEqual(car, tc.output) {
			t.Errorf("\n[TEST %v] Failed \n got %v\nGot \n Expected %v", i, car, tc.output)
		}
	}
}

// TestCreatecar function to test store layer Create function
func TestCreatecar(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	a := New(db)

	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	id := uuid.New()
	car := models.Car{ID: id, Name: "AQ", Year: 2015, Brand: "Ferrari",
		FuelType: "electric", Engine: models.Engine{EngineID: id}}

	car1 := models.Car{ID: uuid.Nil, Name: "AQ", Year: 2015, Brand: "Ferrari",
		FuelType: "electric", Engine: models.Engine{EngineID: id}}

	queryErr := errors.New("query error")

	testCases := []struct {
		desc           string
		input          models.Car
		expectedOutput models.Car
		err            error
	}{
		{"success", car, car, nil},
		{"fail", car1, models.Car{}, queryErr},
	}

	mock.ExpectExec("INSERT INTO Car (id,engine_id,name,year,brand,fuel_type) VALUES(?,?,?,?,?,?)").
		WithArgs(sqlmock.AnyArg(), car.Engine.EngineID, car.Name, car.Year, car.Brand, car.FuelType).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO Car (id,engine_id,name,year,brand,fuel_type) VALUES(?,?,?,?,?,?)").
		WithArgs(uuid.Nil, car.Engine.EngineID, car.Name, car.Year, car.Brand, car.FuelType).
		WillReturnError(queryErr)

	for i, tc := range testCases {
		res, err := a.CreateCar(context.TODO(), &tc.input)

		if res != tc.expectedOutput {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, res, tc.expectedOutput)
		}

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

// TestUpdatecar function to test store layer update function
func TestUpdatecar(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	a := New(db)

	id := uuid.New()
	id1 := uuid.Nil

	car := models.Car{ID: id, Name: "BMW", Year: 2018, Brand: "Rolls-Royce", FuelType: "petrol"}
	updateFail := errors.New("update failed")

	car1 := models.Car{ID: id1, Name: "BMW", Year: 2018, Brand: "Rolls-Royce", FuelType: "petrol"}

	testCases := []struct {
		desc  string
		input models.Car
		err   error
	}{
		{"success", car, nil},
		{"failure", car1, updateFail},
	}

	defer db.Close()

	mock.ExpectExec("UPDATE Car SET name=?,year=?,brand=?,fuel_type=? WHERE id=?").
		WithArgs(car.Name, car.Year, car.Brand, car.FuelType, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE Car SET name=?,year=?,brand=?,fuel_type=? WHERE id=?").
		WithArgs(car1.Name, car1.Year, car1.Brand, car1.FuelType, id1).
		WillReturnError(updateFail)

	for i, tc := range testCases {
		_, err := a.UpdateCar(context.TODO(), tc.input.ID.String(), tc.input)
		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

// TestDeleteCar function to test store layer  delete function
func TestDeleteCar(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	a := New(db)

	defer db.Close()

	id1 := uuid.New()
	er := errors.New("delete failed")

	testCases := []struct {
		desc         string
		id           uuid.UUID
		rowsEffected int
		err          error
	}{
		{"Success", id1, 1, nil},
		{"ID does not exists", uuid.Nil, 0, er},
	}

	mock.ExpectExec("DELETE FROM Car WHERE ID=?").WithArgs(id1.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM Car WHERE ID=?").WithArgs(uuid.Nil).
		WillReturnError(er)

	for i, tc := range testCases {
		_, err := a.DeleteCar(context.TODO(), tc.id.String())

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}
