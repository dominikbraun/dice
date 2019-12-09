#
# Dockerfile for the Dice load balancer
#

# Start the build stage.
FROM golang:1.13-alpine as build
LABEL maintainer="braun@sternentstehung.de"

# Accepted build arguments.
ARG BUILD_DATE
ARG VERSION

# Important image labels.
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date=${BUILD_DATE}
LABEL org.label-schema.name="Dice"
LABEL org.label-schema.description="Dice load balancer"
LABEL org.label-schema.url="https://github.com/dominikbraun/dice"
LABEL org.label-schema.vcs-url="https://github.com/dominikbraun/dice"
LABEL org.label-schema.version=${VERSION}
LABEL org.label-schema.docker.cmd="docker container run -d -p 8080:8080 dice"

# Add git and clone the Dice repository.
RUN apk add git --no-cache
RUN git clone https://github.com/dominikbraun/dice dice

WORKDIR dice

# Build the binary withouth the symbol table and debug information.
RUN go build -v -ldflags="-s -w" -o .target/dice cmd/dice/main.go

# Start the execution stage.
FROM alpine:3.10 as exec

ENV PORT 8080

RUN mkdir /dice
WORKDIR /dice

COPY --from=build /go/dice/.target/dice .
COPY dice.yml .

EXPOSE ${PORT}
CMD ["./dice"]