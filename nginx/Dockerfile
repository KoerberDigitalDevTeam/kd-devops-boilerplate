FROM amd64/nginx:stable AS server-build

WORKDIR /usr/share/nginx/html

ENV SYSTEM_ENV_BASE_ZONE=""

COPY ["nginx.conf", "/etc/nginx/conf.d/default.conf"]

COPY ["entrypoint.sh", "changeDnsZone.sh"]

RUN nginx -t && chmod +x changeDnsZone.sh

CMD [ "sh", "./changeDnsZone.sh" ]