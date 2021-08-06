FROM alpine
RUN apk add ca-certificates
ENTRYPOINT ["/usr/local/bin/katapult"]
WORKDIR /katapult
COPY katapult /usr/local/bin/katapult
