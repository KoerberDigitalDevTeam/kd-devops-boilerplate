#!/bin/sh

ROOT_DIR=/etc/nginx/conf.d/

echo "Replacing env constants in JS"
for file in $(find . -name '*.conf');
do
  echo "Processing $file ...";

  sed -i 's|SYSTEM_ENV_BASE_ZONE|'${SYSTEM_ENV_BASE_ZONE}'|g' $file

done

# Starting NGINX
nginx -g 'daemon off;'