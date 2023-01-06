package storage

import (
	"PostgresSharding/config"
	"PostgresSharding/models"
	"context"
	"github.com/thanhpk/randstr"
	"log"
	"testing"
)

func TestStorage_AddItems(t *testing.T) {
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

	items := []models.Item{}
	for i := 0; i < 10; i++ {
		s := randstr.String(6)
		items = append(items, models.Item{
			ID:   int64(i),
			Code: s,
		})
	}

	type args struct {
		ctx   context.Context
		items []models.Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "first", args: args{ctx: ctx, items: nil}, wantErr: false},
		{name: "second", args: args{ctx: ctx, items: items}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := stor.AddItems(tt.args.ctx, tt.args.items...); (err != nil) != tt.wantErr {
				t.Errorf("AddItems() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
