docker compose exec -T gostrategy psql -U gostrategy -d gostrategy < code/backend/schema.sql
./code/backend/seed.sh