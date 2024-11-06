package main

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelFunc()

	router, err := message.NewRouter(message.RouterConfig{}, watermill.NewStdLogger(true, true))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting %s%s", serviceName, version)

	routes(router)

	if err = router.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
