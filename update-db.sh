docker compose exec -T stratego psql -U stratego -d stratego < code/backend/schema.sql
./code/backend/seed.sh