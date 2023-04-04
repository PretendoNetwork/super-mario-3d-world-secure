package database

import "github.com/PretendoNetwork/super-mario-3d-world-secure/globals"

func initPostgres() {
	var err error

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS meta_binaries (
		data_id serial PRIMARY KEY,
		owner_pid integer,
		name text,
		data_type integer,
		meta_binary bytea,
		permission integer,
		del_permission integer,
		flag integer,
		period integer,
		tags text[],
		persistence_slot_id integer,
		extra_data text[],
		creation_time bigint,
		updated_time bigint,
		referred_time bigint,
		expire_time bigint
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS ratings (
		data_id integer,
		slot integer,
		flag integer,
		internal_flag integer,
		lock_type integer,
		initial_value bigint,
		range_min integer,
		range_max integer,
		period_hour integer,
		period_duration integer,
		total_value bigint,
		count integer,
		PRIMARY KEY(data_id, slot)
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("Postgres tables created")
}
