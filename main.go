package main

import (
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	container.Provide(NewPgConfig)
	container.Provide(NewConfig)
	container.Provide(PostgresSetup)
	container.Provide(NewShapeRepository)
	container.Provide(NewShapeService)
	container.Provide(NewServer)
	container.Provide(NewHub)
	container.Provide(NewOperationManager)

	return container

}

func main() {
    container := BuildContainer()
    
    err := container.Invoke(func(server *Server, hub *Hub) {
		go hub.run()
		server.Run()
	})

	if err != nil {
		panic(err)
	}
}