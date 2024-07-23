package nex

import (
	nex_secure_connection_common "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/globals"
)

func registerCommonProtocols() {
	secureProtocol := secure.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	nex_secure_connection_common.NewCommonProtocol(secureProtocol)

	// TODO: DataStore protocol should use the common protocol
}
