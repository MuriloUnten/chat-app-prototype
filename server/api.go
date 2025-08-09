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
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

type TokenClaim string
const (
	userIdClaim = TokenClaim("userId")
)

func makeHandler(handler APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request from %s\t%s %s\n", r.RemoteAddr, r.Method, r.URL.Path)

		err := handler(w, r)
		if err != nil {
			if e, ok := err.(APIError); ok {
				fmt.Println("API error:", e.Msg)
				writeJSON(w, e.StatusCode, e)
			} else {
				fmt.Println("error:", err)
				writeJSON(w, http.StatusInternalServerError, "Internal Error")
			}
		}
	}
}

var jwtSecret = []byte("TODO: stop using me")

type CustomClaims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

func jwtMiddleware(handler APIFunc) APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var tokenString string
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && !strings.HasPrefix(authHeader, "Bearer ") {
			return writeJSON(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		}
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")

		// TODO refactor this (workaround for receiving token from ws request)
		if authHeader == "" {
			tokenString = r.Header.Get("Sec-WebSocket-Protocol")
			if tokenString == "" {
				return writeJSON(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
			}
		}

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return writeJSON(w, http.StatusUnauthorized, "Invalid token")
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			return writeJSON(w, http.StatusUnauthorized, "Invalid token")
		}

		ctx := context.WithValue(r.Context(), userIdClaim, claims.UserId)
		r = r.WithContext(ctx)

		return handler(w, r)
	}
}

func createJWT(userId int) (string, error) {
	claims := CustomClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func getIdFromToken(r *http.Request) (int, error) {
	id, ok := r.Context().Value(userIdClaim).(int)
	if !ok {
		return 0, fmt.Errorf("unable to retrieve id from context")
	}
	return id, nil
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
	port         string
	db           *pgxpool.Pool
	websocketHub *Hub
	rooms        map[int]RoomMemberMap
	roomsMutex   sync.RWMutex
}

func NewServer(port string) *Server {
	s := &Server{
		port: port,
		rooms: make(map[int]RoomMemberMap),
	}
	s.websocketHub = NewHub(s)

	s.initDB()
	s.populateRooms()

	http.HandleFunc("GET /api/room", makeHandler(s.handleGetRooms))
	http.HandleFunc("GET /api/room/{id}", makeHandler(s.handleGetRoomById))
	http.HandleFunc("GET /api/room/{roomId}/users", makeHandler(jwtMiddleware(s.handleGetRoomMembers)))
	http.HandleFunc("POST /api/room", makeHandler(jwtMiddleware(s.handleCreateRoom)))
	http.HandleFunc("POST /api/join-room", makeHandler(jwtMiddleware(s.handleJoinRoom)))
	http.HandleFunc("DELETE /api/delete-room", makeHandler(jwtMiddleware(s.handleDeleteRoom)))

	http.HandleFunc("POST /api/user", makeHandler(s.handleCreateUser))
	http.HandleFunc("GET /api/user/{userId}", makeHandler(s.handleGetUserById))
	http.HandleFunc("GET /api/user/rooms", makeHandler(jwtMiddleware(s.handleGetJoinedRoomsByUserId)))
	http.HandleFunc("DELETE /api/user", makeHandler(jwtMiddleware(s.handleDeleteUser)))

	http.HandleFunc("POST /api/login", makeHandler(s.handleUserLogin))

	http.HandleFunc("GET /api/ws", makeHandler(jwtMiddleware(s.handleWebSocket)))

	return s
}

func (s *Server) Run() {
	fmt.Println("Server running on", s.port)
	go s.websocketHub.Run(s)
	err := http.ListenAndServe(s.port, nil)
	s.db.Close()
	log.Fatal(err)
}

func (s *Server) initDB() {
	var err error
	s.db, err = pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err = s.db.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) populateRooms() {
	q := `SELECT room_id, user_id FROM room_user`
	rows, err := s.db.Query(context.Background(), q)
	if err != nil {
		log.Fatal(err)
	}

	var roomId, userId int
	for rows.Next() {
		rows.Scan(&roomId, &userId)
		room, ok := s.rooms[roomId]
		if !ok {
			room = make(RoomMemberMap)
		}
		room[userId] = true
		s.rooms[roomId] = room
	}
}
