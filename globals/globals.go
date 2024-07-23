package globals

import (
	pb "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var Logger = plogger.NewLogger()
var KerberosPassword = "password" // * Default password
var SecureServer *nex.PRUDPServer
var SecureEndpoint *nex.PRUDPEndPoint

var GRPCAccountClientConnection *grpc.ClientConn
var GRPCAccountClient pb.AccountClient
var GRPCAccountCommonMetadata metadata.MD

var S3Client *s3.Client
var S3PresignClient *PresignClient
