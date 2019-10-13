<p align="left">
  <img src="logo.png" height="70">
</p>

[![CircleCI](https://circleci.com/gh/hiconvo/api.svg?style=svg)](https://circleci.com/gh/hiconvo/api) [![codecov](https://codecov.io/gh/hiconvo/api/branch/master/graph/badge.svg)](https://codecov.io/gh/hiconvo/api)

The repo holds the source code for Convo's RESTful API ([Docs 📑](http://api.hiconvo.com/docs)). Learn more about Convo at [hiconvo.com](https://hiconvo.com). The core technologies and APIs used in the project are:

- [Golang](https://golang.org/)
- [Gorilla](https://www.gorillatoolkit.org/)
- [Google App Engine](https://cloud.google.com/appengine/docs/standard/go112/)
- [Google Datastore](https://godoc.org/cloud.google.com/go/datastore)
- [Sendgrid](https://sendgrid.com/docs/index.html)
- [Elasticsearch](https://www.elastic.co/)
- [Docker](https://docs.docker.com/)

## Development

We use docker based development. In order to run the project locally, you need to create an `.env` file and place it at the root of the project. The `.env` file should contain a Google Maps API key and a Sendgrid API key. It should look something like this:

```
GOOGLE_MAPS_API_KEY=<YOUR API KEY>
SENDGRID_API_KEY=<YOUR API KEY>
```

If you don't include this file, the app will panic during startup.

After your `.env` file is ready, all you need to do is run `docker-compose up`. The source code is shared between your machine and the docker container via a volume. The default command runs [`realize`](https://github.com/oxequa/realize), a file watcher that automatically compiles the code and restarts the server when the source changes. By default, the server listens on port `:8080`.

### Running Tests

Run `docker ps` to get the ID of the container running the API. Then run

```
docker exec -it <CONTAINER ID> go test ./...
```

Be mindful that this command will *wipe everything from the database*. There is probably a better way of doing this, but I haven't taken the time to improve this yet.

## Architecture

The architecture is very simple.

![Architecture](architecture.png)

## Code Overview

The best place to start (after looking at `main/main.go`) is `handlers/router.go`. This is where all endpoints and corresponding functions and middlewares are mounted to a [Gorilla Mux Router](https://github.com/gorilla/mux). Some notable middlewares are `bjson.WithJSON` and `bjson.WithJSONReqBody`. These ensure that all endpoints that they wrap are JSON only and make the decoded request body available in the request context.

The next place to look is `handlers/*.go`. All of the handlers are standard handler functions of the form `func(http.ResponseWriter, *http.Request)`. I decided not to alter this signature with something more app specific since the context at `http.Request.Context()` makes it easy enough to standardize common operations.

The handlers are focused on request validation, getting data from the database, and mutating models - all through higher-level APIs. Some of these operations are complicated, especially the ones concerned with validating users to be included in events or threads, but the separation of concens between models and handlers is decent enough.

One code smell in the handlers is validation. Validation is not fully steamlined. Some of it takes place in the handlers and some of it is in the models.

The core business logic is in the models. Most of these files are pretty self explanatory. The messiest of the model files is `mail.go` which handles a portion of the email rendering and sending process. Eventually, it would be nice to break off the email business into its own service, but for now, at this early stage, the current monolithic approach is fine.


