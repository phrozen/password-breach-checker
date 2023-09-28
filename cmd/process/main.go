package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/phrozen/password-breach-checker/pkg/format"
)

const BUFFER_SIZE = 24 * 1024 * 1024 // 24MB

func Process(filename string, minCount int) error {
	// Start timer
	start := time.Now()
	// Open input file for scanning
	inputFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	log.Println("Opened input file:", filename)
	// Check size of input file
	stats, err := inputFile.Stat()
	if err != nil {
		return err
	}
	log.Println("Size:", format.Bytes(uint64(stats.Size())))
	// Create output file changing extension
	outputFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".bin"
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	log.Println("Created output file:", outputFilename)
	// Create new Scanner from the file
	scanner := bufio.NewScanner(inputFile)
	// Create a new buffer for reading
	inputBuffer := make([]byte, BUFFER_SIZE)
	// Assing buffer to our scanner
	scanner.Buffer(inputBuffer, BUFFER_SIZE)
	// Create a new buffer for writing
	outputBuffer := make([]byte, 0, BUFFER_SIZE)
	// Utility counters
	total, processed, written := 0, 0, 0
	// Separator to split hashes and counts
	separator := []byte(":")
	// Create a buffer to process each line
	// 20 bytes hash (SHA1 160 bits) + 4 bytes count (32bit int) = 24 bytes
	lineBuffer := make([]byte, 24)
	// For as long as we can scan...
	log.Println("Processing file... Buffer Size:", format.Bytes(BUFFER_SIZE))
	for scanner.Scan() {
		total++
		// Split line by : (separator)
		s := bytes.Split(scanner.Bytes(), separator)
		if len(s) != 2 {
			return fmt.Errorf("abnormal split of line: %s", scanner.Text())
		}
		// Decode hash into the first 20 bytes of the buffer
		n, err := hex.Decode(lineBuffer, s[0])
		if n != 20 {
			return fmt.Errorf("abnormal length of hash (%d): %s", n, string(s[0]))
		}
		if err != nil {
			return fmt.Errorf("error processing hash <%s>: %w", string(s[0]), err)
		}
		count, err := strconv.Atoi(string(s[1]))
		if err != nil {
			return fmt.Errorf("error processing count (%s): %w", string(s[1]), err)
		}
		// Skip hashes where count does not meet the threshold
		if count < minCount {
			continue
		}
		// Put the count as a 4 byte slice at the end of the line buffer [20:]
		binary.BigEndian.PutUint32(lineBuffer[20:], uint32(count))
		// Add line buffer to our output buffer
		outputBuffer = append(outputBuffer, lineBuffer...)
		// If output buffer is full, flush/write and reset its size
		if len(outputBuffer) == BUFFER_SIZE {
			n, err := outputFile.Write(outputBuffer)
			if err != nil {
				return fmt.Errorf("error writing buffer: %w", err)
			}
			written += n
			fmt.Printf("\r\t%s processed...     ", format.Bytes(uint64(written)))
			outputBuffer = outputBuffer[:0]
		}
		processed++
	}
	// Write the remainder of hashes if any are left in the buffer
	if len(outputBuffer) > 0 {
		n, err := outputFile.Write(outputBuffer)
		if err != nil {
			return fmt.Errorf("error writing buffer: %w", err)
		}
		written += n
		fmt.Printf("\r\t%s processed...     ", format.Bytes(uint64(written)))
	}
	fmt.Println(" Done!")
	log.Println("Hashes:", total)
	log.Println("Processed:", processed)
	log.Println("Elapsed:", time.Since(start))
	t := float64(stats.Size()) / float64(time.Since(start)) * 1000
	log.Printf("Throughput: %.2f MB/s\n", t)
	// error of the scanner if any
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

func main() {
	input := flag.String("f", "", "Input filename of passwords ordered by hash")
	count := flag.Int("c", 1, "Minimum count of breaches to process the hash")
	flag.Parse()

	if err := Process(*input, *count); err != nil {
		log.Fatalln(err)
	}
}
