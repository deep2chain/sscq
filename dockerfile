#
# Written by junying, 2019-04-10
# Includes downloading dependencies & building 
#
# Build HtdfService in a stock Go builder container
FROM golang:1.11-alpine as construction

# Set up dependencies
ENV PACKAGES make git curl build-base gcc musl-dev linux-headers 

# Set working directory for the build
WORKDIR /go/src/github.com/deep2chain/sscq

RUN apk add --update $PACKAGES

# Add source files
COPY . .

# dependency check, build
RUN make all

# Pull HtdfService into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --update ca-certificates
WORKDIR /root
COPY --from=construction /go/src/github.com/deep2chain/sscq/build/* /usr/local/bin/

EXPOSE 1317 26656 26657
CMD ["ssd"]
ENTRYPOINT ["ssd"]