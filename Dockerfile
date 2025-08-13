# ini-file set/get/del in a container
#
# docker run --rm -it -v /tmp:/tmp bitnami/ini-file set -k "title" -v "A wonderful book" -s "My book" /tmp/my.ini
# docker run --rm -it -v /tmp:/tmp bitnami/ini-file get -k "title" -s "My book" /tmp/my.ini
# docker run --rm -it -v /tmp:/tmp bitnami/ini-file del -k "title" -s "My book" /tmp/my.ini
#

FROM bitnami/golang:1.25 as build

WORKDIR /go/src/app
COPY . .

RUN rm -rf out

RUN make

FROM bitnami/minideb:bookworm

COPY --from=build /go/src/app/out/ini-file /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/ini-file"]
