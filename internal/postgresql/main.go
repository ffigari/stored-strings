package postgresql

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ffigari/stored-strings/internal/config"
	"github.com/ffigari/stored-strings/internal/oos"
)

func CreateEmptyDB(name string) error {
	config, err := config.Get()
	if err != nil {
		return err
	}

	for _, c := range []string{
		fmt.Sprintf("DROP DATABASE IF EXISTS %s", name),
		fmt.Sprintf("CREATE DATABASE %s", name),
	} {
		cmd := exec.Command(
			"bash",
			"-c",
			fmt.Sprintf(
				"psql %s -c '%s'",
				config.PostgresServerConnectionString,
				c,
			),
		)

		var stdout, stderr strings.Builder
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("1>", stdout.String())
			fmt.Println("2>", stderr.String())
			return err
		}
	}

	return nil
}

func RunMigrations(ctx context.Context, conn *pgxpool.Conn) error {
	bytes, err := oos.ReadFileAtRoot("schema.sql")
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, string(bytes))
	return err
}
