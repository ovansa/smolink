// test/setup_test.go
package test

import (
	"context"
	"fmt"
	"os"
	"time"

	"smolink/internal/app"
	"smolink/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"
)

type TestApp struct {
	*app.App      // Embed the actual App struct
	pool          *dockertest.Pool
	pgResource    *dockertest.Resource
	redisResource *dockertest.Resource
}

func SetupTestApp() *TestApp {
	testApp := &TestApp{}
	var err error

	// Docker setup
	testApp.pool, err = dockertest.NewPool("")
	if err != nil {
		panic("Docker pool setup failed: " + err.Error())
	}
	testApp.pool.MaxWait = 60 * time.Second

	// PostgreSQL container
	testApp.pgResource, err = testApp.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres", Tag: "15-alpine",
		Env: []string{"POSTGRES_USER=test", "POSTGRES_PASSWORD=test", "POSTGRES_DB=testdb"},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})
	if err != nil {
		panic("Postgres container failed: " + err.Error())
	}

	// Redis container
	testApp.redisResource, err = testApp.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis", Tag: "7-alpine",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})
	if err != nil {
		panic("Redis container failed: " + err.Error())
	}

	// Set environment variables
	os.Setenv("ENV", "test")
	os.Setenv("POSTGRES_DSN", fmt.Sprintf("postgres://test:test@localhost:%s/testdb?sslmode=disable", testApp.pgResource.GetPort("5432/tcp")))
	os.Setenv("REDIS_ADDR", fmt.Sprintf("localhost:%s", testApp.redisResource.GetPort("6379/tcp")))

	// Wait for databases
	testApp.retryConnect()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Test config failed: " + err.Error())
	}

	// Initialize the actual application
	application, err := app.NewApp(cfg, false) // false = no homepage routes
	if err != nil {
		panic("App setup failed: " + err.Error())
	}
	testApp.App = application

	// Initialize database schema
	testApp.initDBSchema()

	return testApp
}

func (ta *TestApp) retryConnect() {
	// PG connection
	err := ta.pool.Retry(func() error {
		pool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_DSN"))
		if err != nil {
			return err
		}
		return pool.Ping(context.Background())
	})
	if err != nil {
		panic("Postgres connection failed: " + err.Error())
	}

	// Redis connection
	err = ta.pool.Retry(func() error {
		client := redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_ADDR"),
			DB:   0,
		})
		return client.Ping(context.Background()).Err()
	})
	if err != nil {
		panic("Redis connection failed: " + err.Error())
	}
}

func (ta *TestApp) initDBSchema() {
	_, err := ta.PGRepo.DB().Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_code VARCHAR(20) UNIQUE NOT NULL,
			original_url TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			click_count INTEGER NOT NULL DEFAULT 0
		);
		CREATE TABLE IF NOT EXISTS url_analytics (
			id SERIAL PRIMARY KEY,
			url_id INTEGER NOT NULL REFERENCES urls(id),
			ip_address TEXT NOT NULL,
			user_agent TEXT NOT NULL,
			accessed_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		panic("DB schema setup failed: " + err.Error())
	}
}

func (ta *TestApp) Cleanup() {
	if ta.pool != nil {
		if ta.pgResource != nil {
			_ = ta.pool.Purge(ta.pgResource)
		}
		if ta.redisResource != nil {
			_ = ta.pool.Purge(ta.redisResource)
		}
	}
	if ta.DBCloser != nil {
		_ = ta.DBCloser()
	}
}
