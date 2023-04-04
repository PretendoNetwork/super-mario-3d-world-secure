package nex_datastore

import (
	"fmt"
	"os"
	"time"

	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/database"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/globals"
)

func PreparePostObject(err error, client *nex.Client, callID uint32, param *nexproto.DataStorePreparePostParam) {
	metaBinary := database.GetMetaBinaryByTypeAndOwnerPIDAndSlotID(param.DataType, client.PID(), uint8(param.PersistenceInitParam.PersistenceSlotId))

	if metaBinary.DataID != 0 {
		// * Meta binary already exists
		if param.PersistenceInitParam.DeleteLastObject {
			// * Delete existing object before uploading new one
			// TODO - Check error
			_ = database.DeleteMetaBinaryByDataID(metaBinary.DataID)
			// TODO - Delete old ratings?
		}
	}

	// TODO - See if this is actually always the case?
	// * Always upload a new object (?)
	dataID := database.InsertMetaBinaryByDataStorePreparePostParamWithOwnerPID(param, client.PID())

	for i := 0; i < len(param.RatingInitParams); i++ {
		ratingInitParam := param.RatingInitParams[i]

		// TODO - Check error
		_ = database.InsertRatingByDataIDAndDataStoreRatingInitParamWithSlot(dataID, ratingInitParam)
	}

	s3Bucket := os.Getenv("PN_SM3DW_CONFIG_S3_BUCKET")
	key := fmt.Sprintf("ghosts/%d.bin", dataID)

	input := &globals.PostObjectInput{
		Bucket:    s3Bucket,
		Key:       key,
		ExpiresIn: time.Minute * 15,
	}

	res, _ := globals.S3PresignClient.PresignPostObject(input)

	fieldKey := nexproto.NewDataStoreKeyValue()
	fieldKey.Key = "key"
	fieldKey.Value = key

	fieldCredential := nexproto.NewDataStoreKeyValue()
	fieldCredential.Key = "X-Amz-Credential"
	fieldCredential.Value = res.Credential

	fieldSecurityToken := nexproto.NewDataStoreKeyValue()
	fieldSecurityToken.Key = "X-Amz-Security-Token"
	fieldSecurityToken.Value = ""

	fieldAlgorithm := nexproto.NewDataStoreKeyValue()
	fieldAlgorithm.Key = "X-Amz-Algorithm"
	fieldAlgorithm.Value = "AWS4-HMAC-SHA256"

	fieldDate := nexproto.NewDataStoreKeyValue()
	fieldDate.Key = "X-Amz-Date"
	fieldDate.Value = res.Date

	fieldPolicy := nexproto.NewDataStoreKeyValue()
	fieldPolicy.Key = "policy"
	fieldPolicy.Value = res.Policy

	fieldSignature := nexproto.NewDataStoreKeyValue()
	fieldSignature.Key = "X-Amz-Signature"
	fieldSignature.Value = res.Signature

	pReqPostInfo := nexproto.NewDataStoreReqPostInfo()

	pReqPostInfo.DataID = 1
	pReqPostInfo.URL = res.URL
	pReqPostInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
	pReqPostInfo.FormFields = []*nexproto.DataStoreKeyValue{
		fieldKey,
		fieldCredential,
		fieldSecurityToken,
		fieldAlgorithm,
		fieldDate,
		fieldPolicy,
		fieldSignature,
	}
	pReqPostInfo.RootCACert = []byte{}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteStructure(pReqPostInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodPreparePostObject, rmcResponseBody)

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

	/*
		rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
		rmcResponse.SetError(nex.Errors.Core.NotImplemented)

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
	*/
}
