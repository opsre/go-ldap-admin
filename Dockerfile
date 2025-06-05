FROM registry.cn-hangzhou.aliyuncs.com/eryajf/golang:1.18.10-alpine3.17  AS builder

WORKDIR /app

ENV GOPROXY https://goproxy.io

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk upgrade && apk add --no-cache --virtual .build-deps \
    ca-certificates gcc g++ curl upx

ADD . .

COPY --from=registry.cn-hangzhou.aliyuncs.com/eryajf/docker-compose-wait /wait .

RUN release_url=$(curl -s https://api.github.com/repos/opsre/go-ldap-admin-ui/releases/latest | grep "browser_download_url" | grep -v 'dist.zip.md5' | cut -d '"' -f 4); wget $release_url && unzip dist.zip && rm dist.zip && mv dist public/static

RUN sed -i 's@localhost:389@openldap:389@g' /app/config.yml \
    && sed -i 's@host: localhost@host: mysql@g'  /app/config.yml && go build -o go-ldap-admin . && upx -9 go-ldap-admin && upx -9 wait

### build final image
FROM registry.cn-hangzhou.aliyuncs.com/eryajf/alpine:3.19

LABEL maintainer eryajf@163.com

WORKDIR /app

COPY --from=builder /app/wait .
COPY --from=builder /app/LICENSE .
COPY --from=builder /app/config.yml .
COPY --from=builder /app/go-ldap-admin .

RUN chmod +x wait go-ldap-admin

# see wait repo: https://github.com/ufoscout/docker-compose-wait
CMD ./wait && ./go-ldap-admin