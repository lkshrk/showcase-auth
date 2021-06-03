package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"harke.me/showcase-auth/pkg/app"
	"harke.me/showcase-auth/pkg/mocks"
)

func TestRoutes(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mocks.NewMockUserService(mockCtrl)

	mockUserRouteHandler := app.NewUserRouteHandler(mockUserService)

	app.RegisterRoutes(mockUserRouteHandler)

	httptest.NewServer()

	t.Run("call login route", func(t *testing.T) {

		req, err := http.NewRequest("POST", "/login", nil)
		if err != nil {
			return
		}

		mockUserRouteHandler.EXPECT().Login(gomock.Any(), gomock.Any())

		http.DefaultClient.Do(req)

	})

	t.Run("call register router", func(t *testing.T) {

		req, err := http.NewRequest("POST", "/register", nil)
		if err != nil {
			return
		}

		mockUserRouteHandler.EXPECT().Login(gomock.Any(), gomock.Any())

		http.DefaultClient.Do(req)

	})

}
