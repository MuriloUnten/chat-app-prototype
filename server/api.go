package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func makeHandler(handler APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			if e, ok := err.(APIError); ok {
				fmt.Println("API error:", e.Msg)
				writeJSON(w, e.StatusCode, e.Msg)
			}
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func getPathId(wildcard string, r *http.Request) (int, error) {
	v := r.PathValue(wildcard)
	if v == "" {
		return 0, errors.New("unable to get path id")
	}

	id, err := strconv.Atoi(v)
	return id, err
}

type Server struct {
	port string
	db *pgxpool.Pool
}

func NewServer(port string) *Server {
	s := &Server{
		port: port,
	}

	var err error
	s.db, err = pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("GET /room", makeHandler(s.handleGetRooms))
	http.HandleFunc("GET /room/{id}", makeHandler(s.handleGetRoomById))
	http.HandleFunc("POST /create-room", makeHandler(s.handleCreateRoom))
	http.HandleFunc("POST /join-room", makeHandler(s.handleJoinRoom))
	http.HandleFunc("DELETE /delete-room", makeHandler(s.handleDeleteRoom))

	http.HandleFunc("POST /user", makeHandler(s.handleCreateUser))
	http.HandleFunc("GET /user/{id}", makeHandler(s.handleGetUserById))
	http.HandleFunc("DELETE /user", makeHandler(s.handleDeleteUser))

	return s
}

func (s *Server) Run() {
	fmt.Println("Server running on", s.port)
	err := http.ListenAndServe(s.port, nil)
	s.db.Close()
	log.Fatal(err)
}
