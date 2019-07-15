# dummyDll

Export some functions. See if they load somewhere.

## Building

Depending on your host OS:

```sh
make {windows,linux}
```

Check the functions work:

```sh
.\rundll32.exe dummy.dll,DllRegisterServer
```

Place the dll somewhere you think it might hijack a hosting executable.

![](https://i.imgur.com/CJ6tx4K.png))
