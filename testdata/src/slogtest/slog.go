package slogtest

import "log/slog"

func test() {
	slog.Info("Starting server")    // want "lowercase"
	slog.Info("запуск сервера")     // want "English"
	slog.Info("server started!")    // want "special"
	slog.Info("user password 1234") // want "sensitive"
	slog.Info("server started")
}
