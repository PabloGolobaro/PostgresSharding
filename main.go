package main

import (
	"PostgresSharding/config"
	"PostgresSharding/storage"
	"context"
	"log"
)

var (
	stor *storage.Storage
)

func main() {

}

func init() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Can't get config %s", err)
	}
	ctx := context.Background()
	dsns := make(map[storage.ShardNum]string, 0)
	for i, dsn := range conf.DSNS {
		dsns[storage.ShardNum(i)] = dsn
	}
	stor = storage.NewStorage(ctx, dsns)
	storage.CreateTablesForShards(stor)
	defer func(stor *storage.Storage) {
		err := stor.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stor)

}
