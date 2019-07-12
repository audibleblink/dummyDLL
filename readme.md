# dummyDll

Export some functions. See if they load somewhere.

## Compilation

```sh
GOOS=windows go build -o dummy.dll -buildmode=c-shared main.go
```

Check the functions work:

```sh
PS> .\rundll32.exe dummy.dll,DllRegisterServer
```
![](https://i.imgur.com/4dmX8lF.png)
