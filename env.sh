#!/bin/sh

# Output list of docker run options to pass required environment variables.

if [ "${TRAVIS}" = "true" ]
then
  echo "-e TRAVIS"
  echo "-e TRAVIS_PULL_REQUEST"
  echo "-e TRAVIS_PULL_REQUEST_SLUG"
fi
