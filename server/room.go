package main

import (
	"context"
	"fmt"
	"net/http"
)

type Room struct {
	Id       int
	Name     string
	Password string
	lastId   int64
	users    []User
}

type RoomInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RoomOutput struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Users []User `json:"users"`
}

type CreateRoomRequest struct {
	Room RoomInput `json:"room"`
}

type DeleteRoomRequest struct {
	Room RoomInput `json:"room"`
}

func (s *Server) handleGetRooms(w http.ResponseWriter, r *http.Request) error {
	q := `SELECT r.room_id, r.name FROM room r`

	output := make([]RoomOutput, 0)
	rows, err := s.db.Query(context.Background(), q)
	if err != nil {
		fmt.Println("db error:", err.Error())
		return InternalError()
	}
	defer rows.Close()

	for rows.Next() {
		var r RoomOutput
		err := rows.Scan(&r.Id, &r.Name)
		if err != nil {
			fmt.Println("scan error:", err.Error())
			return InternalError()
		}

		output = append(output, r)
	}

	return writeJSON(w, http.StatusOK, output)
}

func (s *Server) handleGetRoomById(w http.ResponseWriter, r *http.Request) error {
	id, err := getPathId("id", r)
	if err != nil {
		return BadRequest()
	}

	q := `SELECT r.room_id, r.name FROM room r WHERE r.room_id = $1`
	row := s.db.QueryRow(context.Background(), q, id)
	
	var room RoomOutput
	err = row.Scan(&room.Id, &room.Name)
	if err != nil {
		return writeJSON(w, http.StatusOK, nil)
	}

	return writeJSON(w, http.StatusOK, room)
}

func (s *Server) handleCreateRoom(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (s *Server) handleJoinRoom(w http.ResponseWriter, r *http.Request) error {
	
	return nil
}

func (s *Server) handleDeleteRoom(w http.ResponseWriter, r *http.Request) error {
	
	return nil
}
