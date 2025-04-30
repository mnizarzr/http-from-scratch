package main

import (
	"bytes"
	"compress/gzip"
	"log"
)

func gzipCompress(text string) string {
	var b bytes.Buffer
	zw := gzip.NewWriter(&b)
	_, err := zw.Write([]byte(text))
	if err != nil {
		log.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}
	return b.String()
}

func supportedCompression(schemes []string) string {
	if len(schemes) == 0 {
		return ""
	}
	compressions := []string{"gzip"}
	comp := Intersection(schemes, compressions)
	if comp == nil {
		return ""
	}
	return comp[0]
}
