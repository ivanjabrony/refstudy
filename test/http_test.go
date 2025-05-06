package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUsers(t *testing.T) {
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", getUsers)
	mux.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var responseUsers []User
	if err := json.NewDecoder(resp.Body).Decode(&responseUsers); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(responseUsers) != 2 {
		t.Errorf("Expected 2 users, got %d", len(responseUsers))
	}
}

func TestCreateUser(t *testing.T) {
	newUser := User{ID: 3, Name: "Charlie"}
	body, _ := json.Marshal(newUser)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", createUser)
	mux.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", resp.Status)
	}

	var createdUser User
	if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if createdUser.ID != newUser.ID || createdUser.Name != newUser.Name {
		t.Errorf("Expected user %v, got %v", newUser, createdUser)
	}

	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}
}

func TestCreateUserInvalidJSON(t *testing.T) {
	invalidJSON := `{"id": "aboba", "name": "aboba"}`

	req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", createUser)
	mux.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest, got %v", resp.Status)
	}
}
