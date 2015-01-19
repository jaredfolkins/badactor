![badactor logo](https://raw.githubusercontent.com/jaredfolkins/badactor_logo/master/badactor_logo_300x300.png) [![Build Status](https://travis-ci.org/jaredfolkins/badactor.svg?branch=master)](https://travis-ci.org/jaredfolkins/badactor) [![Coverage Status](https://img.shields.io/coveralls/jaredfolkins/badactor.svg)](https://coveralls.io/r/jaredfolkins/badactor?branch=master)

# BadActor 

BadActor is an in-memory, application driven jailer built in the spirit of fail2ban. It can be used as either middleware or a standalone server (TODO) with its primary goal to increase the expense for "bad actors" who engage in system probing or attacks.

# Install

```bash
$ go get github.com/jaredfolkins/badactor
```

# Use Case

A common use case for BadActor is jailing an offender who fails to login to your website (N) times as this can signal a bruteforce attempt.

# Tutorial

Checkout [badactor.org](http://badactor.org) for a tutorial.

# Design

- speed (subsecond response underload and submillisecond with standard operations)
- no external dependencies 
- solid code coverage and thorough tests
- accessable to multiple systems across the internet using an HTTP API (TODO)

# Does It Scale? 

BadActor can be included in your go application and ran concurrently. It can also be stood up on its own, ran as a service, and accessed via the server's HTTP API (TODO). This allows you an easy way to scale up as BadActor's memory footprint is tiny. Because it leverages a light-weight cache with sharding and reaping, it allows most organizations to be confident that BadActor will not be a bottleneck. 

# Benchmarks 

| Type    |  Value   |
| --- | --- |
| **Model Name** | MacBook Pro |
| **Model Identifier** | MacBookPro11,3 |
| **Processor Name** | Intel Core i7 | 
| **Processor Speed** | 2.3 GHz | 
| **Number of Processors** | 1 |
| **Total Number of Cores** | 4 |
| **L2 Cache (per Core)** | 256 KB | 
| **L3 Cache** | 6 MB | 
| **Memory** | 16 GB |

###### 1.8.2015

```bash
➜  badactor git:(master) ✗ go test -bench=. -cpu=4 -benchmem -benchtime=5s | column -t
PASS
BenchmarkIsJailed-4                50000000                          121       ns/op  0    B/op  0  allocs/op
BenchmarkIsJailedFor-4             50000000                          134       ns/op  0    B/op  0  allocs/op
BenchmarkInfraction-4              5000000                           1390      ns/op  528  B/op  7  allocs/op
BenchmarkInfractionlIsJailed-4     3000000                           2755      ns/op  800  B/op  9  allocs/op
BenchmarkInfractionlIsJailedFor-4  3000000                           2733      ns/op  800  B/op  9  allocs/op
BenchmarkStudioInfraction512-4     3000000                           2215      ns/op  591  B/op  9  allocs/op
BenchmarkStudioInfraction1024-4    3000000                           2357      ns/op  612  B/op  9  allocs/op
BenchmarkStudioInfraction2048-4    5000000                           2617      ns/op  621  B/op  9  allocs/op
BenchmarkStudioInfraction4096-4    5000000                           2566      ns/op  671  B/op  9  allocs/op
BenchmarkStudioInfraction65536-4   3000000                           3309      ns/op  667  B/op  9  allocs/op
BenchmarkStudioInfraction262144-4  2000000                           3644      ns/op  674  B/op  9  allocs/op
ok                                 github.com/jaredfolkins/badactor  178.239s
➜  badactor git:(master) ✗
```

###### 12.30.2014

```bash
➜  badactor git:(master) ✗ go test -benchtime=5s -bench=. -benchmem -cpu=4 | column -t
PASS
BenchmarkIsJailed-4                  50000000                          133        ns/op  0          B/op  0        allocs/op
BenchmarkIsJailedFor-4               50000000                          136        ns/op  0          B/op  0        allocs/op
BenchmarkInfraction-4                10000000                          824        ns/op  116        B/op  5        allocs/op
BenchmarkInfractionMostCostly-4      10000000                          891        ns/op  116        B/op  5        allocs/op
BenchmarkInfractionIsJailed-4        3000000                           2569       ns/op  340        B/op  13       allocs/op
BenchmarkInfractionIsJailedFor-4     3000000                           2611       ns/op  340        B/op  13       allocs/op
Benchmark10000Actors1Infraction-4    1000                              8571335    ns/op  1162931    B/op  50023    allocs/op
Benchmark100000Actors1Infraction-4   100                               87687224   ns/op  11630938   B/op  500248   allocs/op
Benchmark1000000Actors1Infraction-4  10                                841989544  ns/op  116292788  B/op  5002740  allocs/op
Benchmark10000Actors4Infractions-4   200                               30728688   ns/op  4522659    B/op  170013   allocs/op
ok                                   github.com/jaredfolkins/badactor  93.868s
➜  badactor git:(master) ✗

```

###### 12.24.2014

```bash
➜  badactor git:(master) ✗ go test -bench=. -benchtime=5s -benchmem | column -t
PASS
BenchmarkIsJailed                 50000000                          138        ns/op  0         B/op  0       allocs/op
BenchmarkIsJailedFor              50000000                          140        ns/op  0         B/op  0       allocs/op
BenchmarkInfraction               10000000                          943        ns/op  128       B/op  4       allocs/op
BenchmarkInfractionMostCostly     10000000                          1008       ns/op  128       B/op  4       allocs/op
Benchmark10000Actors              100                               140566388  ns/op  13150354  B/op  150598  allocs/op
Benchmark10000Actors4Infractions  50                                241030802  ns/op  17278074  B/op  210614  allocs/op
ok                                github.com/jaredfolkins/badactor  73.592s
➜  badactor git:(master) ✗

```

###### 12.16.2014

This was **before** a serious refactoring. I am keeping it here because **(a)** I'd like to encourage others to *benchmark* their code and **(b)** I learned many valuable lessons while doing it. 
 
```bash
➜  badactor git:(master) go test -bench=. -benchtime=5s -benchmem | column -t
PASS
BenchmarkInfraction1                  2000                              2679694   ns/op  518  B/op  10  allocs/op
BenchmarkInfraction10                 2000                              3050845   ns/op  516  B/op  10  allocs/op
BenchmarkInfraction100                2000                              3430051   ns/op  516  B/op  10  allocs/op
BenchmarkInfraction1000               2000                              3738125   ns/op  516  B/op  10  allocs/op
BenchmarkInfraction10000              2000                              4004534   ns/op  516  B/op  10  allocs/op
BenchmarkInfractionWithIsJailed1      3000                              1832770   ns/op  193  B/op  3   allocs/op
BenchmarkInfractionWithIsJailed10     3000                              1968030   ns/op  193  B/op  3   allocs/op
BenchmarkInfractionWithIsJailed100    3000                              2120179   ns/op  193  B/op  3   allocs/op
BenchmarkInfractionWithIsJailed1000   3000                              1955656   ns/op  193  B/op  3   allocs/op
BenchmarkInfractionWithIsJailed10000  3000                              1943728   ns/op  193  B/op  3   allocs/op
ok                                    github.com/jaredfolkins/badactor  109.879s
➜  badactor git:(master)
```

# Httprouter & Negroni Example

```go
package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/jaredfolkins/badactor"
	"github.com/julienschmidt/httprouter"
)

var st *badactor.Studio

func main() {

	//runtime.GOMAXPROCS(4)

	// studio capacity
	var sc int32
	// director capacity
	var dc int32

	sc = 1024
	dc = 1024

	// init new Studio
	st = badactor.NewStudio(sc)

	// define and add the rule to the stack
	ru := &badactor.Rule{
		Name:        "Login",
		Message:     "You have failed to login too many times",
		StrikeLimit: 10,
		ExpireBase:  time.Second * 1,
		Sentence:    time.Second * 10,
	}
	st.AddRule(ru)

	err := st.CreateDirectors(dc)
	if err != nil {
		log.Fatal(err)
	}

	// Start the reaper
	st.StartReaper()

	// router
	router := httprouter.New()
	router.POST("/login", LoginHandler)

	// middleware
	n := negroni.Classic()
	n.Use(NewBadActorMiddleware())
	n.UseHandler(router)
	n.Run(":9999")

}

//
// HANDLER
//

// this is a niave login function for example purposes
func LoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var err error

	un := r.FormValue("username")
	pw := r.FormValue("password")

	// snag the IP for use as the actor's name
	an, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		panic(err)
	}

	// mock authentication
	if un == "example_user" && pw == "example_pass" {
		http.Redirect(w, r, "", http.StatusOK)
		return
	}

	// auth fails, increment infraction
	err = st.Infraction(an, "Login")
	if err != nil {
		log.Printf("[%v] has err %v", an, err)
	}

	// auth fails, increment infraction
	i, err := st.Strikes(an, "Login")
	log.Printf("[%v] has %v Strikes %v", an, i, err)

	http.Redirect(w, r, "", http.StatusUnauthorized)
	return
}

//
// MIDDLEWARE
//
type BadActorMiddleware struct {
	negroni.Handler
}

func NewBadActorMiddleware() *BadActorMiddleware {
	return &BadActorMiddleware{}
}

func (bam *BadActorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// snag the IP for use as the actor's name
	an, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		panic(err)
	}

	// if the Actor is jailed, send them StatusUnauthorized
	if st.IsJailed(an) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// call the next middleware in the chain
	next(w, r)
}
```
