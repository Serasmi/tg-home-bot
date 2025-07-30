package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"

	"tg-home-bot/internal/config"
	"tg-home-bot/internal/echo"
	"tg-home-bot/internal/middleware"
	"tg-home-bot/internal/sensor"
	ha "tg-home-bot/pkg/home-assistant"

	"golang.org/x/sync/errgroup"
	tele "gopkg.in/telebot.v3"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var (
		app = &application{}
		err error
	)

	app.cfg, err = config.Init()
	if err != nil {
		slog.Error("init config", "error", err)
		return
	}

	initLogger(app.cfg.App.LogLevel)

	slog.Info("conf: permitted users", "users", app.cfg.Telegram.PermitUsers)

	haProvider := ha.NewService(app.cfg.HomeAssistant.URL, app.cfg.HomeAssistant.Token)

	_, err = app.initBot(ctx, haProvider)
	if err != nil {
		slog.Error("init bot", "error", err)
		return
	}

	slog.Info("app starting...")

	group, ctx := errgroup.WithContext(ctx)
	group.Go(
		func() error {
			return app.run(ctx)
		},
	)
	group.Go(func() error {
		<-ctx.Done()
		stop()

		slog.Info("interrupted. app is shutting down...")

		doneCtx, cancel := context.WithTimeout(context.Background(), app.cfg.App.ShutdownTimeout)
		defer cancel()

		closeErr := app.close(doneCtx)
		if closeErr != nil {
			slog.Error("close application", "error", closeErr)
		} else {
			slog.Info("close application")
		}

		return nil
	})

	err = group.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("wait application", "error", err)
	}

	slog.Info("application stopped")
}

func initLogger(level string) {
	var lvl slog.Level

	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		slog.Default().Error("unmarshal log level", "level", level, "error", err)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))

	slog.SetDefault(log)
}

type application struct {
	cfg *config.Config
	bot *tele.Bot
}

func (app *application) run(_ context.Context) error {
	app.bot.Start()

	return nil
}

func (app *application) close(_ context.Context) error {
	app.bot.Stop()

	return nil
}

func (app *application) initBot(ctx context.Context, haProvider ha.Service) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  os.Getenv("TG_API_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		OnError: func(err error, c tele.Context) {
			if c != nil {
				slog.Error("bot error", "update_id", c.Update().ID, "error", err.Error())
				return
			}

			slog.Error("bot error", "error", err.Error())
		},
	}

	var err error

	app.bot, err = tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("new bot: %w", err)
	}

	app.bot.Use(middleware.PermitUsers(app.cfg.Telegram.PermitUsers))
	echo.RegisterHandler(app.bot)

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, err
	}

	sensor.RegisterHandler(ctx, app.bot, sensor.NewService(haProvider, loc, app.cfg.HomeAssistant.Timeout))

	return app.bot, nil
}
