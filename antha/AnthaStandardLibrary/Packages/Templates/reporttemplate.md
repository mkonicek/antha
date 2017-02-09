## Aim:



##Status
 


##Next steps:



##Execution instructions:


#### Get required repos

1. branch of antha-lang/antha :

```bash
cd $GOPATH/src/github.com/antha-lang/antha

git fetch 
git checkout ***ANTHACOMMIT****
cd -
```


2.  branch of antha-lang/elements


```bash
cd $GOPATH/src/github.com/antha-lang/elements
git fetch
git checkout ***ELEMENTSCOMMIT****
cd -
```

3. Other Dependencies:

```bash
***OTHERDEPENDENCIES***
```

4. (A) Pipetmaxdriver

```bash
cd $GOPATH/src/github.com/Synthace/PipetMaxDriver
git fetch
git checkout ***PIPETMAXDRIVERCOMMIT****
```

Or

4. (B) CybioDriver

```bash
cd $GOPATH/src/github.com/Synthace/CybioXMLDriver
git fetch
git checkout ***CYBIODRIVERCOMMIT****
cd -
```

#### Run whenever any source code is changed  (e.g. plate definitions, antha element changes, liquid class changes)

5. Build 

```bash
make current -C $GOPATH/src/github.com/antha-lang/elements
```

or

```
anthabuild
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


