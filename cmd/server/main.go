package main

import (
	"context"
	"sync"
	"time"

	"github.com/diontr00/serverlessgo/config"
)

var (
	app_    *config.Applications
	appOnce sync.Once
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	appOnce.Do(func() {
		app_ = config.NewApp(ctx)
	})

	app_.Start()
}
