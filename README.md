# PostgreSQL Sharding

This project is an **experiment** realization of sharded storage on *some* Postgres entities

### Main features
- sqlx to make SQL operations
- viper to get config from .env file
- docker-compose to deploy DB instances
***
## How to use:
- rename *"example.env"* file to *".env"* and fill it with your variables
- use methods:

        Storage.GetItems(ctx context.Context, itemIDs ...int64)
        Storage.AddItems(ctx context.Context, items ...models.Item)

***
#[Inspired By](https://habr.com/ru/company/ozontech/blog/705912/)




