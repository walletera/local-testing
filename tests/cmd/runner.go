package main

import (
    "tests/scenarios"

    "github.com/form3tech-oss/f1/v2/pkg/f1"
)

func main() {
    // Create a new f1 instance, add all the scenarios and execute the f1 tool.
    // Any scenario added here can be executed like: `go run main.go run constant dinopayOutboundSucceed`
    f1.New().Add("dinopayOutboundSucceed", scenarios.DinopayOutboundSucceed).Execute()
}
