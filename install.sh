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

rel=$(curl -s https://api.github.com/repos/at-wat/gh-pr-comment/releases/${1:-latest})
echo "${rel}" | sed -n 's/.*"browser_download_url":\s*"\([^"]*\)"/\1/p' | while read url
do
  if echo ${url} | grep -q "_${os}_${arch}${ext}"
  then
    echo ${url} 
    tmpdir=$(mktemp -d)
    curl -sL ${url} | tar xzfv - -C ${tmpdir}/
    cp ${tmpdir}/gh-pr-comment ~/.local/bin/
    cp ${tmpdir}/gh-pr-upload ~/.local/bin/
    rm -rf ${tmpdir}
    exit 99
  fi
done

if [ $? -eq 99 ]
then
  exit 0
fi

echo "supported binary not found" >&2
exit 1
