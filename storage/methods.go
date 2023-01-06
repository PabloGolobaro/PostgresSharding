package storage

func (s *Storage) shardByItemID(itemID int64) ShardNum {
	return ShardNum(itemID % bucketQuantity)
}
