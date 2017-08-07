#!/bin/bash

function gh-pr-comment()
{
  if [ $# -lt 2 ];
  then
    echo "usage: gh-pr-comment title comment" >&2
    echo "  env:" >&2
    echo "     - TRAVIS_REPO_SLUG        : user/repos" >&2
    echo "     - TRAVIS_PULL_REQUEST     : pull request number" >&2
    echo "     - TRAVIS_BOT_GITHUB_TOKEN : token with comment write permission" >&2
    echo ' note: the script escapes special charactors for json encoding.' >&2

    return
  fi

  if [ ${TRAVIS_PULL_REQUEST} != "false" ];
  then
    text=`echo "$2" | sed 's/\\\\/\\\\\\\\/g' | sed -n '1h;1!H;${x;s/\n/\\\\n/g;p;}' | sed 's/\"/\\\\"/g' | sed 's/\t/\\\\t/g' | sed 's/\//\\\\\//g'`
    curl -vs -X POST -H 'Content-Type:application/json' -d "{\"body\":\"## $1\n\n$text\"}" \
      https://api.github.com/repos/${TRAVIS_REPO_SLUG}/issues/${TRAVIS_PULL_REQUEST}/comments?access_token=${TRAVIS_BOT_GITHUB_TOKEN} 2> /dev/null
  else
    echo 'gh-pr-comment: TRAVIS_PULL_REQUEST is false.' >&2
  fi
}

