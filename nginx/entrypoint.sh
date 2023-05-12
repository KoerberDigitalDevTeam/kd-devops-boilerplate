#!/bin/sh

echo "Replacing env constants in conf file"
for file in $(find /etc/nginx/conf.d -name '*.conf');
do
  echo "Processing $file ...";

  sed -i 's|SYSTEM_ENV_BASE_ZONE|'${SYSTEM_ENV_BASE_ZONE}'|g' $file

done