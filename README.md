# amzn-es-proxy

amzn-es-proxy is reverse proxy to add authorization header for Amazon Elastisearch Service.

## Installing

amzn-es-proxy uses [dep](https://github.com/golang/dep) to solve the depencency.

```bash
$ go get -d github.com/ikeisuke/amzn-es-proxy
$ dep ensure
$ go install
```

## Usage

```bash
$ amzn-es-proxy --domain=logs --region=ap-northeast-1 &
Using endpoint search-logs-xxxxxx.ap-northeast-1.es.amazonaws.com
Listen 127.0.0.1:9200
```

```bash
$ curl http://localhost:9200/
{
  "name" : "rSKb5Pf",
  "cluster_name" : "247280120152:logs",
  "cluster_uuid" : "xxxxxxxx",
  "version" : {
    "number" : "5.1.1",
    "build_hash" : "5395e21",
    "build_date" : "2016-12-15T22:47:19.049Z",
    "build_snapshot" : false,
    "lucene_version" : "6.3.0"
  },
  "tagline" : "You Know, for Search"
}
```
