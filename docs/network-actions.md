# Network actions

## Listing
Lists all of the current networks for the user. You can do this with `nets list`. Takes either `--id` or `--subdomain` for the organization:

```
$ katapult networks list --subdomain debug-inc
Networks:
NAME                    ID                    
Public Network          netw_gVRkZdSKczfNg34P   
Public Network - NYC    netw_q0lBvtutvOjujgyO   
Public Network - AZP    netw_NONCdbcLHfrIeloe   
Virtual Networks:
NAME    ID                    
Testing vnet_OEzVM9GftFGIKelfD
```

