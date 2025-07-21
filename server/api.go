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
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/golang-jwt/jwt/v5"
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
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return writeJSON(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

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
	port string
	db *pgxpool.Pool
	websocketHub *Hub
}

func NewServer(port string) *Server {
	s := &Server{
		port: port,
		websocketHub: NewHub(),
	}

	var err error
	s.db, err = pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err = s.db.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("GET /room", makeHandler(s.handleGetRooms))
	http.HandleFunc("GET /room/{id}", makeHandler(s.handleGetRoomById))
	http.HandleFunc("POST /room", makeHandler(jwtMiddleware(s.handleCreateRoom)))
	http.HandleFunc("POST /join-room", makeHandler(jwtMiddleware(s.handleJoinRoom)))
	http.HandleFunc("DELETE /delete-room", makeHandler(jwtMiddleware(s.handleDeleteRoom)))

	http.HandleFunc("POST /user", makeHandler(s.handleCreateUser))
	http.HandleFunc("GET /user/{userId}", makeHandler(s.handleGetUserById))
	http.HandleFunc("DELETE /user", makeHandler(jwtMiddleware(s.handleDeleteUser)))

	http.HandleFunc("GET /ws", makeHandler(jwtMiddleware(s.handleWebSocket)))

	return s
}

func (s *Server) Run() {
	fmt.Println("Server running on", s.port)
	go s.websocketHub.Run(s)
	err := http.ListenAndServe(s.port, nil)
	s.db.Close()
	log.Fatal(err)
}
