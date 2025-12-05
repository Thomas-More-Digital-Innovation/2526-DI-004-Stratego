docker compose -f docker-compose.dev.yml exec -T db psql -U stratego -d stratego < code/backend/schema.sql
./code/backend/seed.sh