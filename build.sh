#/bin/bash

#
# Author: Daniel Struck
# Date: 2021-05-21
#

IMAGE_NAME=kafkatool

docker build -f Dockerfile --rm -t "${IMAGE_NAME}" .

if [[ ! -d $(pwd)/build ]]
then
  mkdir $(pwd)/build
else
  rm -fr $(pwd)/build
  mkdir $(pwd)/build
fi

CONTAINER_ID=$(docker create "${IMAGE_NAME}")
docker cp "${CONTAINER_ID}":/build/kafkatool $(pwd)/build/
docker rm "${CONTAINER_ID}"
