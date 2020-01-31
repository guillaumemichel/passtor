# Passtör

Authors: Cedric Maire, Guillaume Michel, Xavier Pantet

## Files

- ### [client/](client)

This folder contains the client implementation. To launch the client, just type:
```
go build
./client
````

- ### [server/](server)

This folder contains the passtor (dht node) implementation. The server should be used as follow:

```
Usage of ./server:
  -addr string
    	address used to communicate other passtors instances (default "127.0.0.1:5000")
  -name string
    	name of the Passtor instance
  -peers string
    	bootstrap peer addresses
  -v int
    	verbose mode (default 1)
```

To test our program it is recommended to run at least 5 passtor instances.

- ### [scripts/](scripts)

Contains scripts that automatically launch multiple instances of passtors.

- ### [tests/](tests)

Contains a few tests we created for the project.

- ### System constants

All system constants are defined in [const.go](const.go).

## How to run this project

To run this project you can run the following code:

```
git clone https://gitlab.gnugen.ch/gmichel/passtor
cd passtor
go build
```

Then you need to launch multiple [passtor](server/server.go) instances to start the system. Once it is done, you can connect to any passtor instance with the [client](client/client.go), and create an account, store, modify and download your credentials.
