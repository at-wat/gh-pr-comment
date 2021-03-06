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
ep=latest
if [ "${tag}" != "latest" ]
then
  ep=tags/${tag}
fi

curl_err_file=$(mktemp)
rel=$(eval curl \
  ${api_auth} \
  -sSL --retry 4 \
  ${gh_api_base}/repos/at-wat/gh-pr-comment/releases/${ep} 2> ${curl_err_file} || true)
curl_err="$(cat ${curl_err_file})"
rm -f ${curl_err_file}
if [ -n "${curl_err}" ]
then
  echo "failed to fetch releases: ${curl_err}" >&2
  exit 1
fi

url=$(echo "${rel}" \
  | sed -n 's/.*"browser_download_url":[ \t\r\n]*"\([^"]*\)"/\1/p' \
  | grep -e "_${os}_${arch}${ext}$" | head -n1)
echo ${url} 
if [ -z "${url}" ]
then
  echo "supported binary not found for ${os} ${arch}" >&2
  exit 1
fi

tmpdir=$(mktemp -d)
dir=${2:-~/.local/bin/}
curl -sSL ${url} | tar xzfv - -C ${tmpdir}/

mkdir -p ${dir}
cp ${tmpdir}/gh-pr-comment ${dir}
cp ${tmpdir}/gh-pr-upload ${dir}
rm -rf ${tmpdir}
