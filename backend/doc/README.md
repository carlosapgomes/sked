## GoBackend

### Intro

This program is meant to be an starter kit for a web application in Go.

It provides basic users flows (new account, login, logout, reset password) and
implements server side sessions with a Postgres database and http cookies.

It also implements a context middleware that allow you to provide a custom
authorization layer.

It's structure is based on the Hexagon Architecture (Ports & Adapters),
which is an application pattern that provides a good decoupling of its
components making them easier to test.

Basically, the _Ports & Adapters_ Architecture defines an application and
its actors.

Inside the application are the data models and the business logic (code that
implements use cases and is also called a _Port_).

Outside the application are the actors that interact with it.

An actor that send commands or request _things_ from the application is called
a _driver actor_ and can be a graphical user interface, REST api, web/console
interface or even a test suit.

An actor that receives commands or data from the application is called
a _driven actor_ and can be a mailer or a database.

Each actor interacts with a _Port_ using an _Adapter_.

A _Driver Adapter_ is, for example, a code that understands the REST
API/Protocol and translates it into calls to specific use cases implementators.

On the other hand a _Driven Adapter_ is a code that receives commands and
requests from a use case and translates them into SQL commands, for example.

If you want to read more about the Hexagon Architecture, I suggest you
to follow these links:

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)

- [Ports and Adapters](https://jmgarridopaz.github.io/content/hexagonalarchitecture.html)
- [DDD, Hexagonal, Onion, Clean, CQRS, … How I put it all together](https://herbertograca.com/2017/11/16/explicit-architecture-01-ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together/#components)


## Building the application

From the `backend/` folder type:

`go build -o skedbacked cmd/main.go`
