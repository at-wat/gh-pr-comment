#!/bin/sh

set -eu

arch=
case $(uname -p) in
  x86_64)
    arch=amd64
    ;;
  *)
    echo "unsupported arch" >&2
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
  *)
    echo "unsupported OS" >&2
    exit 1
    ;;
esac

tag=${1:-latest}
rel=
if [ ${tag} = "latest" ]
then
  rel=$(curl -s --retry 4 https://api.github.com/repos/at-wat/gh-pr-comment/releases/latest)
else
  rel=$(curl -s --retry 4 https://api.github.com/repos/at-wat/gh-pr-comment/releases/tags/${tag})
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
