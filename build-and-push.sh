#!/bin/bash

# Переменные
USERNAME="xomrkob"
IMAGE_NAME="go-http-server"
VERSION="1.0.0"

echo "Building image ${USERNAME}/${IMAGE_NAME}:${VERSION}"
docker build \
  --build-arg VERSION=${VERSION} \
  --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  -t ${USERNAME}/${IMAGE_NAME}:${VERSION} \
  -t ${USERNAME}/${IMAGE_NAME}:latest \
  .

if [ $? -eq 0 ]; then
  echo "Build successful"

  echo "Pushing to Docker Hub..."
  docker push ${USERNAME}/${IMAGE_NAME}:${VERSION}
  docker push ${USERNAME}/${IMAGE_NAME}:latest

  echo "Image pushed successfully!"
else
  echo "Build failed"
  exit 1
fi
