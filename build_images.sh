#!/usr/bin/env bash
set -e

export COMMIT=`git rev-parse HEAD`
export BASE_TAG="joelsaunders91"

export BACKEND_TAG=$BASE_TAG/gothebookofjoel-backend:$COMMIT
export MIGRATE_TAG=$BASE_TAG/gothebookofjoel-migrate:$COMMIT
export NGINX_TAG=$BASE_TAG/gothebookofjoel-nginx:$COMMIT

echo "building frontend"
npm run --prefix ./ui build
rm -rf nginx/www
mkdir -p nginx/www/
cp -a ui/build/* nginx/www/
echo "frontend build finished"

echo "building backend image"
docker build -t $BACKEND_TAG -q ./api
echo "finished building backend image"
echo "$BACKEND_TAG"

echo "building nginx image"
docker build -t $NGINX_TAG -q ./nginx
echo "finished building nginx image"
echo "$NGINX_TAG"

echo "building migrations image"
docker build -t $MIGRATE_TAG -q ./api/migrations
echo "finished building migrations image"
echo "$MIGRATE_TAG"

if [[ "$#" -eq 1 ]]; then
    echo "Pushing images to registry"
    docker push $BACKEND_TAG
    docker push $NGINX_TAG
    docker push $MIGRATE_TAG

fi

if [ "$1" == "dev" ]; then
    echo "adding tags to kubectl for dev"
    sed -i "s#image: backend#image: ${BACKEND_TAG}#" ./k8s/dev/deployment.yaml
    sed -i "s#image: nginx#image: ${NGINX_TAG}#" ./k8s/dev/deployment.yaml
    sed -i "s#image: migrate-image#image: ${MIGRATE_TAG}#" ./k8s/dev/deployment.yaml
#    kubectl apply -Rf ./k8s/dev/
fi

if [ "$1" == "prod" ]; then
    echo "adding tags to kubectl for prod"
    sed -i "s#image: backend#image: ${BACKEND_TAG}#" ./k8s/prod/deployment.yaml
    sed -i "s#image: nginx#image: ${NGINX_TAG}#" ./k8s/prod/deployment.yaml
    sed -i "s#image: migrate-image#image: ${MIGRATE_TAG}#" ./k8s/prod/deployment.yaml
#    kubectl apply -Rf ./k8s/prod/
fi
