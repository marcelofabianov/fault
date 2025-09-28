package main

import (
	"log/slog"
	"os"

	"github.com/marcelofabianov/fault"

	"github.com/marcelofabianov/example/service"
)

// handleError gera um único log de erro estruturado em JSON.
func handleError(err error) {
	if err == nil {
		return
	}

	fErr, ok := fault.AsFault(err)
	if !ok {
		slog.Error("Operação falhou com um erro inesperado", "error", err.Error())
		return
	}

	slog.Error(
		"Operação falhou",
		"message", fErr.Message,
		"code", fErr.Code,
		"context", fErr.Context,
	)
}

func main() {
	// Configura o logger slog para uma saída em JSON.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	userService := service.NewUserService()

	// --- Cenário 1: Sucesso (não produzirá nenhuma saída) ---
	_, err := userService.FindUserByID("0199910e-d028-7e10-aa14-3954e781a9bc")
	handleError(err)

	// --- Cenário 2: Not Found (produzirá uma saída de erro) ---
	_, err = userService.FindUserByID("0199910f-a9b8-7440-8438-2c49c7e0f21c")
	handleError(err)

	// --- Cenário 3: Inativo (produzirá uma saída de erro) ---
	_, err = userService.FindUserByID("0199910f-4238-7882-a075-b7842b4d41f0")
	handleError(err)

	// --- Cenário 4: Conflito (produzirá uma saída de erro) ---
	err = userService.RegisterUser("0199910e-d028-7e10-aa14-3954e781a9bc", "Jane Doe", "jane.doe@example.com")
	handleError(err)
}
