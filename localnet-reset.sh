make zetanode
docker build -t zetanode .
cd contrib/localnet/orchestrator/
docker compose down --remove-orphans
docker compose up -d


