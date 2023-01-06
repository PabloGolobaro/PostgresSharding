package storage

import (
	"PostgresSharding/config"
	"PostgresSharding/models"
	"context"
	"log"
	"reflect"
	"testing"
)

func TestStorage_GetItems(t *testing.T) {
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
	type args struct {
		ctx     context.Context
		itemIDs []int64
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Item
		wantErr bool
	}{
		{name: "first", args: args{ctx: ctx, itemIDs: []int64{0, 1}}, want: []models.Item{{
			0,
			"1bhdY8",
		}, {
			1,
			"nfzIx4",
		}}, wantErr: false},
		{name: "second", args: args{ctx: ctx, itemIDs: []int64{7, 8}}, want: []models.Item{{
			7,
			"3ruhRm",
		}, {
			8,
			"f27hKo",
		}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := stor.GetItems(tt.args.ctx, tt.args.itemIDs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("GetItems() got = %v, want %v", got, tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItems() got = %v, want %v", got, tt.want)
			}
		})
	}
}
