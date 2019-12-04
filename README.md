<p align="center">
<br>   
<img src="https://sternentstehung.de/dice-black-100.png">
<br>
<br>
</p>

<h3 align="center">Dice &ndash; Simple load balancing for non-microservice infrastructures.</h3>

<p align="center">
<br>
<img src="https://circleci.com/gh/dominikbraun/foodunit.svg?style=shield">
<img src="https://goreportcard.com/badge/github.com/dominikbraun/foodunit">
<img src="https://www.codefactor.io/repository/github/dominikbraun/dice/badge?s=0f13518b90c29be6bc3ec4ff537581a2e5c51c6a" />
<img src="https://img.shields.io/github/v/release/dominikbraun/foodunit?sort=semver">
<img src="https://img.shields.io/github/license/dominikbraun/foodunit">
<br>
<br>
</p>

---

<br>
Dice is an ergonomic, flexible, easy to use load balancer designed for non-microservice infrastructures.

Features include:
* Different load balancing methods for each application
* Periodic health checks for deployed instances
* Service updates without downtime
* Configuration changes without restart
* Nodes with less computing resources receive less requests
* Attachment and detachment of instances on the fly
* Dice is passive: Making a service available for load balancing is up to you

### Example

There are the services _A_, _B_ and _C_. These might be web applications or authentication services. Our infrastructure is simple: We've got a _MainServer_ and _AnotherServer_. Servers, virtual machines etc. are known to Dice as _nodes_.

<p align="center">
<br>
<img src="https://sternentstehung.de/dice-example.png">
<br>
<br>
</p>

Service _A_, _B_ and _C_ have an instance deployed on the _MainServer_. A service instance is a concrete executable. Additionally, there are instances of _A_ and _B_ deployed on _AnotherServer_ because they're under heavy load.

Let's make the servers available for Dice:

````shell script
$ dice node create MainServer 172.21.21.1
$ dice node create AnotherServer 172.21.21.2
````

Then we tell Dice about our services _A_ and _B_:

````shell script
$ dice service create A --attach
$ dice service create B
````

Attaching the node with `--attach` will make the node available as a target for load balancing instantly.
After that, we register a service deployment.

````shell script
$ dice instance create A MainServer 127.21.21.1:8080 --name MyInstance
````

This registers an instance available under `127.21.21.1:8080` that is deployed on `MainServer` and associated with service `A`.

Now we attach the created instance to Dice:

````shell script
$ dice instance attach MyInstance
````

We could also use the instance URL here. `MyInstance` will now receive incoming requests for service `A`.

### Installation

Download the [latest release of Dice](https://github.com/dominikbraun/dice/releases).

... Install instructions here ...

### Usage

Dice is a passive tool. This means that you 