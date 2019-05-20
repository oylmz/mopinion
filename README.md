# mopinion

[![Build Status](https://travis-ci.com/oylmz/mopinion.svg?branch=master)](https://travis-ci.com/oylmz/mopinion)
[![codecov](https://codecov.io/gh/oylmz/mopinion/branch/master/graph/badge.svg)](https://codecov.io/gh/oylmz/mopinion)

Go client library for mopinion. The aim is to make the communication easier with the [Mopinion API](https://developer.mopinion.com/api/). Go doc can be found at https://godoc.org/github.com/oylmz/mopinion

## Usage ##

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/oylmz/mopinion"
)

func main() {
	basicCredentialProvider := mopinion.NewBasicCredentialProvider(os.Getenv("MOPINION_PUBLIC_KEY"),
		os.Getenv("MOPINION_PRIVATE_KEY"))
	client, _ := mopinion.NewClient(nil, basicCredentialProvider)

	ctx := context.TODO()
	var err error
	if _, _, err = client.Token.Get(ctx); err != nil {
		log.Fatalf("get the token: %s", err)
	}

	account, _, err := client.Account.Get(ctx)
	if err != nil {
		log.Fatalf("get account: %s", err)
	}

	log.Println("account: %+v", account)
}
```

## Contributing ##
Please feel free to contribute if any updates or changes happen in the Mopinion API.

## License ##

This library is distributed under the BSD-style license found in the [LICENSE](./LICENSE)
file.
