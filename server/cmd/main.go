package main

import (
	"github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/server"
)

func main() {
	/**
	 * Connect to db
	 */
	defaultDB := db.Database{
		DataBaseRefName: "default",
		URL:             config.GetSecret("MONGO_URI"),
		DataBaseName:    config.GetSecret("MONGO_DATABASE"),
		RetryWrites:     config.GetSecret("MONGO_RETRYWRITES"),
	}
	defaultDB.Init()
	defer defaultDB.Disconnect()

	server.Start()
}
