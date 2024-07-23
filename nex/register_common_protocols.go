package nex

import (
	"os"

	nex_datastore_common "github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
	nex_secure_connection_common "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/globals"
)

func registerCommonProtocols() {
	secureProtocol := secure.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	nex_secure_connection_common.NewCommonProtocol(secureProtocol)

	// TODO: DataStore protocol should use the common protocol
	datastoreProtocol := datastore.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(datastoreProtocol)
	datastore_protocol := nex_datastore_common.NewCommonProtocol(datastoreProtocol)
	datastore_protocol.S3Bucket = os.Getenv("PN_SM3DW_CONFIG_S3_BUCKET")
}
