package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/briandowns/openweathermap"
	"github.com/intelux/hubitat"
	"github.com/spf13/cobra"
)

var (
	rootCtx, rootCtxCancel = withInterrupt(context.Background())
	endpoint               = ":8081"
)

var rootCmd = &cobra.Command{
	Use:  "hubitat-prometheus-exporter",
	Args: cobra.NoArgs,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		rootCtxCancel()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := hubitat.NewClientFromEnv()

		if err != nil {
			return fmt.Errorf("instanciating Hubitat client: %s", err)
		}

		fmt.Fprintf(cmd.OutOrStderr(), "Listening on %s.\n", endpoint)

		var currentWeatherData *openweathermap.CurrentWeatherData

		owmAPIKey := os.Getenv("OWM_API_KEY")
		owmCityID, _ := strconv.Atoi(os.Getenv("OWM_CITY_ID"))

		if owmAPIKey != "" && owmCityID != 0 {
			if currentWeatherData, err = openweathermap.NewCurrent("C", "EN", owmAPIKey); err != nil {
				return fmt.Errorf("failed to initialize Open Weather Map API: %s", err)
			}
		}

		server := &http.Server{
			Addr: endpoint,
			Handler: &Handler{
				Client:             client,
				CurrentWeatherData: currentWeatherData,
				CityID:             owmCityID,
			},
		}

		go func() {
			<-rootCtx.Done()

			shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			server.Shutdown(shutdownCtx)
			cancel()

			server.Close()
		}()

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}

		fmt.Fprintf(cmd.OutOrStderr(), "Server shutdown.\n")

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&endpoint, "endpoint", "l", endpoint, "the endpoint to listen on")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func withInterrupt(ctx context.Context) (context.Context, func()) {
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		cancel()
	}()

	go func() {
		<-ctx.Done()

		signal.Stop(ch)
		close(ch)
	}()

	return ctx, cancel
}
