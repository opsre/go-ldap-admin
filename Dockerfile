FROM docker.cnb.cool/znb/images/alpine

WORKDIR /app
LABEL maintainer=eryajf@163.com

ENV TZ=Asia/Shanghai
ENV BINARY_NAME=go-ldap-admin

ARG TARGETOS
ARG TARGETARCH

COPY LICENSE .
COPY config.yml .
COPY bin/${BINARY_NAME}_${TARGETOS}_${TARGETARCH} ${BINARY_NAME}
COPY --from=docker.cnb.cool/znb/images/docker-compose-wait /wait .

RUN chmod +x wait go-ldap-admin &&\
    sed -i 's@localhost:389@openldap:389@g' /app/config.yml \
    && sed -i 's@host: localhost@host: mysql@g'  /app/config.yml

# see wait repo: https://github.com/ufoscout/docker-compose-wait
CMD ./wait && ./go-ldap-admin