<p align="center">
<br>
<br>
<img src="https://sternentstehung.de/dice-colored-100.png">
<br>
<br>
</p>

<h3 align="center">Dice &ndash; Simple load balancing for non-microservice infrastructures.</h3>

<p align="center">
<img src="https://circleci.com/gh/dominikbraun/foodunit.svg?style=shield">
<img src="https://goreportcard.com/badge/github.com/dominikbraun/foodunit">
<img src="https://www.codefactor.io/repository/github/dominikbraun/dice/badge?s=0f13518b90c29be6bc3ec4ff537581a2e5c51c6a" />
<img src="https://img.shields.io/github/v/release/dominikbraun/foodunit?sort=semver">
<img src="https://img.shields.io/badge/license-Apache--2.0-brightgreen">
<br>
<br>
<br>
</p>

---

:game_die: Dice is an ergonomic, flexible, easy to use load balancer designed for non-microservice, static infrastructures.

**Project status:** In development, some features are experimental yet.

## <img src="https://sternentstehung.de/dice-dot.png"> Included Features


* Different load balancing methods for each application
* Periodic health checks for deployed instances
* Service and server updates without downtime
* Configuration changes without restart
* Nodes with less computing resources receive less requests
* Attachment and detachment of instances on the fly
* Deployments managed by logical and physical affiliation
* Dice is passive and explicit: Make services available for load balancing yourself

## <img src="https://sternentstehung.de/dice-dot.png"> Simple Example

### The Scenario

Our infrastructure is quite simple: We've got two servers, _main-server_ and _another-server_. Servers, virtual machines etc. are known to Dice as _nodes_. Also, we have the services _A_, _B_ and _C_. These services might be web applications, REST APIs or authentication services for example.

<p align="center">
<br>
<img src="https://sternentstehung.de/dice-example-scenario.png">
<br>
<br>
</p>

Each service _A_, _B_ and _C_ has an instance deployed to _main-server_. An instance is a concrete executable instance of a service, like a PHP application running on Apache or a standalone Go binary. Additionally, there are instances of _A_ and _B_ deployed to _another-server_ because they're under heavy load.

### Setting up our environment

Let's make our infrastructure available to Dice. After starting Dice, we can register our servers:

````shell script
$ dice node create --attach --weight=2 main-server
$ dice node create --attach another-server
````

Registering these servers will help Dice choosing an appropriate service instance later. `--weight=2` indicates that `main-server` has double computing capacities.

After that, we have to tell Dice about our services. For this example, we'll just create service _A_.

````shell script
$ dice service create --url=example.com --enable A
````

By using `--url=example.com`, we specify a public URL that the service is associated with. We can add or remove these URLs later as well. When a request for `example.com` hits Dice, it will forward the request to an instance of service _A_.

### Start load balancing

We can register our instances of _A_ like so:

````shell script
$ dice instance create --name=first-instance A main-server 172.21.21.1:8080
$ dice instance create --name=second-instance A another-server 172.21.21.2:8080
````

For example, the first command tells Dice to register an instance of service `A` that has been deployed to `main-server` and is available at `172.21.21.1:8080`.

Attaching the created instances to Dice will make them available for load balancing:

````shell script
$ dice instance attach first-instance
$ dice instance attach second-instance
````

We could also use the full instance URL here, but names like `first-instance` are more convenient. Incoming requests for `example.com` will now be balanced among our instances.

## <img src="https://sternentstehung.de/dice-dot.png"> Installation

Download the [latest release of Dice](https://github.com/dominikbraun/dice/releases) [...]

## <img src="https://sternentstehung.de/dice-dot.png"> Getting started

[...]
