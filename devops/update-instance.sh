#!/bin/bash
set -euo pipefail

usage() {
  cat << EOF
Usage: ./devops/update-instance.sh <user@ip>
EOF
  exit -1
}

CONNECTION_STRING=${1:-""}

if [ -z $CONNECTION_STRING ]; then
  usage
fi

ssh -t $CONNECTION_STRING << EOF
export NVM_DIR="\$HOME/.nvm"
\. "\$NVM_DIR/nvm.sh"
\. "\$NVM_DIR/bash_completion"
cd stored-strings
git pull
nohup ./devops/deploy.sh 1>api.stdout 2>api.stderr &
EOF
