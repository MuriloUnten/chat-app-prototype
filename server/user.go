package main

import (
	"context"
	// "encoding/json"
	"net/http"
)

type UserHandler struct {
	
}

type User struct {
	Id       int64
	Name     string
	Password string
	// ... conn stuff maybe
}

type UserInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserOutput struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CreateUserRequest struct {
	User UserInput
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	
	return nil
}

func (s *Server) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
	id, err := getPathId("id", r)
	if err != nil {
		return BadRequest()
	}

	q := `SELECT u.user_id, u.name from app_user u WHERE u.user_id = $1`
	row := s.db.QueryRow(context.Background(), q, id)
	
	var u UserOutput
	err = row.Scan(&u.Id, &u.Name)
	if err != nil {
		return writeJSON(w, http.StatusOK, nil)
	}

	return writeJSON(w, http.StatusOK, u)
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {

	return nil
}
