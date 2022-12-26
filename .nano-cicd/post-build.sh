CI_COMMIT_SHA=$(git rev-parse --short HEAD)

docker tag gitea.notnulldev.com/notnulldev/nnd:latest gitea.notnulldev.com/notnulldev/nnd:$CI_COMMIT_SHA

docker push gitea.notnulldev.com/notnulldev/nnd:$CI_COMMIT_SHA
docker push gitea.notnulldev.com/notnulldev/nnd:latest