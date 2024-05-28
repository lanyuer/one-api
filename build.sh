VERSION=1.1
docker build -t one-api:${VERSION} .
docker tag one-api:${VERSION} mylanyuer/one-api:${VERSION}
docker push mylanyuer/one-api:${VERSION}