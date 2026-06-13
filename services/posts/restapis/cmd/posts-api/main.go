package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	hatchetsdk "github.com/hatchet-dev/hatchet/sdks/go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"

	infrastructure_applications "github.com/steve-rodrigue/aabs/services/posts/restapis/infrastructure/applications"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/servers"
)

func main() {
	root := &cobra.Command{
		Use: "posts-api",
	}

	root.AddCommand(
		installCommand(),
		startCommand(),
		stopCommand(),
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func installCommand() *cobra.Command {
	return &cobra.Command{
		Use: "install",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			pool, err := openPool(ctx)
			if err != nil {
				return err
			}
			defer pool.Close()

			return infrastructure_applications.Install(ctx, pool)
		},
	}
}

func startCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "start",
		Aliases: []string{"serve"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			pool, err := openPool(ctx)
			if err != nil {
				return err
			}
			defer pool.Close()

			if envBool("POSTS_API_INSTALL_ON_START", true) {
				if err := infrastructure_applications.Install(ctx, pool); err != nil {
					return err
				}
			}

			hatchetClient, err := newHatchetClient()
			if err != nil {
				return err
			}

			application := infrastructure_applications.New(
				pool,
				hatchetClient,
			)

			server := &http.Server{
				Addr:              envString("POSTS_API_ADDR", ":8200"),
				Handler:           servers.New(application),
				ReadHeaderTimeout: 5 * time.Second,
			}

			if err := writePID(); err != nil {
				return err
			}
			defer removePID()

			errs := make(chan error, 1)

			go func() {
				errs <- server.ListenAndServe()
			}()

			shutdown := make(chan os.Signal, 1)
			signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

			select {
			case signal := <-shutdown:
				fmt.Printf("received %s, shutting down\n", signal)

				shutdownCtx, cancel := context.WithTimeout(
					context.Background(),
					10*time.Second,
				)
				defer cancel()

				return server.Shutdown(shutdownCtx)

			case err := <-errs:
				if errors.Is(err, http.ErrServerClosed) {
					return nil
				}

				return err
			}
		},
	}
}

func stopCommand() *cobra.Command {
	return &cobra.Command{
		Use: "stop",
		RunE: func(cmd *cobra.Command, args []string) error {
			pidBytes, err := os.ReadFile(pidFile())
			if err != nil {
				return err
			}

			pid, err := strconv.Atoi(string(pidBytes))
			if err != nil {
				return err
			}

			process, err := os.FindProcess(pid)
			if err != nil {
				return err
			}

			return process.Signal(syscall.SIGTERM)
		},
	}
}

func openPool(
	ctx context.Context,
) (*pgxpool.Pool, error) {
	dsn := os.Getenv("POSTS_POSTGRES_DSN")
	if dsn == "" {
		dsn = os.Getenv("POSTS_POSTGRES_TEST_DSN")
	}

	if dsn == "" {
		return nil, errors.New("POSTS_POSTGRES_DSN is not set")
	}

	return pgxpool.New(ctx, dsn)
}

func newHatchetClient() (*hatchetsdk.Client, error) {
	if os.Getenv("HATCHET_CLIENT_TOKEN") == "" {
		return nil, nil
	}

	return hatchetsdk.NewClient()
}

func writePID() error {
	return os.WriteFile(
		pidFile(),
		[]byte(strconv.Itoa(os.Getpid())),
		0644,
	)
}

func removePID() {
	_ = os.Remove(pidFile())
}

func pidFile() string {
	return envString(
		"POSTS_API_PID_FILE",
		"/tmp/aabs-posts-api.pid",
	)
}

func envString(
	key string,
	defaultValue string,
) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func envBool(
	key string,
	defaultValue bool,
) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value == "true" ||
		value == "1" ||
		value == "yes"
}
