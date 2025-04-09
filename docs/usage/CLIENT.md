## Client Test Usage

### Build and Run Client

```
go build -o cli //cli.exe for windows
```

```
cli send "Message"
```

For config during testing, 
change fetchDbTest()

### Build and Run Test Db and Node

```
cd test
```
```
set PORT=8080 && npm run dev
```