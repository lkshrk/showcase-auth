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

	mockUserRouteHandler := mocks.NewMockUserRouteHandler(mockCtrl)

	handler := app.SetupHandlers(mockUserRouteHandler)
	w := httptest.NewRecorder()

	t.Run("call login route", func(t *testing.T) {

		req, err := http.NewRequest("POST", "/login", nil)
		if err != nil {
			return
		}

		mockUserRouteHandler.EXPECT().Login(gomock.Any(), gomock.Any())

		handler.ServeHTTP(w, req)

	})

	t.Run("call register router", func(t *testing.T) {

		req, err := http.NewRequest("POST", "/register", nil)
		if err != nil {
			return
		}

		mockUserRouteHandler.EXPECT().CreateUser(gomock.Any(), gomock.Any())

		handler.ServeHTTP(w, req)

	})

}
