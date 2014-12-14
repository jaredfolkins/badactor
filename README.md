![badactor logo](https://raw.githubusercontent.com/jaredfolkins/badactor_logo/master/badactor_logo_300x300.png) [![Build Status](https://travis-ci.org/jaredfolkins/badactor.svg?branch=master)](https://travis-ci.org/jaredfolkins/badactor) [![Coverage Status](https://coveralls.io/repos/jaredfolkins/badactor/badge.png?branch=master)](https://coveralls.io/r/jaredfolkins/badactor?branch=master)

# BadActor 

BadActor is an application driven in-memory jailer, in the spirit of fail2ban. It can be used as either middleware or a standalone server with its primary goal to increase the expense for "bad actors" who engage in system probing or attacks.

# Use Case

A common use case for BadActor is jailing an offender who fails to login to your website (N) times as this can signal a bruteforce attempt.

Egor Homakov's [otp calculator](http://sakurity.com/otp)
provides a good description.

# Design

- observer pattern
- concurrent, non blocking, self governed workers
- accessable to multiple systems across the internet
- no external dependencies 
- 100% code coverage and thorough tests

# Does It Scale?

BadActor can be included in your go application and ran concurrently. It can also be stood up on its own, ran as a service, and accessed via the server's HTTP API (TODO). This allows you an easy way to scale up as BadActor's memory footprint is tiny. Because it uitlizes Go's goroutines and channels, it allows most organizations to be confident that BadActor will not be a bottleneck. 
