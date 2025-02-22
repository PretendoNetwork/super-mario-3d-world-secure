package nex

import (
	"fmt"
	"os"
	"strconv"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/super-mario-3d-world/globals"
)

var serverBuildString string

func StartAuthenticationServer() {
	serverBuildString = "branch:ngs_2_30 build:2_22_11148_30_2"

	globals.AuthenticationServer = nex.NewPRUDPServer()

	globals.AuthenticationEndpoint = nex.NewPRUDPEndPoint(1)
	globals.AuthenticationEndpoint.ServerAccount = globals.AuthenticationServerAccount
	globals.AuthenticationEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.AuthenticationEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.AuthenticationServer.BindPRUDPEndPoint(globals.AuthenticationEndpoint)
	globals.AuthenticationServer.ByteStreamSettings.UseStructureHeader = false

	globals.AuthenticationServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(3, 4, 0))
	globals.AuthenticationServer.AccessKey = "86a89872"

	globals.AuthenticationEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("=== SM3DW - Auth ===")
		fmt.Printf("Protocol ID: %d\n", request.ProtocolID)
		fmt.Printf("Method ID: %d\n", request.MethodID)
		fmt.Println("==================")
	})

	registerCommonAuthenticationServerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_SM3DW_AUTHENTICATION_SERVER_PORT"))

	globals.AuthenticationServer.Listen(port)
}
