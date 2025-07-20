package main

import (
	"context"
	"encoding/json"

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

type CreateUserResponse struct {
	Token string     `json:"token"`
	User  UserOutput `json:"user"`
}

func (r CreateUserRequest) validate() map[string]string {
	errs := make(map[string]string)
	if len(r.User.Name) < 3 {
		errs["name"] = "name must be at least 3 characters long"
	}

	if len(r.User.Password) < 12 {
		errs["password"] = "password must be at least 12 characters long"
	}
	// TODO add better validation for weak passwords

	return errs
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req.User)
	if err != nil {
		return BadRequest()
	}

	errs := req.validate()
	if len(errs) > 0 {
		return NewAPIError(http.StatusUnprocessableEntity, errs)
	}
	
	q := `INSERT INTO app_user(name, password_hash) VALUES($1, $2) RETURNING user_id, name`
	row := s.db.QueryRow(context.Background(), q, &req.User.Name, &req.User.Password)

	var user UserOutput
	err = row.Scan(&user.Id, &user.Name)
	if err != nil {
		return err
	}

	jwt, err := createJWT(user.Id)
	if err != nil {
		return err
	}

	response := CreateUserResponse{
		Token: jwt,
		User: user,
	}
	
	return writeJSON(w, http.StatusOK, response)
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
	id, err := getIdFromToken(r)
	if err != nil {
		return BadRequest()
	}

	// TODO manage jwt expiration

	q := `DELETE FROM app_user WHERE user_id = $1`
	_, err = s.db.Exec(context.Background(), q, id)
	if err != nil {
		return err
	}

	// TODO should probably double check what happened if no rows affected
	// if cmdTag.RowsAffected == 0 {
	//
	// }

	return nil
}
