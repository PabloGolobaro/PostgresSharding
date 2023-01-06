package storage

import (
	"PostgresSharding/config"
	"context"
	_ "github.com/lib/pq"
	"log"
	"testing"
)

func TestCreateTablesForShards(t *testing.T) {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Can't get config %s", err)
	}
	ctx := context.Background()
	dsns := make(map[ShardNum]string, 0)
	for i, dsn := range conf.DSNS {
		dsns[ShardNum(i)] = dsn
	}
	stor := NewStorage(ctx, dsns)
	tests := []struct {
		name    string
		storage *Storage
	}{
		{name: "Main", storage: stor},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateTablesForShards(tt.storage)
		})
	}
}
