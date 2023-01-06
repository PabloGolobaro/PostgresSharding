package storage

import (
	"github.com/jmoiron/sqlx"
)

var schema string = `
create table if not exists items (
    id serial primary key,
    code text
);
`

func CreateTablesForShards(storage *Storage) {
	for _, shard := range storage.shardMap {
		createShardTable(shard)
	}
}

func createShardTable(shard *sqlx.DB) {
	shard.MustExec(schema)
}
