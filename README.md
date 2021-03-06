# Todo Microservice

There was a list of technologies I wanted to use for a project, and this todo
app is a test to see how that would work.

This is just the server, but you can interact with it with
[GraphiQL](https://github.com/graphql/graphiql), which is also served on
`localhost:8080/`.  As I'm going to use Elm for the frontend in said project,
there will soon(-ish) be a client available in another repo for speaking with
this server.

### Contributing
If you want to, in any way, contribute to this experiment, or if you have
questions about anything, feel free to open issues or PR's.

### Technology/library roadmap
- [x] go-kit
- [x] graphql with [graphql-go](https://github.com/graphql-go/graphql) and [granate](https://github.com/granateio/granate)
- [ ] paging of todos
- [x] logging (go-kit)
- [x] instrumenting (prometheus)
- [ ] opentracing (with appdash)
- [x] client authentication
- [ ] authorisation *(in progress)*
- [ ] Proxying?
- [ ] Load balancing?
- [ ] Circuit breaking?
- [ ] Throtling?

### Roadmap otherwise
- [x] Todos
- [x] Users
- [x] Login

While I'm at it I want to try out a whole bunch of things. The `todo service`
is exposed with graphql, but I want to expose it via a REST-ful api as well.
This is both to see if the architecture would allow it without making
compromises, and to test whether abstractions are done properly, or if logic
leaks where it shouldn't. The REST stuff has low priority, though.

### Running
```
$ cd $PROJECT_ROOT
$ go get . # Install dependencies
$ npm install # or `yarn`
$ gulp
```

Visit `localhost:8080/` to interact with the server through GraphiQL.
You can visit `localhost:8080/metrics` to see data from
instrumentation.
Use `localhost:8080/graphql` if you need to speak directly with the GraphQL API.

### Inspiration
- [go-kit/kit](github.com/go-kit/kit) (shipping examle in particular)
- [narqo/test-graphql](https://github.com/narqo/test-graphql)

### Resources
- [graphql.org](http://graphql.org/learn/)
- [opentracing.io](http://opentracing.io/documentation/)
- [graphql-go](https://github.com/graphql-go/graphql)
- [granate](https://github.com/granateio/granate)
- [Microservice pattern](http://microservices.io/patterns/microservices.html)
- [Hexagonal architecture](http://alistair.cockburn.us/Hexagonal+architecture)
