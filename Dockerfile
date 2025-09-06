FROM docker.cnb.cool/znb/images/golang:1.25.0-alpine3.22  AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.io

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk upgrade && apk add --no-cache --virtual .build-deps \
    ca-certificates gcc g++ curl upx

ADD . .

COPY --from=docker.cnb.cool/znb/images/docker-compose-wait /wait .
COPY --from=docker.cnb.cool/opsre/go-ldap-admin-ui /app/dist public/static/dist

RUN sed -i 's@localhost:389@openldap:389@g' /app/config.yml \
    && sed -i 's@host: localhost@host: mysql@g'  /app/config.yml && go build -o go-ldap-admin . && upx -9 go-ldap-admin && upx -9 wait

### build final image
FROM docker.cnb.cool/znb/images/alpine:latest

LABEL maintainer=eryajf@163.com

WORKDIR /app

COPY --from=builder /app/wait .
COPY --from=builder /app/LICENSE .
COPY --from=builder /app/config.yml .
COPY --from=builder /app/go-ldap-admin .

RUN chmod +x wait go-ldap-admin

# see wait repo: https://github.com/ufoscout/docker-compose-wait
CMD ./wait && ./go-ldap-admin