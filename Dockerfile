
# STAGE 1: Compilation Environment 
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN mkdir -p db && \
    wget -qO db/GeoLite2-City.mmdb "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb" && \
    wget -qO db/GeoLite2-ASN.mmdb "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-ASN.mmdb"

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ipinfo-service ./main.go

# STAGE 2: Minimal Runtime Environment 
FROM gcr.io/distroless/static-debian12:latest AS runner

WORKDIR /dist


COPY --from=builder /app/ipinfo .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/db/ ./db/


EXPOSE 9755

USER nonroot:nonroot

ENTRYPOINT ["./ipinfo"]