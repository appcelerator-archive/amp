[![Build Status](https://travis-ci.org/holys/safe.svg)](https://travis-ci.org/holys/safe)
[![Coverage Status](https://coveralls.io/repos/holys/safe/badge.svg?branch=master&service=github)](https://coveralls.io/github/holys/safe?branch=master)
[![GoDoc](https://godoc.org/github.com/holys/safe?status.svg)](https://godoc.org/github.com/holys/safe)

# Safe

Is your password safe? This is a Golang fork of [Safe](https://github.com/lepture/safe)


## How it works?

[link](https://github.com/lepture/safe#how-it-works)

## Installation 

```
go get github.com/holys/safe
```


## Usage 

```
import  "github.com/holys/safe"

s := safe.New(8, 0, 3, safe.Strong)
s.SetWords("/path/to/words")
s.Check("password")
```



## Performance

```
$ go test  -bench .
PASS
BenchmarkIsAsdf-4          	 5000000	       290 ns/op
BenchmarkIsByStep-4        	20000000	        74.4 ns/op
BenchmarkIsCommonPassword-4	100000000	        21.2 ns/op
BenchmarkReverse-4         	 3000000	       484 ns/op
BenchmarkCheck-4           	  500000	      2379 ns/op
ok  	github.com/holys/safe	8.670s
```


## License

MIT




