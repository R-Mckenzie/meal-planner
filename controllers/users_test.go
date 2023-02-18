package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/R-Mckenzie/meal-planner/assert"
	"github.com/R-Mckenzie/meal-planner/models"
)

type mockService struct {
	models.UserDB
}

func (m mockService) GenerateRemember(u *models.User) error {
	u.Remember = "test remember"
	return nil
}
func (m mockService) Authenticate(email, pass string) (*models.User, error) {
	return nil, nil
}

func TestUpdateRememberToken(t *testing.T) {
	mock := mockService{}
	uc := User{us: mock}
	rr := httptest.NewRecorder()

	u := &models.User{}
	uc.updateRememberToken(rr, u)
	assert.Equals(t, u.Remember, "test remember")
}
