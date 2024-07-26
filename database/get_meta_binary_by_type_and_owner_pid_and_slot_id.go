package database

import (
	"database/sql"
	"time"

	nex_types "github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/globals"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/types"
	"github.com/lib/pq"
)

func GetMetaBinaryByTypeAndOwnerPIDAndSlotID(dataType uint16, pid uint32, slotID uint8) *types.MetaBinary {
	metaBinary := types.NewMetaBinary()

	metaBinary.CreationTime = nex_types.NewDateTime(0)
	metaBinary.UpdatedTime = nex_types.NewDateTime(0)
	metaBinary.ReferredTime = nex_types.NewDateTime(0)
	metaBinary.ExpireTime = nex_types.NewDateTime(0)

	var creationTimestamp int64
	var updatedTimestamp int64
	var referredTimestamp int64
	var expireTimestamp int64

	err := Postgres.QueryRow(`
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
	FROM meta_binaries WHERE data_type=$1 AND owner_pid=$2 AND persistence_slot_id=$3`, dataType, pid, slotID).Scan(
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

	return metaBinary
}
