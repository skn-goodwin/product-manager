source .env
docker build -f Dockerfile --build-arg PIPELINE_SSH_KEY="${PIPELINE_SSH_KEY}" -t "product-manager" .
docker save -o product-manager.docker "product-manager"