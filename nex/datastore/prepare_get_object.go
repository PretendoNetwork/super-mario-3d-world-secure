package nex_datastore

import (
	"fmt"
	"os"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/database"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/globals"
)

func PrepareGetObject(_ error, client *nex.Client, callID uint32, param *datastore.DataStorePrepareGetParam) {
	// TODO - Check error
	metaBinary, _ := database.GetMetaBinaryByDataID(uint32(param.DataID))

	bucket := os.Getenv("PN_SM3DW_CONFIG_S3_BUCKET")
	key := fmt.Sprintf("ghosts/%d.bin", metaBinary.DataID)

	objectSize, _ := globals.S3ObjectSize(bucket, key)

	pReqGetInfo := datastore.NewDataStoreReqGetInfo()

	pReqGetInfo.URL = fmt.Sprintf("https://%s.b-cdn.net/%s", bucket, key)
	pReqGetInfo.RequestHeaders = []*datastore.DataStoreKeyValue{}
	pReqGetInfo.Size = uint32(objectSize)
	pReqGetInfo.RootCA = []byte{}
	pReqGetInfo.DataID = uint64(metaBinary.DataID)

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteStructure(pReqGetInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodPrepareGetObject, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
