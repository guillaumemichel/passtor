# Scripts

Author: Guillaume

## [clean.sh](clean.sh)

Kill all running Passtor instances

## [launch.sh](launch.sh)

Starts multiple server instances connected with each other, forming a DHT. The addresses start at `127.0.0.1:6000` (udp for Passtor-Passtor connections, and tcp for Client-Passtor connections), and the port number increases (`6001`, `6002` ...). Verbose level is not a mandatory parameter. The logs of the instances are stored in the [logs](logs) folder. The processes never stop until you kill them manually or with [clean.sh](clean.sh).

```
Usage:      ./launch.sh [nb_hosts] [verbose_lvl]
Example:    ./launch.sh 10 2
```

## [start.sh](start.sh)

Only work on linux with tilix installed
