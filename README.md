# tsv2mdtbl

## Usage

```
Usage of tsv2mdtbl:
   tsv2mdtbl [<inputFileName>]
  -H	use header
  -d string
    	delimiter character (default "\t")
  -delimiter string
    	delimiter character (default "\t")
  -header
    	use header
```

## Install

```
go get -u github.com/watarukura/tsv2mdtbl
```

## Example

```
$ echo num1\tnum2\tnum3\n1\t2\t3 | tsv2mdtbl --header                                                                
| NUM1 | NUM2 | NUM3 |
|------|------|------|
|    1 |    2 |    3 |
```

| NUM1 | NUM2 | NUM3 |
|------|------|------|
|    1 |    2 |    3 |

```
$ cat testdata/TEST2.txt
"num1","num2"
"1","2"
"3","4
5"
$ tsv2mdtbl --header --delimiter "," testdata/TEST2.txt
| NUM1 | NUM2   |
|------|--------|
|    1 |      2 |
|    3 | 4<br>5 |
```

| NUM1 | NUM2   |
|------|--------|
|    1 |      2 |
|    3 | 4<br>5 |
