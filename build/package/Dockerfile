FROM golang:1.21.0-alpine
WORKDIR /app
COPY . .
RUN go mod download
ENV SHORTENER_CONFIG_PATH=/app/config/conf.yaml
RUN go build -o app.exe cmd/main.go
ENTRYPOINT ["/app/app.exe"]
CMD ["-config=SHORTENER_CONFIG_PATH"]