# hadolint ignore=DL3007
FROM alpine:latest

WORKDIR /app/log-parser

COPY target/log-parser .
COPY integration/testdata/sample.log .

ENTRYPOINT [ "./log-parser" ]

CMD [ ]
