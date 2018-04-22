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
  -redmine
    	output redmine table
```

## Install

```
go get -u github.com/watarukura/tsv2mdtbl
```

## Example

```
$ echo num1\tnum2\tnum3\n1\t2\t3 | tsv2mdtbl --header                                                                
| num1 | num2 | num3 |
| --- | --- | --- |
| 1 | 2 | 3 |
$ echo num1,num2,num3\n1,2,3 | tsv2mdtbl --delimiter=,
| num1 | num2 | num3 |
| 1 | 2 | 3 |
```
