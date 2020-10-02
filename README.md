# Large file reader

`emoji.txt` were combined from [Kaggle dataset](https://www.kaggle.com/praveengovi/emotions-dataset-for-nlp)

```bash
cat *.txt > emoji.txt
``` 

### first approach
- faster, lower memory usage, but not accurate
```bash
$ go run reader.go

4 goroutine has been completed
4 goroutine has been completed
4 goroutine has been completed
4 goroutine has been completed
Alloc = 2 MiB	TotalAlloc = 2 MiB	Sys = 69 MiB	NumGC = 0
emoji.txt is 2069616 bytes
uniqueWords: 10929, wordCount: 101011
time taken: 52.895164ms
top 5 words

i, 4124
feel, 3859
and, 3292
to, 3097
the, 2848
a, 2163
```

### second approach
- slower, higher memory usage, accurate
```bash
$ go run scanner.go

Alloc = 4 MiB	TotalAlloc = 10 MiB	Sys = 71 MiB	NumGC = 3
emoji.txt is 2069616 bytes
uniqueWords: 23929, wordCount: 382701
time taken: 35.120529ms
top 5 words

i, 32221
feel, 13938
and, 11983
to, 11151
the, 10454
a, 7732

```
