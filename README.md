# Large file reader

`emoji.txt` were combined from [Kaggle dataset](https://www.kaggle.com/praveengovi/emotions-dataset-for-nlp)

```bash
$ cat *.txt > emoji.txt
``` 

### First attempt
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

### Second attempt
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
then we can simulate by reading 1 GB file, with max 0.16 GB (160 MB) memory. If the usage exceed 160 MB, the pod will be killed.

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
```text
http://IP_ADDRESS/word?file=file-url.txt
```

#### Create larger text file
```bash
$ cat /usr/share/dict/words | sort -R | head -100000 > file.txt
$ cat *.txt > big.txt # repeat for 10 times until get 1.6GB txt file
```

#### create cloud resources
- create public bucket
- upload `big10.txt` to bucket
- create GKE cluster
> make sure to download Service Account file with appropriate permissions from cloud console
```bash
$ cd terraform
$ terraform init
$ terraform apply
```

#### Build container and deploy
> make sure you enable cloud build access to GKE at [setting](https://console.cloud.google.com/cloud-build/settings/service-account)
```bash
$ gcloud builds submit . \
   --substitutions SHORT_SHA=$(git rev-parse --short HEAD)
```

#### Get the IP address
> can get from cloud console or CLI
```bash
$ gcloud container clusters get-credentials CLUSTER_NAME --zone ZONE --project PROJECT_ID
$ kubectl get service
NAME                  TYPE           CLUSTER-IP   EXTERNAL-IP   PORT(S)          AGE
file-reader-service   LoadBalancer   10.3.254.7   34.87.25.50   3000:31465/TCP   24m
```

#### Test the API
```bash
$ curl 'http://34.87.25.50:3000/word?file=https://storage.googleapis.com/temp-read-large-file-bucket/big10.txt' | jq

  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   402  100   402    0     0      4      0  0:01:40  0:01:36  0:00:04   104
{
  "file": {
    "MiB": 1526,
    "name": "file-1601703943.txt",
    "size": 1600501760,
    "url": "https://storage.googleapis.com/temp-read-large-file-bucket/big10.txt"
  },
  "mem": {
    "alloc": "22 MiB",
    "numGC": "506",
    "sys": "71 MiB",
    "totalAlloc": "5562 MiB"
  },
  "timeTaken": "1m15.587011379s",
  "words": {
    "top5": [
      "i, 16497152",
      "feel, 7136768",
      "and, 6135296",
      "to, 5709312",
      "the, 5352960",
      "a, 3958784"
    ],
    "uniqueWords": 119960,
    "wordCount": 247142912
  }
}
```
Pod usage metrics
> After hit the API for 3 times, rest in between, to show cool down period

> Memory: Max usage:44.4 MiB during file processing

![usage](screenshot/pod-usage.png)


### Clean up 

```
cd terraform
terraform destroy
```
