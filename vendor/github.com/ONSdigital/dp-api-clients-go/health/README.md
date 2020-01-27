health
=======

health package creates a generic health check client to determine an applications (not just APIs) ability to perform the required function.

```
import (
    ...
    "context"

    "github.com/ONSdigital/dp-api-clients-go/health"
    ...
)

func main() {
    ...
    // Create new health check (rchttp) client, this will set the '/health' 
    // and '/healthcheck' as endpoints that are not retriable
    hcClient := health.NewClient(<name>, <url>)

    ctx := context.Background()

    // Check state of external service
    checkObj, err := hcClient.Checker(ctx)
    if err != nil {
        ...
    }
    
    ...
}