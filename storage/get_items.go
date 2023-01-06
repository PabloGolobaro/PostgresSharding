package storage

import (
	"PostgresSharding/models"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
	"sync"
)

func (s *Storage) GetItems(ctx context.Context, itemIDs ...int64) ([]models.Item, error) {
	shardToItems := s.sortItemsIDsByShard(itemIDs...)

	respChan := make(chan []models.Item, len(shardToItems))
	errChan := make(chan error, len(shardToItems))
	wg := &sync.WaitGroup{}

	for shardID, ids := range shardToItems {
		wg.Add(1)
		shard := s.shardMap[shardID]
		go s.asyncGetItemsByID(ctx, shard, ids, wg, respChan, errChan)
	}

	wg.Wait()
	close(respChan)
	close(errChan)

	result := make([]models.Item, 0)
	for items := range respChan {
		result = append(result, items...)
	}

	errs := make([]error, 0, len(errChan))
	for e := range errChan {
		errs = append(errs, e)
	}
	err := multierr.Combine(errs...)

	return result, err
}

func (s *Storage) sortItemsIDsByShard(itemIDs ...int64) map[ShardNum][]int64 {
	shardToItems := make(map[ShardNum][]int64)

	for _, id := range itemIDs {
		shardID := s.shardByItemID(id)
		if _, ok := shardToItems[shardID]; !ok {
			shardToItems[shardID] = make([]int64, 0)
		}

		shardToItems[shardID] = append(shardToItems[shardID], id)
	}

	return shardToItems
}

func (s *Storage) asyncGetItemsByID(ctx context.Context, shard *sqlx.DB, itemsIDs []int64, wg *sync.WaitGroup, resp chan<- []models.Item, errs chan<- error) {
	defer wg.Done()
	items, err := s.getItemsByID(ctx, shard, itemsIDs)
	if err != nil {
		errs <- fmt.Errorf("%s [getItemsByID] can't select from shard %v", err, shard)
	}
	resp <- items
}

func (s *Storage) getItemsByID(ctx context.Context, shard *sqlx.DB, itemsIDs []int64) ([]models.Item, error) {
	items := make([]models.Item, 0)
	//goland:noinspection ALL
	query, args, err := sqlx.In("SELECT * FROM items WHERE id IN (?);", itemsIDs)
	query = shard.Rebind(query)
	err = shard.SelectContext(ctx, &items, query, args...)
	return items, err
}
