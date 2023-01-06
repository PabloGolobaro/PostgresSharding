package storage

import (
	"PostgresSharding/models"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
	"sync"
)

func (s *Storage) AddItems(ctx context.Context, items ...models.Item) error {
	itemsByShardMap := s.itemsByShard(items...)
	errChan := make(chan error, len(itemsByShardMap))
	wg := &sync.WaitGroup{}

	for shardID, items := range itemsByShardMap {
		wg.Add(1)
		shard := s.shardMap[shardID]
		go s.asyncAddItems(ctx, errChan, wg, shard, items...)
	}

	wg.Wait()
	close(errChan)

	errs := make([]error, 0, len(errChan))
	for e := range errChan {
		errs = append(errs, e)
	}

	return multierr.Combine(errs...)
}

func (s *Storage) itemsByShard(items ...models.Item) map[ShardNum][]models.Item {
	itemsByShard := make(map[ShardNum][]models.Item)
	for _, item := range items {
		shardID := s.shardByItemID(item.ID)
		if _, ok := itemsByShard[shardID]; !ok {
			itemsByShard[shardID] = make([]models.Item, 0)
		}

		itemsByShard[shardID] = append(itemsByShard[shardID], item)
	}

	return itemsByShard
}

func (s *Storage) asyncAddItems(ctx context.Context, errChan chan<- error, wg *sync.WaitGroup, shard *sqlx.DB, items ...models.Item) {
	defer wg.Done()
	err := s.addItems(ctx, shard, items...)
	if err != nil {
		errChan <- fmt.Errorf("[asyncAddItems] can't insert to shard %s", err)
	}

}

func (s *Storage) addItems(ctx context.Context, shard *sqlx.DB, items ...models.Item) error {
	//goland:noinspection SqlResolve
	query, args, err := sqlx.Named(`INSERT INTO items (id, code)
        VALUES (:id, :code)`, items)
	if err != nil {
		return err
	}
	query = shard.Rebind(query)
	_, err = shard.DB.ExecContext(ctx, query, args...)
	return err
}
