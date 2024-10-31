package main

// Think overcomplicated? I agree.

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func getLogLevel() slog.Level { // is this necessary? maybe not. But it's nice to have one.
	lvl := os.Getenv("LOG_LEVEL")
	switch strings.ToUpper(lvl) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func migrate(db *sql.DB) error {
	slog.Info("Starting database migration...")
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS results (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp DATETIME,
        r1 TEXT, r2 TEXT, r3 TEXT, r4 TEXT, r5 TEXT,
        r6 TEXT, r7 TEXT, r8 TEXT, r9 TEXT, r10 TEXT
    )`)
	if err != nil {
		slog.Error("Migration failed", "error", err)
		return err
	}
	slog.Info("Migration completed successfully")
	return nil
}

func crawl(db *sql.DB) error {
	slog.Debug("Starting archiving...")

	resp, err := http.Get("https://search.namu.wiki/api/ranking")
	if err != nil {
		slog.Error("Request failed", "error", err)
		return err
	}
	defer resp.Body.Close()

	var items []string
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		slog.Error("JSON decode failed", "error", err)
		return err
	}
	slog.Debug("API access succuss", "resp_code", resp.StatusCode, "len", len(items))
	if len(items) != 10 {
		slog.Warn("Item len not 10. Perhapse the API has changed?", "len", len(items))
	}

	slog.Debug("Database insert started", "data", items)
	_, err = db.Exec(`INSERT INTO results 
        (timestamp, r1, r2, r3, r4, r5, r6, r7, r8, r9, r10) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		time.Now(), items[0], items[1], items[2], items[3], items[4],
		items[5], items[6], items[7], items[8], items[9])
	if err != nil {
		slog.Error("Database insert failed", "error", err)
		return err
	}
	slog.Info("Data saved successfully")
	return nil
}

func main() {
	slogTextHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLogLevel(),
	})
	logger := slog.New(slogTextHandler)
	slog.SetDefault(logger)

	app := &cli.App{
		Name:  "namu-rank-archive",
		Usage: "Archives Namu wiki's search rank",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "db",
				Aliases: []string{"d"},
				Usage:   "Use `FILE` as a database or newly create one",
				Value:   "./data.db",
				EnvVars: []string{"NAMU_RANK_DB"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "migrate",
				Usage: "Create database tables",
				Action: func(c *cli.Context) error {
					slog.Debug("Opening database connection...")
					dbName := c.String("db")
					if dbName == "" {
						slog.Error("Database is not specified!")
						os.Exit(1)
					}
					db, err := sql.Open("sqlite3", dbName)
					if err != nil {
						slog.Error("Database connection failed", "error", err)
						return err
					}
					defer db.Close()
					return migrate(db)
				},
			},
			{Name: "archive", Usage: "Fetch and save search rank", Action: func(c *cli.Context) error {
				dbName := c.String("db")
				if _, err := os.Stat(dbName); err != nil && errors.Is(err, os.ErrNotExist) {
					slog.Error("Database file does not exist. Run migration first.", "db", dbName)
					os.Exit(1)
				} else if err != nil {
					slog.Error("Something went wrong during database checking", "error", err)
					os.Exit(1)
				}
				slog.Debug("Opening database connection...", "database", dbName)
				if dbName == "" {
					slog.Error("Database is not specified!")
					os.Exit(1)
				}
				db, err := sql.Open("sqlite3", dbName)
				if err != nil {
					slog.Error("Database connection failed", "error", err)
					return err
				}
				defer db.Close()
				slog.Debug("Database connection succeeded")
				return crawl(db)
			}},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("Something went wrong during execution", "error", err)
		os.Exit(1)
	}
}
