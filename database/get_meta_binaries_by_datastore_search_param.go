package database

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/globals"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/types"
	"github.com/lib/pq"
)

/*
DataStoreSearchParam: &{SearchTarget:1 OwnerIds:[] OwnerType:1 DestinationIds:[] DataType:209 CreatedAfter:0xc0004ac060 CreatedBefore:0xc0004ac070 UpdatedAfter:0xc0004ac080 UpdatedBefore:0xc0004ac090 ReferDataId:0 Tags:[EnterCatMario] ResultOrderColumn:0 ResultOrder:1 ResultRange:0xc0001d0018 ResultOption:5 MinimalRatingFrequency:0 UseCache:false Structure:{StructureInterface:<nil>}}
DataStoreSearchParam.ResultRange: &{Offset:0 Length:20 Structure:{StructureInterface:<nil>}}

*/

func GetMetaBinariesByDataStoreSearchParam(param *nexproto.DataStoreSearchParam) []*types.MetaBinary {
	metaBinaries := make([]*types.MetaBinary, 0, param.ResultRange.Length)

	rows, err := Postgres.Query(`
		SELECT
		data_id,
		owner_pid,
		name,
		data_type,
		meta_binary,
		permission,
		del_permission,
		flag,
		period,
		tags,
		persistence_slot_id,
		extra_data,
		creation_time,
		updated_time,
		referred_time,
		expire_time
		FROM meta_binaries WHERE data_type=$1 AND tags && $2 LIMIT $3`,
		param.DataType,
		pq.Array(param.Tags),
		param.ResultRange.Length,
	)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return metaBinaries
	}

	for rows.Next() {
		metaBinary := types.NewMetaBinary()

		metaBinary.CreationTime = nex.NewDateTime(0)
		metaBinary.UpdatedTime = nex.NewDateTime(0)
		metaBinary.ReferredTime = nex.NewDateTime(0)
		metaBinary.ExpireTime = nex.NewDateTime(0)

		var creationTimestamp int64
		var updatedTimestamp int64
		var referredTimestamp int64
		var expireTimestamp int64

		err := rows.Scan(
			&metaBinary.DataID,
			&metaBinary.OwnerPID,
			&metaBinary.Name,
			&metaBinary.DataType,
			&metaBinary.Buffer,
			&metaBinary.Permission,
			&metaBinary.DeletePermission,
			&metaBinary.Flag,
			&metaBinary.Period,
			pq.Array(&metaBinary.Tags),
			&metaBinary.PersistenceSlotID,
			pq.Array(&metaBinary.ExtraData),
			&creationTimestamp,
			&updatedTimestamp,
			&referredTimestamp,
			&expireTimestamp,
		)

		if err != nil && err != sql.ErrNoRows {
			globals.Logger.Critical(err.Error())
		}

		if err == nil {
			_ = metaBinary.CreationTime.FromTimestamp(time.Unix(0, creationTimestamp))
			_ = metaBinary.UpdatedTime.FromTimestamp(time.Unix(0, updatedTimestamp))
			_ = metaBinary.ReferredTime.FromTimestamp(time.Unix(0, referredTimestamp))
			_ = metaBinary.ExpireTime.FromTimestamp(time.Unix(0, expireTimestamp))
		}

		metaBinaries = append(metaBinaries, metaBinary)
	}

	return metaBinaries
}
