FROM ubuntu:20.04 

COPY dist/chip-distribution /opt/chip-distribution/
COPY source/ /opt/chip-distribution/source

ENV GIN_MODE=release

EXPOSE 8080

WORKDIR /opt/chip-distribution/

CMD ["./chip-distribution"]
