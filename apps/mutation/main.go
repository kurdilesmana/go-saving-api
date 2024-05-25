package main

import (
	"github.com/kurdilesmana/go-saving-api/apps/mutation/cmd"
	"github.com/kurdilesmana/go-saving-api/apps/mutation/deps"
)

func main() {
	// init API mutation
	dependency := deps.SetupDependencies()
	cmd.ExecuteConsumeStream(dependency)
}
