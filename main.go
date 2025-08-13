package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/robertobff/nexpos/adapter"
	"github.com/robertobff/nexpos/application"
	"go.uber.org/fx"
)

func main() {
	if os.Getenv("ENV") != "production" {
		LoadConfig()
	}

	fx.New(
		adapter.Module,
		application.Module,
	).Run()
}

func LoadConfig() {
	_, b, _, _ := runtime.Caller(0)

	basepath := filepath.Dir(b)

	err := godotenv.Load(fmt.Sprintf("%v/.env", basepath))
	if err != nil {
		panic(err)
	}
}
