package controller

import (
	"context"
	databaseClient "github.com/golang-base-template/util/database/client"
)

var (
	dbGbtMaster databaseClient.DatabaseList
)

type (
	GbtEmployee interface {
		NewGbtEmployee(ctx context.Context, name, gender, address string) GbtEmployee
		SetName(ctx context.Context, name string)
		SetGender(ctx context.Context, gender string)
		SetAddress(ctx context.Context, name string)
		Save(ctx context.Context) (err error)
	}
	gbtEmployee struct {
		Data GbtEmployeeData
	}

	GbtEmployeeData struct {
		EmployeeID int64  `json:"employee_id"`
		Name       string `json:"name"`
		Gender     string `json:"gender"`
		Address    string `json:"address"`
	}
)

func InitGbtEmployeeCore(ctx context.Context) GbtEmployee {
	return &gbtEmployee{}
}

func (ge *gbtEmployee) NewGbtEmployee(ctx context.Context, name, gender, address string) GbtEmployee {
	return &gbtEmployee{
		Data: GbtEmployeeData{
			Name:    name,
			Gender:  gender,
			Address: address,
		}}
}

func (ge *gbtEmployee) SetName(ctx context.Context, name string) {
	ge.Data.Name = name
	return
}

func (ge *gbtEmployee) SetGender(ctx context.Context, gender string) {
	ge.Data.Gender = gender
	return
}

func (ge *gbtEmployee) SetAddress(ctx context.Context, address string) {
	ge.Data.Address = address
	return
}

func (ge *gbtEmployee) Save(ctx context.Context) (err error) {
	//TODO: validate data

	//TODO: build query

	//TODO: exec query

	//TODO: set return result
	return nil
}
