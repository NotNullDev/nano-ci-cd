# docker compose rm -f $APP_NAME
echo $CONTAINER_REGISTRY_PASSWORD | docker login -u $CONTAINER_REGISTRY_USERNAME --password-stdin $CONTAINER_REGISTRY_URL
