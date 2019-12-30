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

## <img src="https://sternentstehung.de/dice-dot.png"> Included Features


* Different load balancing methods for each application
* Periodic health checks for deployed instances
* Service and server updates without downtime
* Configuration changes without restart
* Nodes with less computing resources receive less requests
* Attachment and detachment of instances on the fly
* Manage deployments by logical and physical affiliation
* Dice is passive: Making a service available for load balancing is up to you

## <img src="https://sternentstehung.de/dice-dot.png"> Quick Example

### The Scenario

Our infrastructure is quite simple: We've got two servers, _main-server_ and _another-server_. Servers, virtual machines etc. are known to Dice as _nodes_. Also, we have the services _A_, _B_ and _C_. These services might be web applications, REST APIs or authentication services for example.

<p align="center">
<br>
<img src="https://sternentstehung.de/dice-example-scenario.png">
<br>
<br>
</p>

Each service _A_, _B_ and _C_ has an instance deployed to _main-server_. An instance is a concrete executable instance of a service, like a PHP application running on Apache or a standalone Go binary. Additionally, there are instances of _A_ and _B_ deployed to _another-server_ because they're under heavy load.

### Getting started

Let's make our infrastructure available to Dice. After starting the Dice service, we can register our servers:

````shell script
$ dice node create main-server --weight 2
$ dice node create another-server
````

Registering these servers will help Dice choosing an appropriate service instance later. `--weight 2` indicates that `main-server` has double computing capacaties compared to our other servers.

Before we're able to register our service instances, we have to tell Dice about the services itself. For this example, we'll just create service _A_.

````shell script
$ dice service create A
````

We also have to specify the public URL that belongs to our service. Mapping an URL to a service is fairly simple:

````shell script
$ dice service url A example.com
````

When an request for `example.com` hits Dice, it will forward the request to an instance of service _A_. We can register such an instance like so:

````shell script
$ dice instance create A main-server 172.21.21.1:8080 --name main-instance
````

This tells Dice to register an instance of service `A` that has been deployed to `main-server` and is available at `172.21.21.1:8080`.

Attaching the created instance to Dice will make it available as a target for load balancing:

````shell script
$ dice instance attach main-instance
````

We could also use the full instance URL here, but names like `main-instance` are more convenient. The created instance will now receive incoming request for `example.com`.

## <img src="https://sternentstehung.de/dice-dot.png"> Installation

Download the [latest release of Dice](https://github.com/dominikbraun/dice/releases).

... Install instructions here ...

## <img src="https://sternentstehung.de/dice-dot.png"> Usage

Dice is a passive tool. This means that you 
