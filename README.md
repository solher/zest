# Zest

## Installation

Make sure you have a working Go environment.

Simply install the [Zest app generator](https://github.com/solher/generator-zest) :

    npm install -g yo solher/generator-zest

## Getting Started

Make a new directory and `cd` into it:

    mkdir zest-app
    cd zest-app

Then create a new Zest project:

    yo zest

Add REST resources:

    yo zest:resource

Compile the project:

    go build -v

Create/migrate/seed the database:

    ./zest-app resetDB

Run the app (on `localhost:3000` by default):

    ./zest-app

Enjoy the server freshly created.

## API documentation

Coming soon...

## Features

* Modular design based on the [Clean architecture](https://blog.8thlight.com/uncle-bob/2012/08/13/the-clean-architecture.html).
* Strongly opinionated.
* [GORM](https://github.com/jinzhu/gorm) powered (PostgreSQL and SQLite currently supported).
* Zero runtime reflection.
* High productivity: generate and run.
* Out of the box session management and signin/signup/signout methods.
* Sessions caching.
* Fully dynamic permissions/role management, made insanely fast thanks to out of the box caching.
* App automatically built thanks to the included injector.
* Easy to add business logic thanks to hooks everywhere.

## About

Inspired by [Rails](https://github.com/rails/rails) and [Loopback](https://github.com/strongloop/loopback).

## License

MIT
