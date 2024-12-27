package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Register the PostgreSQL driver
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

var (
	pg *postgres.PostgresContainer
	DB *runner.DB
)

func init() {
	// Setup environment so that Testcontainers gracefully works with colima
	if err := os.Setenv("DOCKER_HOST", "unix:///Users/romain.marcadier/.colima/default/docker.sock"); err != nil {
		log.Fatalln(err)
	}
	if err := os.Setenv("TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE", "/var/run/docker.sock"); err != nil {
		log.Fatalln(err)
	}

	// Start up the test PostgreSQL container
	ctx := context.Background()
	if pgContainer, err := postgres.Run(
		ctx,
		"postgres:alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	); err != nil {
		log.Fatalln(err)
	} else {
		pg = pgContainer
	}

	// Create DB connection
	db, err := sql.Open("pgx", pg.MustConnectionString(ctx))
	if err != nil {
		log.Fatalln(err)
	}

	runner.MustPing(db)

	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(16)

	// Set up table & dummy data
	if _, err := db.ExecContext(ctx, createTalePosts); err != nil {
		log.Fatalln(err)
	}
	if _, err := db.ExecContext(ctx, insertPost); err != nil {
		log.Fatalln(err)
	}

	// Create DB runner
	dat.EnableInterpolation = true
	dat.Strict = false
	runner.LogQueriesThreshold = 10 * time.Millisecond
	DB = runner.NewDB(db, "postgres")
}

const createTalePosts = `
CREATE TABLE posts(
	id         SERIAL    NOT NULL PRIMARY KEY,
	title      TEXT      NOT NULL,
	body       TEXT      NOT NULL,
	user_id    INTEGER   NOT NULL,
	state      TEXT      NOT NULL,
	updated_at TIMESTAMP NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

const insertPost = `
INSERT INTO posts
	("title", "body", "user_id", "state")
VALUES
	('Post #1', 'Body of post #1', 1, 'published'),
	('Post #2', 'Body of post #2', 1, 'published'),
	('Post #3', 'Body of post #3', 1, 'published'),
	('Post #4', 'Body of post #4', 1, 'published'),
	('Post #5', 'Body of post #5', 1, 'published')
;
`
