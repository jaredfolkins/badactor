![badactor logo](https://raw.githubusercontent.com/jaredfolkins/badactor_logo/master/badactor_logo_300x300.png) [![Build Status](https://travis-ci.org/jaredfolkins/badactor.svg?branch=master)](https://travis-ci.org/jaredfolkins/badactor) [![Coverage Status](https://img.shields.io/coveralls/jaredfolkins/badactor.svg)](https://coveralls.io/r/jaredfolkins/badactor?branch=master)

# BadActor 

BadActor is an application driven in-memory jailer, in the spirit of fail2ban. It can be used as either middleware or a standalone server with its primary goal to increase the expense for "bad actors" who engage in system probing or attacks.

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

BadActor can be included in your go application and ran concurrently. It can also be stood up on its own, ran as a service, and accessed via the server's HTTP API (TODO). This allows you an easy way to scale up as BadActor's memory footprint is tiny. Because it uitlizes Go's goroutines and channels, it allows most organizations to be confident that BadActor will not be a bottleneck. 

# Httprouter & Negroni Example

```go
package main

import (
  "log"
  "net/http"
  "time"

  "github.com/codegangsta/negroni"
  "github.com/jaredfolkins/badactor"
  "github.com/julienschmidt/httprouter"
  "gopkg.in/unrolled/render.v1"
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
    ExpireBase:  time.Second * 2,                                                                                                                                                                                                                                                                                                      Sentence:    time.Second * 2,
  }

  err := d.AddRule(ru)
  if err != nil {
    panic(err)
  }

  // run the director
  d.Run()

  // router
  router := httprouter.New()
  router.GET("/success", Success)
  router.GET("/fail", Fail)
  router.POST("/login", Login)

  // middleware
  n := negroni.Classic()
  n.UseHandler(router)
  n.Run(":9999")
}

func Success(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  ren := render.New(render.Options{IndentJSON: true})
  ren.JSON(w, http.StatusOK, map[string]string{"status": "success"})
  return
}

func Fail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  ren := render.New(render.Options{IndentJSON: true})
  ren.JSON(w, http.StatusOK, map[string]string{"status": "fail"})
  return
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

  // this is a niave login function for example purposes
  var err error
  un := r.FormValue("username")
  pw := r.FormValue("password")
  rn := "Login"

  if un == "example_user" && pw == "example_pass" {
    http.Redirect(w, r, "/success", 302)
    return
  }

  err = d.Infraction(un, rn)
  if err != nil {
    log.Printf("[%v] has err %v", un, err)
  }

  i, err := d.Strikes(un, rn)
  log.Printf("[%v] has [%d] strikes, err is %v", un, i, err)

  b := d.IsJailed(un)
  log.Printf("[%v] IsJailed = [%t]", un, b)

  b = d.IsJailedFor(un, rn)
  log.Printf("[%v] IsJailedFor = [%t] [%v]", un, b, rn)

  http.Redirect(w, r, "/fail", 302)
  return
}
```
