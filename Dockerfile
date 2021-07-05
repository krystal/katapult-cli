FROM alpine
RUN apk add ca-certificates
ENTRYPOINT ["/usr/local/bin/katapult-cli"]
WORKDIR /katapult
COPY katapult-cli /usr/local/bin/katapult-cli
