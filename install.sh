#!/bin/sh

set -eu

arch=
case $(uname -m) in
  x86_64)
    arch=amd64
    ;;
  arm64)
    arch=arm64
    ;;
  aarch64)
    arch=arm64
    ;;
  *)
    echo "unsupported arch $(uname -m)" >&2
    exit 1
    ;;
esac

os=
ext=
case $(uname -s) in
  Linux)
    os=linux
    ext=.tar.gz
    ;;
  Darwin)
    os=darwin
    ext=.tar.gz
    ;;
  *)
    echo "unsupported OS $(uname -s)" >&2
    exit 1
    ;;
esac

api_auth=
if [ -n "${GITHUB_TOKEN:-${TRAVIS_BOT_GITHUB_TOKEN:-}}" ]
then
  api_auth="-H \"Authorization: token ${GITHUB_TOKEN:-${TRAVIS_BOT_GITHUB_TOKEN}}\""
fi

gh_api_base=${GITHUB_API_URL_BASE:-https://api.github.com}

tag=${1:-latest}
rel=
if [ ${tag} = "latest" ]
then
  rel=$(eval curl \
    ${api_auth} \
    -s --retry 4 \
    ${gh_api_base}/repos/at-wat/gh-pr-comment/releases/latest)
else
  rel=$(eval curl \
    ${api_auth} \
    -s --retry 4 \
    ${gh_api_base}/repos/at-wat/gh-pr-comment/releases/tags/${tag})
fi

url=$(echo "${rel}" | sed -n 's/.*"browser_download_url":\s*"\([^"]*\)"/\1/p' | grep "_${os}_${arch}${ext}" | head -n1)
echo ${url} 
if [ -z "${url}" ]
then
  echo "supported binary not found" >&2
  exit 1
fi

tmpdir=$(mktemp -d)
dir=${2:-~/.local/bin/}
curl -sL ${url} | tar xzfv - -C ${tmpdir}/

mkdir -p ${dir}
cp ${tmpdir}/gh-pr-comment ${dir}
cp ${tmpdir}/gh-pr-upload ${dir}
rm -rf ${tmpdir}
