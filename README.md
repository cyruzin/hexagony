# Hexagony Records

[![build](https://github.com/cyruzin/hexagony/workflows/build/badge.svg)](https://github.com/cyruzin/hexagony/actions?query=workflow%3Abuild+branch%3Amaster) [![Coverage Status](https://coveralls.io/repos/github/cyruzin/hexagony/badge.svg?branch=master)](https://coveralls.io/github/cyruzin/hexagony?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/cyruzin/hexagony)](https://goreportcard.com/report/github.com/cyruzin/hexagony) [![GitHub license](https://img.shields.io/github/license/Naereen/StrapDown.js.svg)](https://github.com/Naereen/StrapDown.js/blob/master/LICENSE)

Hexagony name was taken from one of the Hate's songs. Hate is a death metal band from Poland.

## Running

First rename the **.env_example** file to **.env** and fill the variables if you want.

Next, make sure you have Docker and Docker Compose installed and then run the following command:

```sh
$ docker-compose up --build
```

Or in Detached Mode:

```sh
$ docker-compose up -d --build 
```

Then, check the app up and running: http://localhost:8000.

## Documentation

Access: http://localhost:8000/docs/index.html

Generate doc: 

```sh
$ swag init -g ./cmd/server/main.go
```

## Schema

Download the schema inside **docs** folder and import in your Insomnia application or another request tool.
You need to specify Authorization header along with the Bearer token.

## Generate Token

Use the following credentials:

**User**: john@doe.com

**Pass**: 12345678

## Contributing

Feel free to send pull requests, let's improve this project.

## License

MIT
