## (a) antharun --driver

1. 
There are many additional flags which may be used with the antharun command. 
To see the full list type antharun --help on the command line.

2. 
To select which driver port to connect to add the --driver flag as shown above (making sure it matches the driver port you’ve served in a separate terminal). 

A driver can be called from a local port if already running (see part b):

```bash
antharun --driver localhost:50051
```

Or called directly from source code:

```bash
antharun --driver go://github.com/antha-lang/manualLiquidHandler/server
```

or

```bash
antharun --driver go://github.com/Synthace/PipetMaxDriver/server
```

### (b) Running a driver 
If running the pipetmax driver this will be launched locally from a binary using the following command in a separate terminal prior to using antharun:

```PipetMax```

By default this will set the port to ```localhost:50051```, output file to ```generated.sqlite``` and protocol name to ```rpctest```

There are various flags available to configure these defaults:



```PipetMax -out nameoffile.sqlite```
```PipetMax -port 50052```
```PipetMax -protocol newprotocolname```

If using the clientdevice (```./clientdevice.exe -—device pipetmax```) to control the pipetmax directly you can add an additional driver flag with the IP address of the remote clientdevice and port.
This will allow direct running and simulation of the Pipetmax:

```antharun --driver localhost:50051 --driver 192.168.1.58:50051```


or if from source code:



```antharun --driver  go://github.com/Synthace/PipetMaxDriver/server --driver 192.168.1.58:50051```



run ```ifconfig``` on pc controlling the pipetmax to find out IP address.



3. The manualLiquidhandlingdriver would work in the same way
You can get this from source code before running in a separate terminal or call directly as seen above

go get github.com/antha-lang/manualLiquidHandler

Running it:

```bash
cd server
go build ./...
./server
```



Again, the default port is 50051

## Excercises

1. Run the protocol with the pipetmax driver

## Next Steps

Now go to [Folder C](../C_differentcomponents/readme_flags.md) to find out about the other configuration options for antharun and a run down of the basics on LHComponents.