package view

import (
	"github.com/MatheusCqb/CRUD-GO/src/controller/model/response"
	"github.com/MatheusCqb/CRUD-GO/src/model"
)

func ConvertDomainToResponse(
	userDomain model.UserDomainInterface,
) response.UserResponse {
	return response.UserResponse{
		ID:    userDomain.GetID(),
		Email: userDomain.GetEmail(),
		Name:  userDomain.GetName(),
		Age:   userDomain.GetAge(),
	}
}
