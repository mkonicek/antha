## Aim:



##Status
 


##Next steps:



##Execution instructions:


#### Get required repos

1. branch of antha-lang/antha :

```bash
cd $GOPATH/src/github.com/antha-lang/antha

git fetch 
git checkout origin/master
cd -
```


2.  branch of Synthace/elements


```bash
cd $GOPATH/src/github.com/Synthace/elements
git fetch
git checkout origin/master
cd -
```


3. (A) Pipetmaxdriver

```bash
cd $GOPATH/src/github.com/Synthace/PipetMaxDriver
git fetch
git checkout 87b4440863ca30181edef2279af9e93411f42885
```

Or

3. (B) CybioDriver

```bash
cd $GOPATH/src/github.com/Synthace/CybioXMLDriver
git fetch
git checkout 0838e814157228144927b3cda42238693a2bdc31
cd -
```

#### Run whenever any source code is changed  (e.g. plate definitions, antha element changes, liquid class changes)

4. Build 

```bash
cd $GOPATH/src/github.com/Synthace/elements
make
go install -v github.com/Synthace/elements/cmd/...
cd -
```

or

```
elementsbuild
```


#### Run when parameters or workflow is changed

5. run


PipetMax:


```bash
antharun --driver  go://github.com/Synthace/PipetMaxDriver/server
```


Cybio:


```bash
cd $GOPATH/src/github.com/Synthace/CybioXMLDriver/server
go build ./...
./server -machine felix
```


```bash
antharun --driver localhost:50051 --inputPlateType pcrplate_skirted
```

6. Rename output file

e.g.

```bash
mv generated.sqlite pipetmaxday1.sqlite
```

or 

```bash
mv cybio.xml felixday1.xml
```


