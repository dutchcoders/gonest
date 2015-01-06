gonest
=====

Golang NEST library.


## Features


## Sample
```
package main

import (
    "github.com/dutchcoders/gonest"
    "log"
)

func main() {
    var err error
    var nest gonest.Nest
    if nest, err := gonest.Connect("clientid", ""); err != nil {
        log.Panic(err)
    }

    if err = nest.Authorize("secret", "pincode"); err != nil {
        log.Panic(err)
    }
    
    var devices gonest.Devices
    if err = nest.Devices(&devices); err != nil {
        log.Panic(err)
    }

    for _, device := range devices {
        fmt.Print(device.Name)
    }
}

```

## References

See github.com/dutchcoders/nest/ for complete implementation of library.

## Contributions

Contributions are welcome.

## Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

## Copyright and license

Code and documentation copyright 2011-2014 Remco Verhoef.
Code released under [the MIT license](LICENSE).

