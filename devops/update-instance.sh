#!/bin/bash
set -euo pipefail

git switch master
git remote update
if [[ $(git status | grep "Your branch is behind" | wc -l) -eq 0 ]]; then
  echo The branch is still up to date, exiting
  exit 0
fi

echo The branch is behind, pulling and deploying...
git pull
./devops/deploy.sh
