FROM golang:1.25-alpine AS budowniczy
WORKDIR /aplikacja
COPY go.mod* ./
RUN go mod download || go mod init weather-app
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o serwer main.go

FROM alpine:latest
LABEL org.opencontainers.image.authors="Kacper Sumera" \
      org.opencontainers.image.title="Aplikacja Pogodowa IMGW" \
      org.opencontainers.image.description="Zadanie 1 – aplikacja Pogodowa" \
      org.opencontainers.image.version="1.0"
RUN apk --no-cache add ca-certificates
RUN addgroup -S grupa-app && adduser -S uzytkownik-app -G grupa-app
USER uzytkownik-app
WORKDIR /home/uzytkownik-app
COPY --from=budowniczy /aplikacja/serwer .
COPY index.html .
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./serwer"]