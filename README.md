![badactor logo](https://raw.githubusercontent.com/jaredfolkins/badactor_logo/master/badactor_logo_300x300.png) [![Build Status](https://travis-ci.org/jaredfolkins/badactor.svg?branch=master)](https://travis-ci.org/jaredfolkins/badactor) [![Coverage Status](https://img.shields.io/coveralls/jaredfolkins/badactor.svg)](https://coveralls.io/r/jaredfolkins/badactor?branch=master)

# BadActor 

BadActor is an in-memory, application driven jailer built in the spirit of fail2ban. It can be used as either middleware or a standalone server (TODO) with its primary goal to increase the expense for "bad actors" who engage in system probing or attacks.

# Use Case

A common use case for BadActor is jailing an offender who fails to login to your website (N) times as this can signal a bruteforce attempt.

Egor Homakov's [otp calculator](http://sakurity.com/otp) provides a good description.

# Design

- observer pattern
- speed (subsecond response underload and submillisecond with standard operations)
- no external dependencies 
- concurrent, non blocking, self governed workers
- 100% code coverage and thorough tests
- accessable to multiple systems across the internet using an HTTP API (TODO)

# Does It Scale? 

BadActor can be included in your go application and ran concurrently. It can also be stood up on its own, ran as a service, and accessed via the server's HTTP API (TODO). This allows you an easy way to scale up as BadActor's memory footprint is tiny. Because it leverages Go's goroutines and channels, it allows most organizations to be confident that BadActor will not be a bottleneck. 

# Benchmarks 

Infraction() is the easiest but also most expensive function to use. You simply pass it an ActorName (IP address in our example) and the RuleName ("Login" in our example).

If we first wrap an Infraction() call around a call to IsJailed(), we can limit the expense of Infraction() because we know that our BadActor is already Jailed.

We benchmark both to show the performance gain.

###Model Name: 
MacBook Pro

###Model Identifier: 
MacBookPro11,3

###Processor Name: 
Intel Core i7

###Processor Speed:
2.3 GHz

###Number of Processors:
1

###Total Number of Cores:
4

###L2 Cache (per Core):
256 KB

###L3 Cache:
6 MB

###Memory:
16 GB

```bash
➜  badactor git:(master) go test -bench=. -benchtime=5s -benchmem 2>/dev/null | column -t
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

var d *badactor.Director

func main() {

  // create new director
  d = badactor.NewDirector()
  // create and add rule
    ru := &badactor.Rule{
    Name:        "Login",
    Message:     "You have failed to login too many times",
    StrikeLimit: 10,
    ExpireBase:  time.Second * 2, // if no activity is detected the infraction will expire after 2 seconds
    Sentence:    time.Minute * 5, // the sentence for breaking the rule is to be jailed for 5 minutes
  }

  // add the rule to the stack
  err := d.AddRule(ru)
  if err != nil {
    panic(err)
  }
  // run the director
  d.Run()

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
// MIDDLEWARE
//
type BadActorMiddleware struct {
  negroni.Handler
}

func NewBadActorMiddleware() *BadActorMiddleware {
  return &BadActorMiddleware{}
}

func (bam *BadActorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

  // snag the IP as the actor's name
  an, _, err := net.SplitHostPort(r.RemoteAddr)
  if err != nil {
    panic(err)
  }

  if d.IsJailed(an) {
    http.Redirect(w, r, "", http.StatusTeapot)
    return
  }

  // call the next middleware in the chain
  next(w, r)
}

//
// HANDLER
//
func LoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

  // this is a niave login function for example purposes
  var err error
  un := r.FormValue("username")
  pw := r.FormValue("password")
  rn := "Login"

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
  err = d.Infraction(an, rn)
  if err != nil {
    log.Printf("[%v] has err %v", an, err)
  }

  // unauthorized
  http.Redirect(w, r, "", http.StatusUnauthorized)
  return
}
```
