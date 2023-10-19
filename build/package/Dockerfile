FROM golang:1.21.0-alpine as base
WORKDIR /app
COPY . .
RUN go mod download
ENV SHORTENER_CONFIG_PATH=/app/config/conf.yaml
RUN go build -o app.exe cmd/main.go
ENTRYPOINT ["/app/app.exe"]
CMD ["-config=SHORTENER_CONFIG_PATH"]

FROM scratch
COPY --from=base /app/app.exe /bin/app.exe
COPY --from=base /app/config/conf.yaml /bin/conf.yaml
ENTRYPOINT ["app.exe"]
CMD ["-config=SHORTENER_CONFIG_PATH"]