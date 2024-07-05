package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"

	"github.com/PretendoNetwork/super-mario-3d-world-secure/globals"
)

var Postgres *sql.DB

func connectPostgres() {
	var err error

	Postgres, err = sql.Open("postgres", os.Getenv("PN_SM3DW_POSTGRES_URI"))
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	globals.Logger.Success("Connected to Postgres!")

	initPostgres()
}
