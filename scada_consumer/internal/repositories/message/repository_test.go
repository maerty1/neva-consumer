package message_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"scada_consumer/internal/app"
	"scada_consumer/internal/repositories/message"
	"scada_consumer/tests"
	"testing"
)

var repositoryTest message.Repository
var appInstance *app.App

func TestMain(m *testing.M) {
	fmt.Println("Start tests...")
	ctx := context.Background()

	tests.Init(ctx)

	app, cleanup, err := tests.GetApp()
	if err != nil {
		log.Fatalf("Ошибка получения app instance: %s", err)
	}

	appInstance = app
	repositoryTest = app.S().MessageRepository()

	exitCode := m.Run()

	cleanup()

	os.Exit(exitCode)
}
