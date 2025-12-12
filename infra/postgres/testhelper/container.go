package testhelper

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // SQL driver for migrations
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	Container *postgres.PostgresContainer
	ConnStr   string
	Pool      *pgxpool.Pool
}

// NewTestDB creates a PostgreSQL container and runs migrations.
func NewTestDB(ctx context.Context) (*TestDB, error) {
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := runMigrations(connStr); err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	// Create pgxpool
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	return &TestDB{
		Container: container,
		ConnStr:   connStr,
		Pool:      pool,
	}, nil
}

func (t *TestDB) Close(ctx context.Context) error {
	t.Pool.Close()
	return t.Container.Terminate(ctx)
}

// Snapshot creates a snapshot of the current DB state.
func (t *TestDB) Snapshot(ctx context.Context, name string) error {
	return t.Container.Snapshot(ctx, postgres.WithSnapshotName(name))
}

// Restore restores DB to a snapshot.
func (t *TestDB) Restore(ctx context.Context, name string) error {
	return t.Container.Restore(ctx, postgres.WithSnapshotName(name))
}

func runMigrations(connStr string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		return err
	}

	// Get migrations path
	_, filename, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(filename), "..", "migrations")

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
