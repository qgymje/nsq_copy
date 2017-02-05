package nsqd

import (
	"crypto/md5"
	"hash/crc32"
	"io"
	"log"
	"os"
	"time"
)

type Options struct {
	ID               int64  `flag:"worker-id" cfg:"id"`
	HTTPAddress      string `flag:"http-address"`
	BroadcastAddress string

	// diskqueue options
	DataPath string `flag:"data-path"`

	// statics options
	StatsdAddress  string        `flag:"statsd-address"`
	StatsdPrefix   string        `flag:"statsd-prefix"`
	StatsdInterval time.Duration `flag:"statsd-interval" arg:"1s"`
	StatsdMemStats bool          `flag:"statsd-mem-stats"`

	MaxDeflateLevel int `flag:"max-deflate-level"`

	Logger Logger
}

func NewOptions() *Options {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	h := md5.New()
	io.WriteString(h, hostname)
	defaultID := int64(crc32.ChecksumIEEE(h.Sum(nil)) % 1024) // 取余数

	return &Options{
		ID:               defaultID,
		HTTPAddress:      "0.0.0.0:4151",
		BroadcastAddress: hostname,

		StatsdPrefix:    "nsq.%s",
		MaxDeflateLevel: 6,
		Logger:          log.New(os.Stderr, "[nsqd]", log.Ldate|log.Ltime|log.Lmicroseconds),
	}
}
