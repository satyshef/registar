ARG APP_NAME="registar"
ARG APP_DIR="/app"
ARG CONFIG_DIR="data"



FROM golang:1.18-alpine3.15 AS builder
LABEL stage=builder
ARG APP_NAME
ARG APP_DIR

COPY --from=satyshef/tdlib:1.8.3-alpine3.15 /usr/local/include/td /usr/local/include/td
COPY --from=satyshef/tdlib:1.8.3-alpine3.15 /usr/local/lib/libtd* /usr/local/lib/
COPY --from=satyshef/tdlib:1.8.3-alpine3.15 /usr/lib/libssl.a /usr/local/lib/libssl.a
COPY --from=satyshef/tdlib:1.8.3-alpine3.15 /usr/lib/libcrypto.a /usr/local/lib/libcrypto.a
COPY --from=satyshef/tdlib:1.8.3-alpine3.15 /lib/libz.a /usr/local/lib/libz.a
RUN apk add build-base

WORKDIR ${APP_DIR}

COPY go.mod go.sum ./
RUN apk add git && \
    apk add make && \
    go mod download && \
    go mod verify
COPY . .
RUN  go mod tidy && make -e APP_PATH=${APP_NAME}

#RUN go build --ldflags "-extldflags '-static -L/usr/local/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -ldl -lm -lssl -lcrypto -lstdc++ -lz'" -o tebot cmd/app/main.go

# finish
FROM alpine:3.15
ARG APP_NAME
ARG APP_DIR
ARG CONFIG_DIR

#ENV APP_PATH=${APP_DIR}/${APP_NAME}

WORKDIR ${APP_DIR}

COPY --from=builder ${APP_DIR}/${CONFIG_DIR} ./${CONFIG_DIR}
COPY --from=builder ${APP_DIR}/${APP_NAME} .
COPY --from=builder ${APP_DIR}/entrypoint.sh .
#COPY --from=builder ${APP_DIR}/profiles ${APP_DIR}/profiles

RUN apk add libstdc++ && chmod +x entrypoint.sh
EXPOSE 7070
CMD [ "./entrypoint.sh" ]
#CMD [ "/app/run.sh" ]
#CMD ${APP_PATH}
#CMD [ ${APP_PATH}, "-c", ${CONF_PATH} ]
