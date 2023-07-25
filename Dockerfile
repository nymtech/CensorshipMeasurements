# OONI Release 3.17 is working with go version 1.19 (see GOVERSION file)
FROM golang:1.19.6-alpine
RUN apk update && apk upgrade
RUN apk add git gcc libc-dev

COPY ./files /nym_files
RUN git clone -b release/3.17 --single-branch https://github.com/ooni/probe-cli.git /ooni-probe
RUN cp /nym_files/registry_nym.go /ooni-probe/internal/registry/nym.go
RUN cp -r /nym_files/experiment_nym /ooni-probe/internal/experiment/nym
WORKDIR /ooni-probe
RUN go build -v -ldflags '-s -w' ./internal/cmd/miniooni

CMD ./miniooni --no-collector --yes --verbose -O NymValidatorURL=${OONI_NYMVALIDATORURL:-https://validator.nymtech.net} -o /results/report.jsonl nym && chmod 666 /results/report.jsonl
