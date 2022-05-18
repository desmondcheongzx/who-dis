package pkg

import (
	"net"
	"time"

	bolt "github.com/boltdb/bolt"
)

type CachedRecord struct {
	timestamp uint32 // unix epoch time in seconds
	ttl       uint32
	addr      net.IP
}

func NewCachedRecord(rr *ResourceRecord) *CachedRecord {
	return &CachedRecord{
		timestamp: uint32(time.Now().Unix()),
		ttl:       rr.ttl,
		addr:      net.IP(rr.rdata),
	}
}

func (cr *CachedRecord) serialize() []byte {
	buf := make([]byte, 0)
	buf = append(buf, htonl(cr.timestamp)...)
	buf = append(buf, htonl(cr.ttl)...)
	buf = append(buf, cr.addr...)
	return buf
}

func (cr *CachedRecord) deserialize(data []byte) {
	cr.timestamp = ntohl(data[0:4])
	cr.ttl = ntohl(data[4:8])
	cr.addr = net.IP(data[8:])
}

type Cache struct {
	db   *bolt.DB
	path string
}

func NewCache(path string) *Cache {
	db, err := bolt.Open(path, 0666, &bolt.Options{
		Timeout: 2 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		getBucket(tx, []byte("record"))
		return nil
	})
	return &Cache{db: db, path: path}
}

func getBucket(tx *bolt.Tx, bucket []byte) *bolt.Bucket {
	var err error
	b := tx.Bucket(bucket)
	if b == nil {
		b, err = tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			panic(err)
		}
	}
	return b
}

func (cache *Cache) Store(name string, rr *CachedRecord) error {
	return cache.db.Update(func(tx *bolt.Tx) error {
		b := getBucket(tx, []byte("record"))
		data := rr.serialize()
		return b.Put([]byte(name), data)
	})
}

func (cache *Cache) Get(name string) (*CachedRecord, bool) {
	tx, err := cache.db.Begin(true)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()
	b := getBucket(tx, []byte("record"))
	data := b.Get([]byte(name))
	if data == nil {
		return nil, false
	}
	rr := &CachedRecord{}
	rr.deserialize(data)
	return rr, true
}
