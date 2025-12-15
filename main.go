package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/viniciuspra/rssagg/internal/db"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *db.Queries
}

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env is missing")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL env is missing")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("can't connect to database:", err)
	}

	apiCfg := apiConfig{
		DB:  db.New(conn),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://", "http://"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handerGetFeeds)

	v1Router.Post("/feedFollows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feedFollows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr: ":"+port,
		Handler: router,
	}
	log.Println("running server on port: ", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("server error: ", err.Error())
	}
}
