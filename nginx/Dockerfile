FROM amd64/nginx:stable-alpine AS server-build

ENV SYSTEM_ENV_BASE_ZONE=""

WORKDIR /usr/share/nginx/html
COPY ["nginx.conf", "/etc/nginx/conf.d/default.conf"]

WORKDIR /
COPY ["entrypoint.sh", "./changeDnsZone.sh"]

RUN nginx -t && chmod +x changeDnsZone.sh

CMD [ "sh", "/changeDnsZone.sh" ]