## EasyStress

Send a lot of requests to test a website's stress resistance

**This program is just made for test stress resistance of websites, not for DDOS attack.**

> To use this program, you must agree these terms:

- All the Consequences made by using this program should be borne by you, not by me (the author of this program).

- You use this program means you agree to all the terms.

- If you don't agree to these terms, you MUST stop using this program AT ONCE.

### Install

You can download from Github Release or build by yourself:

```shell
# Linux
# If you don't have golang, run this. (Ubuntu)
# apt install golang
go build main.go
mv main easystress
# You can also move 'easystress' to a directory in PATH
```

```powershell
# Windows
# If you don't have golang, download from golang.google.cn
go build main.go
mv main.exe easystress.exe
# You can also move 'easystress.exe' to a directory in PATH
```

### Use

```shell
easystress -t 100 -w 5 -f t.csv https://example.com
```

This command will open 5 workers to send 100 requests to https://example.com and record to `t.csv`.

```
GLOBAL OPTIONS:
   --licence, -l             Show the licence (default: false)
   --time value, -t value    The time of sending requests
   --worker value, -w value  The amount of workers to send requests (default: 4)
   --file value, -f value    The name of the csv file which contains every request's time and error (default: none)
   --help, -h                show help
   --version, -v             print the version
```
