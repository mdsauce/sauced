# Sauced
[![CircleCI](https://circleci.com/gh/mdsauce/sauced/tree/master.svg?style=svg)](https://circleci.com/gh/mdsauce/sauced/tree/master)
Managed Sauce Connect tunnels.

[Why](why.md)

### Install and Run
1. `go get github.com/mdsauce/sauced`
2.  clone the repo and put it in `$GOPATH/github.com/mdsauce/sauced`

Create your [config](https://github.com/mdsauce/sauced#config-file) file for sauce labs like below.  You must manually put the access key and username in, environment variables will not be read.

```
# '#' <- This is a comment and the line will be ignored.
# ~/.config/sauced.config
# main tunnel pool
/Users/maxdobeck/workspace/sauce_connect/sc-4.5.1-osx/bin/sc -u sauce.username -k sauce.access.key -v --no-remove-colliding-tunnels -N -i main-tunnel-pool --se-port 0 --pidfile /tmp/sc_client-1.pid 

/Users/maxdobeck/workspace/sauce_connect/sc-4.5.1-osx/bin/sc -u sauce.username -k sauce.access.key -v --no-remove-colliding-tunnels -N -i main-tunnel-pool --se-port 0 --pidfile /tmp/sc_client-1.pid 

/Users/maxdobeck/workspace/sauce_connect/sc-4.5.1-osx/bin/sc -u sauce.username -k sauce.access.key -v --no-remove-colliding-tunnels -N -i main-tunnel-pool --se-port 0 --pidfile /tmp/sc_client-1.pid 
```

Run it: 

`$ sauced start -c ~/.config/sauced.config`

### Config File
The config file should have one line for each SC instance.  The part of the line should be the full path to the SC binary you want to use.  The other arguments should be the flags you would use if you were starting the tunnel manually from a command line.  The scheme is be:

```
/path/to/bin/sc.exe <flags and arguments>
```

An example of a single tunnel:
`/home/myuser/tools/sc-4.5.3-linux/bin/sc -u account-name-here -k api-key-here -v`


An example of a pool of tunnels:

```
/home/user/tools/sc-4.5.1-linux/bin/sc -u account-name -k api-key-here -v --no-remove-colliding-tunnels -N -i main-tunnel-pool --se-port 0 
/home/user/tools/sc-4.5.1-linux/bin/sc -u account-name -k api-key-here -v --no-remove-colliding-tunnels -N -i main-tunnel-pool --se-port 0 
```

### Testing
Run `$ go test ./...`.

### Building for another platform
Run the `build.sh` script and find the binary for the OS of your choosing in the `builds/` directory. 
