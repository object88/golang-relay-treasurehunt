# Go Relay Starter Kit

This app is based on the [Relay Treasure](https://github.com/relayjs/relay-starter-kit) from Facebook and the [Go/Golang relay server](https://github.com/graphql-go/relay) from graphql-go.  

This kit includes an Golang app server, a GraphQL server, and transpiler-powered JS frontend, that you can use to get started building an app with Relay. For a JS-based walkthrough, see the [Relay tutorial](https://facebook.github.io/relay/docs/tutorial.html).

## Notes

* Start using Go Vendors.
* Talk a little about app signing on a Mac.

## Installation

```
npm install
```

Go packages must also be installed.  Instructions TBD.

## Running

Start a local server:

```
npm start
```

## Developing

Any changes you make to files in the `js/` directory will cause the server to
automatically rebuild the app and refresh your browser.

If at any time you make changes to `data/schema.go`, stop the server,
regenerate `data/schema.json`, and restart the server:

```
npm run update-schema
npm start
```

## License

Relay Starter Kit is [BSD licensed](./LICENSE). We also provide an additional [patent grant](./PATENTS).
