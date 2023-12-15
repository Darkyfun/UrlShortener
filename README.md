# UrlShortener
Simple url shortener with cache and RDB.

#### This project is intended for educational purposes

Stack:
- pgx/pgxpool
- go-redis
- Zap
- Viper
- Gin

RDB should have "service" and "services" databases created before start for correct work.

This repo contains two branches: master and container. 

##### master

You need to specify the environment variable that contains a path to config file with 'config' flag. I use SHORTENER_CONFIG_PATH environment variable.
So my usage is:

go run main.go -config=SHORTENER_CONFIG_PATH

#### Container

Just use docker compose to run multiple containers

