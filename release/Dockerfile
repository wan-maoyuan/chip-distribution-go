FROM ubuntu:18.04 

COPY dist/chip-distribution /opt/chip-distribution/
COPY source/ /opt/chip-distribution/source

ENV GIN_MODE=release

EXPOSE 8080

WORKDIR /go/src/app/

CMD ["/opt/chip-distribution/chip-distribution"]