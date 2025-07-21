package main

import (
	"golang.org/x/crypto/bcrypt"
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

	if r.Room.Private && len(r.Room.Password) > 72 {
		errs["password"] = "password must not exceed 72 characters"
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

type JoinRoomRequest struct {
	RoomId   int    `json:"room_id"`
	Password string `json:"password"`
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

	var hash string = ""
	if req.Room.Private {
		hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Room.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		hash = string(hashBytes)
	}

	q := `INSERT INTO room(name, password_hash, private, owner_id) VALUES($1, $2, $3, $4) RETURNING room_id, name, private, owner_id`
	row := s.db.QueryRow(context.Background(), q, req.Room.Name, hash, req.Room.Private, userId)

	var resp CreateRoomResponse
	err = row.Scan(&resp.Room.Id, &resp.Room.Name, &resp.Room.Private, &resp.Room.OwnerId)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleJoinRoom(w http.ResponseWriter, r *http.Request) error {
	userId, err := getIdFromToken(r)
	if err != nil {
		return UserNotAuthenticated()
	}

	var req JoinRoomRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return BadRequest()
	}

	q := `SELECT 1 FROM room_user WHERE room_id = $1 AND user_id = $2`
	err = s.db.QueryRow(context.Background(), q, req.RoomId, userId).Scan(nil)
	if err == nil {
		// TODO Handle user already in the room
		return writeJSON(w, http.StatusOK, "already in the room")
	}

	var room Room
	q = `SELECT r.private, r.password_hash FROM room r WHERE room_id = $1`
	err = s.db.QueryRow(context.Background(), q, req.RoomId).Scan(&room.Private, &room.Password)
	if err != nil {
		return BadRequest()
	}

	if room.Private {
		err := bcrypt.CompareHashAndPassword([]byte(room.Password), []byte(req.Password))
		if err != nil {
			return NewAPIError(http.StatusUnauthorized, "invalid room password")
		}
	}

	q = `INSERT INTO room_user(room_id, user_id) VALUES($1, $2)`
	tag, err := s.db.Exec(context.Background(), q, req.RoomId, userId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return InternalError()
	}

	// TODO Handle sync with chat service
	return nil
}

func (s *Server) handleDeleteRoom(w http.ResponseWriter, r *http.Request) error {
	
	return nil
}
