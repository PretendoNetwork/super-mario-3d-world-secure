package nex_datastore

import (
	"errors"
	"fmt"
	"os"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/database"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/globals"
)

func CompletePostObject(_ error, client *nex.Client, callID uint32, param *datastore.DataStoreCompletePostParam) {
	bucket := os.Getenv("PN_SM3DW_CONFIG_S3_BUCKET")
	key := fmt.Sprintf("ghosts/%d.bin", param.DataID)

	_, err := globals.S3HeadRequest(bucket, key)

	// * DELETE FAILED UPLOAD
	// TODO - LOG THIS AND FIX THE ISSUE WITH UPLOADS FAILING
	if err != nil {
		var re *awshttp.ResponseError
		if errors.As(err, &re) {
			if re.Response.StatusCode == 404 {
				_ = database.DeleteMetaBinaryByDataID(uint32(param.DataID))
			} else {
				globals.Logger.Error(err.Error())
			}
		} else {
			globals.Logger.Error(err.Error())
		}
	}

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodCompletePostObject, nil)

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
