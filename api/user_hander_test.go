package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/api"
	"github.com/mkabdelrahman/hotel-reservation/db"
	"github.com/mkabdelrahman/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testmongodbURI = "mongodb://localhost:27017"
)

type testdb struct {
	db.UserStore
}

func setupTestDB(t *testing.T) *testdb {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(testmongodbURI))
	if err != nil {
		t.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client, "hotel-reservation-test", "user"),
	}
}

func (tdb *testdb) teardown(t *testing.T) {
	err := tdb.UserStore.Drop(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostUser(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.teardown(t)

	t.Run("PostUser_Success", func(t *testing.T) {
		engine := gin.New()

		userHandler := api.NewUserHandler(tdb.UserStore)
		engine.POST("/", userHandler.HandlePostUser)

		params := types.CreateUserParams{
			FirstName: "firstName",
			LastName:  "lastName",
			Email:     "test@gmail.com",
			Password:  "password",
		}

		b, _ := json.Marshal(params)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		// Check the response status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var userResponse types.User
		err := json.Unmarshal(w.Body.Bytes(), &userResponse)
		if err != nil {
			t.Errorf("Error unmarshalling response: %v", err)
		}

		// Additional assertions based on the response content
		assertEqual(t, params.FirstName, userResponse.FirstName, "First name mismatch")
		assertEqual(t, params.LastName, userResponse.LastName, "Last name mismatch")
		assertEqual(t, params.Email, userResponse.Email, "Email mismatch")

		if userResponse.ID.IsZero() {
			t.Error("Expected non-empty user ID, got empty")
		}

		if len(userResponse.EncryptedPassword) > 0 {
			t.Error("Expected encrypted password not to be present in the response")
		}

		fmt.Printf("%+v",userResponse)
	})

}

func assertEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s\nExpected: %v\nActual: %v", message, expected, actual)
	}
}
