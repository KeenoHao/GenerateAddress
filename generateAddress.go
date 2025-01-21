package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/zlyuancn/zpbar"
	"log"
	"net"
	"os"
	"strings"
)

var (
	BitSet = make([]byte, 1<<30)
)

/**
 * @create by: Keeno
 * @description: 根据IPv6地址模式生成所有情况地址
 * @create time: 2024/7/30 13:46
 *
 */
func expandWildcard(input string) []string {
	var results []string
	expandWildcardHelper("", input, &results)
	return results
}

func expandWildcardHelper(prefix, remaining string, results *[]string) {
	if len(remaining) == 0 {
		if containsInBloomFilter(prefix) {
			return
		}
		*results = append(*results, prefix)
		return
	}

	firstWildcard := strings.Index(remaining, "*")
	if firstWildcard == -1 {
		// Keys of Bloom filters
		if containsInBloomFilter(prefix + remaining) {
			return
		}
		*results = append(*results, prefix+remaining)
		return
	}

	// Find the part before the first '*'
	partBeforeWildcard := remaining[:firstWildcard]
	remaining = remaining[firstWildcard+1:]

	// Generate all combinations for '*'
	for i := 0; i <= 15; i++ {
		hexValue := fmt.Sprintf("%01x", i)
		expandWildcardHelper(prefix+partBeforeWildcard+hexValue, remaining, results)
	}
}

func murmur3(data []byte, seed uint32) uint32 {
	hash := seed
	for i := 0; i < len(data); i = i + 4 {
		k := binary.BigEndian.Uint32(data[i : i+4])
		k = k * 0xcc9e2d51
		k = (k << 15) | (k >> 17)
		k = k * 0x1b873593
		hash = hash ^ k
		hash = (hash << 13) | (hash >> 19)
		hash = hash*5 + 0xe6546b64
	}
	hash = hash ^ (hash >> 16)
	hash = hash * 0x85ebca6b
	hash = hash ^ (hash >> 13)
	hash = hash * 0xc2b2ae35
	hash = hash ^ (hash >> 16)
	return hash
}

func containsInBloomFilter(address string) bool {
	IPv6 := net.ParseIP(address).To16()
	i := murmur3(IPv6, 0x12345678)
	j := murmur3(IPv6, 0x87654321)
	// Check if the ip is in BitSet
	if BitSet[i/8]&(1<<(i%8)) != 0 && BitSet[j/8]&(1<<(j%8)) != 0 {
		return true
	}
	BitSet[i/8] |= (1 << (i % 8))
	BitSet[j/8] |= (1 << (j % 8))
	return false
}

func readPatterns(patternFile string) []string {

	file, err := os.Open(patternFile)
	if err != nil {
		fmt.Println("打开文件时出错:", err)
		return nil
	}
	defer file.Close()

	// 创建一个新的扫描器用于读取文件
	scanner := bufio.NewScanner(file)
	var patterns []string
	// 逐行读取文件
	for scanner.Scan() {
		line := scanner.Text()
		patterns = append(patterns, line)
	}
	return patterns
}

func generateAddress(patternFile, addressFile string) {

	patterns := readPatterns(patternFile)
	if len(patterns) == 0 {
		return
	}

	outputFile, err := os.OpenFile(addressFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	writer := bufio.NewWriter(outputFile)
	if err != nil {
		log.Fatalf("无法创建输出文件: %v", err)
	}
	defer outputFile.Close()
	p := zpbar.NewPbar(
		zpbar.WithTotal(int64(len(patterns))),
	)
	p.Start()
	var countAddress float64
	for _, pattern := range patterns {
		addresses := expandWildcard(pattern)
		for _, ip := range addresses {
			countAddress++
			fmt.Fprintf(writer, "%s\n", ip)
		}
		p.Done()
	}
	writer.Flush()
	p.Close()
	println("地址数量为:", countAddress)
}
