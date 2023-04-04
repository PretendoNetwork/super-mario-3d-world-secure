package database

func ConnectAll() {
	// * Future proofing for when/if we need Mongo
	connectPostgres()
}
