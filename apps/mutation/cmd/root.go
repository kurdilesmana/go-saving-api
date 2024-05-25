package cmd

import (
	"github.com/kurdilesmana/go-saving-api/apps/mutation/deps"
	"github.com/kurdilesmana/go-saving-api/apps/mutation/server"
)

func ExecuteConsumeStream(dependency deps.Dependency) {
	// Start Redis server
	handler := server.SetupHandler(dependency)
	server.Redis(handler, dependency.Redis, dependency.Logger)
}
