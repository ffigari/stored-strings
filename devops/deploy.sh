#!/bin/bash
set -euo pipefail

lsof -ti:3000 && kill $(lsof -ti:3000)
nohup node runner.js spawn &

NGINX_CONFIG_PATH=/etc/nginx/sites-enabled/stored-strings
rm -rf $NGINX_CONFIG_PATH
cp devops/nginx-config $NGINX_CONFIG_PATH
nginx -t
nginx -s reload
