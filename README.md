# Ping Reboot

## Background
My Linux server in the garage on occasion loses network connectivity to the router. 

Even though the machine is hardwired to the router via ethernet cable, it stops responding to ICMP and SSH requests.

This little utility will ping a target host or IP until the remote host no longer responds. After a number of retry attempts, the reboot instruction will be sent to the underlying OS.

## Prerequisites
In order to build this library, you will need to have GIT and GO installed. In addition the machine needs to have access to the internet in order to download the additional libraries.

## How to use
1. Clone the repository.
2. Execute `make build` - this will build the pingreboot binary for the machine OS and architecture type
3. `cd bin` and run the tool `./pingreboot` and provide the required flags

## The command line
```bash
# Specify the target
pingreboot --target 1.1.1.1 

# Add interval in minutes between pings to target
pingreboot --target 1.1.1.1 --interval 2

# Specify the failure retry count before restart
pingreboot --target 1.1.1.1 --retry-count 5
`
```
```
