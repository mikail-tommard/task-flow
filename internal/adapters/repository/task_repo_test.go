package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/mikail-tommard/task-flow/internal/config"
	"github.com/mikail-tommard/task-flow/internal/domain"

	_ "github.com/lib/pq"
)

//"postgres://taskflow:taskflow@localhost:5433/taskflow?sslmode=disable"

func TestRepo_CreateAndGetByID(t *testing.T) {
	cfg := config.New()

	t.Log("DB_PORT:", cfg.DBPort)
	t.Log("DB_USER:", cfg.DBUser)
	t.Log("DB_NAME:", cfg.DBName)

	dsn := fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := New(db)

	ctx := context.Background()

	task, err := domain.New("test title", "test desc", 1)
	if err != nil {
		t.Fatal(err)
	}

	id, err := repo.Create(ctx, task)
	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID() != id {
		t.Fatalf("expected id=%d, got=%d", id, got.ID())
	}

	if got.Title() != "test title" {
		t.Fatalf("expected title=%q, got=%q", "test title", got.Title())
	}

	if got.Done() != false {
		t.Fatalf("expected done=%t, got=%t", false, got.Done())
	}
}
