> [!IMPORTANT]
> This software is now being maintained on [Codeberg](https://codeberg.org/scip/pgidler/).

# pgidler

Create idle postgres connections for testing

## Introduction

This program start a number of  idle (or idle in transaction) sessions
on a postgres database for testing purposes. It uses go routines (like
threads) to execute independent clients. Each client executes a number
of harmless SELECT's  and then stops doing anything.  The clients then
hang either forever or until a specified timeout is reached.

This is an implementation of a tool described at [AWS Blog: Performance impact of idle PostgreSQL connections](https://aws.amazon.com/blogs/database/performance-impact-of-idle-postgresql-connections/).

## Building

Run  `make`  to compile.  You'll  need  Golang.  The repo  contains  a
pre-compiled binary for linux amd.

## Usage

`pgidler` offers two idle modes:

- **idle**: create normal idle sessions
- **idle in transaction**: create sessions hanging in a transaction

The following commandline options are available:
```
Usage of ./pgidler:
  -c, --client int        Number of concurrent users (default 500)
  -d, --database string   Database (default "postgres")
  -i, --idletransaction   Wether to stay in idle in transaction state
  -p, --password string   Password of the database user
  -P, --port int          TCP Port (default 5432)
  -s, --server string     Server (default "localhost")
  -t, --timeout int       Wether to stop the clients after N seconds
  -u, --user string       Database user (default "postgres")
```

# Report bugs

[Please open an issue](https://codeberg.org/scip/pgidler/issues). Thanks!

# License

This work is licensed under the terms of the General Public Licens
version 3.

# Author

Copyleft (c) 2024 Thomas von Dein
