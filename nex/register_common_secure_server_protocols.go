package nex

import (
	//"context"

	"github.com/PretendoNetwork/nex-go/v2/types"
	common_secure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	local_globals "github.com/PretendoNetwork/super-mario-3d-world/globals"
)

func registerCommonSecureServerProtocols() {
	secureProtocol := secure.NewProtocol()
	local_globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	secure := common_secure.NewCommonProtocol(secureProtocol)
	secure.CreateReportDBRecord = func(pid types.PID, reportID types.UInt32, reportData types.QBuffer) error {
		return nil
	}

}
