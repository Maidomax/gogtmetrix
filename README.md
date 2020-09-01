# gogtmetrix
A Golang SDK for GTmetrix

## Usage
Firtly, get an instance of an authenticated GTmetrix client with:
``` go
import (
    "github.com/Maidomax/gogtmetrix"
)

c := gogtmetrix.GetClient("my-username", "my-password")
```

After that the easiest way to get a site tested is to use:
``` go
tm, err := c.TestAndWaitForResults("https://mysite.com")

if err != nil {
    log.Println(err)
}

```

This function will automatically poll GTmetrix for you once a second until it returns the TestModel holding the results or the error. This function gives up after 5 minutes.

Alternatively, you can use other methods provided to handle polling yourself. Check the docs for more details. Make sure to stay within the constraints of your GTmetrix account in terms of API usage and concurent tests running.