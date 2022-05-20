# who-dis

**who-dis** is a DNS resolver that we built to better understand DNS, both from a user and name server perspective. Similar to **dig**, **who-dis** takes in a domain name and contacts name servers recursively, until it finds the records for the desired domain name, or bottoms out in some other way.

**who-dis** supports both recursive and non-recursive requests, result caching, and can decode A/AAAA/CNAME records, with extensibility to support more record types.

## Usage

```
./who-dis [-trace] [-no-cache] <domain-name>
```

- `-trace` enables a recursive DNS query. By default, **who-dis** simply queries 8.8.8.8. But with `-trace` set, **who-dis** will query a root server, followed by the name server for the top-level domain, and so on.

- `-nocache` disables the use of the cache for answering requests. By default **who-dis** will consult its cache before making a DNS query. The cache maps domain names to an IP address, a TTL, and the timestamp when the record was cached.

## Examples

Without setting `-trace`, the response from **who-dis** mimicks the behavior of `dig`:

```
$ ./who-dis www.brown.edu
nochache = false; trace = false
### Cached response ###
brown.edu	128.111.1.1

### Response from 128.111.1.1 ###
;; QUESTION SECTION:
www.brown.edu			IN	A

;; ANSWER SECTION:
www.brown.edu		3600	IN	CNAME	www.brown.edu.cdn.cloudflare.net

;; AUTHORITY SECTION:

;; ADDITIONAL SECTION:
```

A recursive query is similar to `dig +trace`:
```
$ ./who-dis -trace -nocache www.brown.edu
nochache = true; trace = true

### Response from 199.9.14.201 ###
;; QUESTION SECTION:
www.brown.edu			IN	A

;; ANSWER SECTION:

;; AUTHORITY SECTION:
edu		172800	IN	NS	b.edu-servers.net
edu		172800	IN	NS	h.edu-servers.net
edu		172800	IN	NS	i.edu-servers.net
edu		172800	IN	NS	d.edu-servers.net
edu		172800	IN	NS	k.edu-servers.net
edu		172800	IN	NS	a.edu-servers.net
edu		172800	IN	NS	e.edu-servers.net
edu		172800	IN	NS	f.edu-servers.net
edu		172800	IN	NS	c.edu-servers.net
edu		172800	IN	NS	l.edu-servers.net
edu		172800	IN	NS	g.edu-servers.net
edu		172800	IN	NS	j.edu-servers.net
edu		172800	IN	NS	m.edu-servers.net

;; ADDITIONAL SECTION:
m.edu-servers.net		172800	IN	A	192.55.83.30
l.edu-servers.net		172800	IN	A	192.41.162.30
k.edu-servers.net		172800	IN	A	192.52.178.30
j.edu-servers.net		172800	IN	A	192.48.79.30
i.edu-servers.net		172800	IN	A	192.43.172.30
h.edu-servers.net		172800	IN	A	192.54.112.30
g.edu-servers.net		172800	IN	A	192.42.93.30
f.edu-servers.net		172800	IN	A	192.35.51.30
e.edu-servers.net		172800	IN	A	192.12.94.30
d.edu-servers.net		172800	IN	A	192.31.80.30
c.edu-servers.net		172800	IN	A	192.26.92.30
b.edu-servers.net		172800	IN	A	192.33.14.30
a.edu-servers.net		172800	IN	A	192.5.6.30
m.edu-servers.net		172800	IN	AAAA	2001:501:b1f9::30

### Response from 192.55.83.30 ###
;; QUESTION SECTION:
www.brown.edu			IN	A

;; ANSWER SECTION:

;; AUTHORITY SECTION:
brown.edu		172800	IN	NS	ns1.ucsb.edu
brown.edu		172800	IN	NS	bru-ns1.brown.edu
brown.edu		172800	IN	NS	bru-ns2.brown.edu
brown.edu		172800	IN	NS	bru-ns3.brown.edu

;; ADDITIONAL SECTION:
ns1.ucsb.edu		172800	IN	A	128.111.1.1
ns1.ucsb.edu		172800	IN	AAAA	2607:f378::1
bru-ns1.brown.edu		172800	IN	A	128.148.248.11
bru-ns2.brown.edu		172800	IN	A	128.148.248.12
bru-ns3.brown.edu		172800	IN	A	128.148.2.13

### Response from 128.111.1.1 ###
;; QUESTION SECTION:
www.brown.edu			IN	A

;; ANSWER SECTION:
www.brown.edu		3600	IN	CNAME	www.brown.edu.cdn.cloudflare.net

;; AUTHORITY SECTION:

;; ADDITIONAL SECTION:
```

Performing the query again gives us a cache hit so we don't have to go back to the root server:

```
$ ./who-dis -trace www.brown.edu
nochache = false; trace = true
### Cached response ###
brown.edu	128.111.1.1

### Response from 128.111.1.1 ###
;; QUESTION SECTION:
www.brown.edu			IN	A

;; ANSWER SECTION:
www.brown.edu		3600	IN	CNAME	www.brown.edu.cdn.cloudflare.net

;; AUTHORITY SECTION:

;; ADDITIONAL SECTION:
```
