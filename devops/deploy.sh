#!/bin/bash
set -euo pipefail

lsof -ti:3000 && kill $(lsof -ti:3000)
node index.js build
nohup node index.js start &

cp devops/nginx-config /etc/nginx/sites-enabled/mi-rincon
nginx -t
nginx -s reload
