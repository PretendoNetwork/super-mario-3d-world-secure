package nex

import (
	"fmt"
	"os"
	"strconv"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/super-mario-3d-world/globals"
)

func StartSecureServer() {
	globals.SecureServer = nex.NewPRUDPServer()

	globals.SecureEndpoint = nex.NewPRUDPEndPoint(1)
	globals.SecureEndpoint.IsSecureEndPoint = true
	globals.SecureEndpoint.ServerAccount = globals.SecureServerAccount
	globals.SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.SecureServer.BindPRUDPEndPoint(globals.SecureEndpoint)
	globals.SecureServer.ByteStreamSettings.UseStructureHeader = false

	globals.SecureServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(3, 4, 0))
	globals.SecureServer.AccessKey = "86a89872"

	globals.SecureEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("=== SM3DW - Secure ===")
		fmt.Printf("Protocol ID: %d\n", request.ProtocolID)
		fmt.Printf("Method ID: %d\n", request.MethodID)
		fmt.Println("====================")
	})

	registerCommonSecureServerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_SM3DW_SECURE_SERVER_PORT"))

	globals.SecureServer.Listen(port)
}
