FROM quay.io/projectquay/golang:1.22 as builder

WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o dice
# -ldflags="-X="dice/main.AppVersion=${VERSION}

FROM scratch
ENV OTLPMETRICHTTP_ENDPOINT="http://collector:3030"
ENV OTEL_DICE_ENV="prod"
WORKDIR /
COPY --from=builder /usr/src/app/dice .
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT [ "./dice" ]
