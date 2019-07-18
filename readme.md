# dummyDLL

Export some functions. See if they load somewhere.

## Building

* Have `mingw-w64`

```sh
make
```

Check the functions work:

```sh
rundll32.exe dummy.dll,DllRegisterServer
```

Place the dll somewhere you think it might hijack a hosting executable.

![](https://i.imgur.com/9fMiAQG.png))
