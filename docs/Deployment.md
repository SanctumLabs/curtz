# Deployment

Deployment of the application should be simple enough to perform without much setup.

## Build

If you intend to build & run the application in a VM, that can be done with a few simple steps:

```bash
make build # or go build -o ./bin/curtz app/cmd/main.go
```

> This will create a binary `curtz` in the [bin](../bin) directory

If building in a Docker container, then simply run:

``` bash
docker build -t <IMAGE_NAME>:<IMAGE_TAG>
```

> <IMAGE_NAME> is the name of the image, for example `curtz` <IMAGE_TAG> is the tag/version of the image for example `1.0.0`

## Environment

Then setup the necessary environment variables as specified in the sample [.env.sample](../.env.sample) file.

```.env
ENV=development
LOG_LEVEL=debug
LOG_JSON_OUTPUT=true
PORT=8085
DATABASE_HOST=localhost
DATABASE=curtzdb
DATABASE_USERNAME=curtzUser
DATABASE_PASSWORD=curtzPassword
DATABASE_PORT=27017
CACHE_HOST=localhost
CACHE_PORT=6370
CACHE_USERNAME=curtzUser
CACHE_PASSWORD=curtzPassword
CACHE_REQUIRE_AUTH=false
SENTRY_DSN=""
SENTRY_ENV=development
SENTRY_SAMPLE_RATE=0.3
SENTRY_ENABLED=true
```

> Sample env file

These are all provided environment variable defaults for running locally, most are self explanatory, so, this document will explain the ones that may not be clear.

`ENV` can either be `testing`, `development` or `release`. Release here is equivalent to `production`, but the underlying http framework used is [Gin](https://gin-gonic.com/) will recognize `release` over `production`

`CACHE_HOST`, `CACHE_PORT`, `CACHE_USERNAME`, `CACHE_PASSWORD`, `CACHE_REQUIRE_AUTH` all relate to caching of urls when performing redirects. Underling technology used here is [Redis](https://redis.io/). The `CACHE_REQUIRE_AUTH` is used to determine whether to configure authentication to the caching service. If set to true, then the username and password are required variables. Setting this to false will not require authentication, just ensure that's the case for the caching service :).

Similar case applies to [Sentry](https://sentry.io/welcome/). This has been setup to enable application tracking and monitoring. But has been left as optional with the `SENTRY_ENABLED` environment variable. This can be set to true or false depending on your setup. Setting it to true, will require the `SENTRY_DSN` which you can get from Sentry's website once you have signed up and configured a project. The `SENTRY_ENV` is entirely up to you as this will depend on the deployment strategy employed, the default is `development`.

Lastly, the Database. The database used here is a [NoSQL](https://en.wikipedia.org/wiki/NoSQL) database. [MongoDB](https://www.mongodb.com/) has been picked as the NoSQL database of choice. Therefore the `DATABASE_` environment variables should be set to enable this connection, however, as you may have noticed, these environment variables don't necesarily tell us what the underlying database type is. This is by choice, allowing us to change the values & only simply change the underlying database client connection without affecting how the application really runs. Details on this can be found in the [Architecture](./Architecture.md).

## Running

Once you have all that setup, you can now run the application. If running in a VM or a hosted environment then simply run it with `go run ./bin/curtz`. That should be it. Of course this depends on the infrastracture setup you have.

If running it with Docker, run the docker container with:

```bash
docker run -p 8085:8085 --name=curtz-api <IMAGE_NAME>:<IMAGE_TAG>
```

> The <IMAGE_NAME>:<IMAGE_TAG> are the name and the tag of the image used as specified in [Build](#build) what you used to build

That should be it for deployment. Depending on your infrastructure you should be able to view the logs of the running application.
