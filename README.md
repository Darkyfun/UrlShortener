# UrlShortener
Simple url shortener with cache and SQL base.

You need to specify the environment variable that contains a path to config file with 'config' flag. I am using SHORTENER_CONFIG_PATH environment variable.
So my usage is:

go run main.go -config=SHORTENER_CONFIG_PATH
