# Zest

## Installation

Make sure you have a working Node and Go environment.

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

    go build -v -o zest-app

Create/migrate/seed the database (SQLite is default if the environment variable `DATABASE_URL` is not set):

    ./zest-app resetDB

Generate the Swagger documentation:

    ./zest-app generateDoc

Run the app (on `localhost:3000` by default):

    ./zest-app

Enjoy the server freshly created.

## API documentation

The Swagger API documentation is available on `/explorer`.

## Why Zest ?

The main purpose of Zest is to speed up the development of simple resource oriented apps, without doing compromises on the performances and the scalability.

I see Zest more like a boilerplate/framework than a "toolbox" framework like [Beego](https://github.com/astaxie/beego) or [Revel](https://github.com/revel/revel).

A good practice could be to clone the Zest repo and directly build an app on it.

## Features

* Modular design based on the [Clean architecture](https://blog.8thlight.com/uncle-bob/2012/08/13/the-clean-architecture.html).
* Strongly opinionated.
* [GORM](https://github.com/jinzhu/gorm) powered (PostgreSQL and SQLite currently supported).
* Zero runtime reflection.
* High productivity: generate and run.
* Auto API documentation.
* Out of the box session management and signin/signup/signout methods.
* Sessions caching.
* Fully dynamic permissions/role management, made insanely fast thanks to out of the box caching.
* App automatically built thanks to the included injector.
* Easy to add business logic thanks to hooks everywhere.

## About

Inspired by [Rails](https://github.com/rails/rails) and [Loopback](https://github.com/strongloop/loopback).

This project is still in an alpha state. **DO NOT USE IT IN PRODUCTION**.

## License

MIT
