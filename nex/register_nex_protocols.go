package nex

import (
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/globals"
	nex_datastore "github.com/PretendoNetwork/super-mario-3d-world-secure/nex/datastore"
)

func registerNEXProtocols() {
	datastoreProtocol := datastore.NewDataStoreProtocol(globals.NEXServer)

	datastoreProtocol.SearchObject(nex_datastore.SearchObject)
	datastoreProtocol.PreparePostObject(nex_datastore.PreparePostObject)
	datastoreProtocol.PrepareGetObject(nex_datastore.PrepareGetObject)
	datastoreProtocol.CompletePostObject(nex_datastore.CompletePostObject)
}
