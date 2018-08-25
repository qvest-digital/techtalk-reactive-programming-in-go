package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tarent/techtalk-reactive-programming-in-go/crawl"
	"github.com/tarent/techtalk-reactive-programming-in-go/persistence"
)

func main() {
	
}

func getFileName(in string) string {
	return in[strings.Index(in, "://")+3:]
}
