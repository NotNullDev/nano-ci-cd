CI_COMMIT_SHA=$(git rev-parse --short HEAD)

echo docker tag $CONTAINER_REGISTRY_URL/$CONTAINER_REGISTRY_USERNAME/$APP_NAME:latest $CONTAINER_REGISTRY_URL/$CONTAINER_REGISTRY_USERNAME/$APP_NAME:$CI_COMMIT_SHA

echo docker push $CONTAINER_REGISTRY_URL/$CONTAINER_REGISTRY_USERNAME/$APP_NAME:$CI_COMMIT_SHA
echo docker push $CONTAINER_REGISTRY_URL/$CONTAINER_REGISTRY_USERNAME/$APP_NAME:latest
env

echo mounted file:

ls .mounted
cat .mounted

echo mounted env file:
echo $BASE_64_ENV_FILE
echo $BASE_64_ENV_FILE | base64 -d > .env
cat .env
