package main

import (
	"github.com/kurdilesmana/go-saving-api/apps/account/cmd"
	"github.com/kurdilesmana/go-saving-api/apps/account/deps"
)

func main() {
	// init API Account
	dependency := deps.SetupDependencies()
	cmd.ExecuteHTTPAccount(dependency)
}
