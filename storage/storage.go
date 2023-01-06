package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
)

// Обозначим количество шардов.
const bucketQuantity = 2

// Для лучшей семантики.
type ShardNum int
type shardMap map[ShardNum]*sqlx.DB

type Storage struct {
	shardMap shardMap
}

func initShardMap(ctx context.Context, dsns map[ShardNum]string) shardMap {
	m := make(shardMap, len(dsns))
	for sh, dsn := range dsns {
		m[sh] = discoveryShard(ctx, dsn)
	}

	return m
}

func discoveryShard(ctx context.Context, dsn string) *sqlx.DB {
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		panic(err)
	}

	return db
}

func NewStorage(ctx context.Context, dsns map[ShardNum]string) *Storage {
	return &Storage{
		shardMap: initShardMap(ctx, dsns),
	}
}

func (s *Storage) Close() error {
	for _, db := range s.shardMap {
		return db.Close()
	}
	return nil
}
