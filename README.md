# Gocache

Gocache is a Golang library for caching.

## Installation

You need [golang 1.13](https://golang.org/dl/) to use gocache.


```go
 go get github.com/ahmetask/gocache/v2
```
## Example 
[Server](https://github.com/ahmetask/gocache-server-example)

[Client](https://github.com/ahmetask/gocache-client-example)

## Usage

```go
package main

import (
	"fmt"
	"github.com/ahmetask/gocache/v2"
	"time"
)

func main() {
	cache := gocache.NewCache(5*time.Minute, 5*time.Second)

	cache.Set("foo", "bar", gocache.Eternal)
	cache.Set("foo2", 1, 5*time.Second)

	res, success := cache.Get("foo")
	fmt.Println(res)
	fmt.Println(success)
}


```

## Cache Server Usage
```go
package main

import (
	"github.com/ahmetask/gocache/v2"
	"time"
)

func main() {
	cache := gocache.NewCache(5*time.Minute, 5*time.Second)

	serverConfig := gocache.ServerConfig{
		CachePtr: cache,
		Port:     "8080",
	}

	gocache.NewCacheServer(serverConfig)
}

```

## Cache Client Usage
```go
package main

import (
	"github.com/ahmetask/gocache/v2"
)

func main() {
	request := gocache.AddCacheRequest{Key: "key", Value: "any interface", Life: 4}
	gocache.SaveCache("serviceIp", request)

	var r interface{}
	gocache.GetCache("serviceIp", "key", &r)
}


```
## Worker Pool Usage
```go
package main

import (
	"github.com/ahmetask/gocache/v2"
	"time"
)

func main()  {
	cache := gocache.NewCache(5*time.Minute, 5*time.Second)
	cache.StartWorkerPoolWith(2)
	
	cache.AddBackground("key", "value",  gocache.Eternal)
}

```
## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
