# Large file reader

`emoji.txt` were combined from [Kaggle dataset](https://www.kaggle.com/praveengovi/emotions-dataset-for-nlp)

```bash
cat *.txt > emoji.txt
``` 
Create more text file
```bash
cat /usr/share/dict/words | sort -R | head -100000 > file.txt
cat *.txt > big.txt # repeat for 10 times until get 1.6GB txt file
```

### first approach
- faster, lower memory usage, but not accurate
```bash
$ go run cmd/reader.go

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
$ go run cmd/scanner.go

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

Looks like second approach is better. Let's improve second approach further. To simulate low memory environment, we containerise it and make run it on http server.

Assume we need to read 100 GB file, with max 16 GB memory. 
then we can simulate by reading 1 GB file, with max 0.16 GB (160 MB) memory.

```yaml
 resources:
    requests:
      memory: "32Mi"
      cpu: "100m"
    limits:
      memory: "160Mi"
      cpu: "500m"
```

The API will look like this
```bash
http://IP_ADDRESS/word?file=file-url.txt
```

#### create cloud resources
- create public bucket
- upload `big10.txt` to bucket
- create GKE cluster
```bash
cd terraform
terraform init
terraform apply
```
#### Build container and deploy
```bash
gcloud builds submit . \
   --substitutions SHORT_SHA=$(git rev-parse --short HEAD)
```

