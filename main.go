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

	return container

}

func main() {
    container := BuildContainer()
    
    err := container.Invoke(func(server *Server) {
		server.Run()
	})

	if err != nil {
		panic(err)
	}
}