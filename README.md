# Go BankID
This library makes BankID authentications in Golang a breeze.

## Requirements

- For Production usage, a pfx certificate issued by a certified Bank.

## Quickstart

```
go get github.com/fouxarn/go-bankid
```

## Example
~~~ go
package main

import (
  "bufio"
  "fmt"
  "os"

  bankid "github.com/fouxarn/go-bankid"
)

func main() {
  client, _ := bankid.NewTestClient()
  
  u := &bankid.EndUserInfo{
    UserInfoType: "IP_ADDR",
    Value:        "192.168.0.1",
  }
  authResp, _ := s.Authenticate("190101010593", u)
  
  reader := bufio.NewScanner(os.Stdin)
  fmt.Println("Please open your bankid-app and verify authentication request")
  scanner.Scan()
  
  collResp, _ := s.Collect(authResp.OrderRef)
  
  if collResp.Status == bankid.StatusComplete {
    fmt.Printf("%v is now authenticated!\n", collResp.UserInfo.GivenName)
  } else {
    fmt.Printf("Failed authentication!\n")
  }
}
~~~
