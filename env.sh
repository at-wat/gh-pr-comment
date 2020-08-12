#!/bin/sh

# Output list of docker run options to pass required environment variables.

if [ "${TRAVIS}" = "true" ]
then
  echo "-e TRAVIS"
  echo "-e TRAVIS_PULL_REQUEST"
  echo "-e TRAVIS_PULL_REQUEST_SLUG"
  echo "-e TRAVIS_REPO_SLUG"
fi

if [ -n "${GITHUB_API_URL_BASE}" ]
then
  echo "-e GITHUB_API_URL_BASE"
fi

if [ -n "${GITHUB_TOKEN}" ]
then
  echo "-e GITHUB_TOKEN"
fi

if [ -n "${TRAVIS_BOT_GITHUB_TOKEN}" ]
then
  echo "-e TRAVIS_BOT_GITHUB_TOKEN"
fi
