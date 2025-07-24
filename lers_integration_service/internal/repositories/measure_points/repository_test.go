package measure_points_test

import (
	"context"
	"fmt"
	"lers_integration_service/internal/app"
	"lers_integration_service/internal/repositories/measure_points"
	"lers_integration_service/tests"
	"log"
	"os"
	"testing"
	"time"
)

var repositoryTest measure_points.Repository
var appInstance *app.App

func TestMain(m *testing.M) {
	fmt.Println("Start tests...")
	ctx := context.Background()

	tests.Init(ctx)

	time.Sleep(5 * time.Second)

	app, cleanup, err := tests.GetApp()
	if err != nil {
		log.Fatalf("Не удалось получить экземпляр приложения: %s", err)
	}
	appInstance = app
	repositoryTest = app.S().MeasurePointsRepository()

	exitCode := m.Run()

	cleanup()

	os.Exit(exitCode)
}
