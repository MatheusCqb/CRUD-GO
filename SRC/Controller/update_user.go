package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/HunCoding/meu-primeiro-crud-go/src/configuration/rest_err"
	"github.com/HunCoding/meu-primeiro-crud-go/src/controller/model/request"
	"github.com/HunCoding/meu-primeiro-crud-go/src/model"
	"github.com/HunCoding/meu-primeiro-crud-go/src/tests/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

func TestUserControllerInterface_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockUserDomainService(ctrl)
	controller := NewUserControllerInterface(service)

	t.Run("validation_got_error", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		context := GetTestGinContext(recorder)

		userRequest := request.UserUpdateRequest{
			Name: "Test User",
			Age:  -10, // Age is invalid
		}

		b, _ := json.Marshal(userRequest)
		stringReader := io.NopCloser(strings.NewReader(string(b)))

		param := []gin.Param{
			{
				Key:   "userId",
				Value: "test",
			},
		}

		MakeRequest(context, param, url.Values{}, "PUT", stringReader)
		controller.UpdateUser(context)

		assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("user_id_is_invalid_returns_error", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		context := GetTestGinContext(recorder)

		userRequest := request.UserUpdateRequest{
			Name: "Test User",
			Age:  20,
		}

		b, _ := json.Marshal(userRequest)
		stringReader := io.NopCloser(strings.NewReader(string(b)))

		param := []gin.Param{
			{
				Key:   "userId",
				Value: "test",
			},
		}

		MakeRequest(context, param, url.Values{}, "PUT", stringReader)
		controller.UpdateUser(context)

		assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("object_is_valid_but_service_returns_error", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		context := GetTestGinContext(recorder)

		userRequest := request.UserUpdateRequest{
			Name: "Test User",
			Age:  20,
		}

		domain := model.NewUserUpdateDomain(
			userRequest.Name,
			userRequest.Age,
		)

		b, _ := json.Marshal(userRequest)
		stringReader := io.NopCloser(strings.NewReader(string(b)))

		param := []gin.Param{
			{
				Key:   "userId",
				Value: primitive.NewObjectID().Hex(),
			},
		}

		service.EXPECT().UpdateUser(primitive.NewObjectID().Hex(), domain).Return(
			rest_err.NewInternalServerError("error test"))

		MakeRequest(context, param, url.Values{}, "PUT", stringReader)
		controller.UpdateUser(context)

		assert.EqualValues(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("object_is_valid_and_service_returns_success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		context := GetTestGinContext(recorder)

		userRequest := request.UserUpdateRequest{
			Name: "Test User",
			Age:  20,
		}

		domain := model.NewUserUpdateDomain(
			userRequest.Name,
			userRequest.Age,
		)

		b, _ := json.Marshal(userRequest)
		stringReader := io.NopCloser(strings.NewReader(string(b)))

		id := primitive.NewObjectID().Hex()

		param := []gin.Param{
			{
				Key:   "userId",
				Value: id,
			},
		}

		service.EXPECT().UpdateUser(id, domain).Return(nil)

		MakeRequest(context, param, url.Values{}, "PUT", stringReader)
		controller.UpdateUser(context)

		assert.EqualValues(t, http.StatusOK, recorder.Code)
	})
}
