package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Room struct {
	Id       int
	Name     string
	Password string
	Private  bool
	OwnerId  int
	users    []User
}

type RoomInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Private  bool   `json:"private"`
}

type RoomOutput struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Private bool   `json:"private"`
	OwnerId int    `json:"owner_id"`
	Users   []User `json:"users"`
}

type CreateRoomRequest struct {
	Room RoomInput `json:"room"`
}

func (r CreateRoomRequest) validate() map[string]string {
	errs := make(map[string]string)

	if r.Room.Private && len(r.Room.Password) == 0 {
		errs["password"] = "private room must have a password"
	}
	if len(r.Room.Name) == 0 {
		errs["name"] = "room must have a name"
	}

	return errs
}

type CreateRoomResponse struct {
	Room RoomOutput `json:"room"`
}

type DeleteRoomRequest struct {
	RoomId int `json:"room_id"`
}

func (s *Server) handleGetRooms(w http.ResponseWriter, r *http.Request) error {
	q := `SELECT r.room_id, r.name, r.private, r.owner_id FROM room r`

	output := make([]RoomOutput, 0)
	rows, err := s.db.Query(context.Background(), q)
	if err != nil {
		fmt.Println("db error:", err.Error())
		return InternalError()
	}
	defer rows.Close()

	for rows.Next() {
		var room RoomOutput
		err := rows.Scan(&room.Id, &room.Name, &room.Private, &room.OwnerId)
		if err != nil {
			fmt.Println("scan error:", err.Error())
			return InternalError()
		}

		output = append(output, room)
	}

	return writeJSON(w, http.StatusOK, output)
}

func (s *Server) handleGetRoomById(w http.ResponseWriter, r *http.Request) error {
	id, err := getPathId("id", r)
	if err != nil {
		return BadRequest()
	}

	q := `SELECT r.room_id, r.name, r.private, r.owner_id FROM room r WHERE r.room_id = $1`
	row := s.db.QueryRow(context.Background(), q, id)
	
	var room RoomOutput
	err = row.Scan(&room.Id, &room.Name, &room.Private, &room.OwnerId)
	if err != nil {
		return writeJSON(w, http.StatusOK, nil)
	}

	return writeJSON(w, http.StatusOK, room)
}

func (s *Server) handleCreateRoom(w http.ResponseWriter, r *http.Request) error {
	userId, err := getIdFromToken(r)
	if err != nil {
		return UserNotAuthenticated()
	}

	var req CreateRoomRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return BadRequest()
	}

	errs := req.validate()
	if len(errs) > 0 {
		return InvalidJSONRequestData(errs)
	}

	if req.Room.Private {
		// TODO encrypt room password
	}

	q := `INSERT INTO room(name, password_hash, private, owner_id) VALUES($1, $2, $3, $4) RETURNING room_id, name, private, owner_id`
	row := s.db.QueryRow(context.Background(), q, req.Room.Name, req.Room.Password, req.Room.Private, userId)

	var resp CreateRoomResponse
	err = row.Scan(&resp.Room.Id, &resp.Room.Name, &resp.Room.Private, &resp.Room.OwnerId)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleJoinRoom(w http.ResponseWriter, r *http.Request) error {
	
	return nil
}

func (s *Server) handleDeleteRoom(w http.ResponseWriter, r *http.Request) error {
	
	return nil
}
