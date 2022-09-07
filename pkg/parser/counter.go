package parser

import "fmt"

const (
	OUTPUT_DIVISOR = 1000
)

var (
	ipCounters      = make(map[string]uint64)
	pathCounters    = make(map[string]uint64)
	pathDuration    = make(map[string]uint64)
	pathDurationAvg = make(map[string]float64)
)

func countIp(ip string) {
	ipCounters[ip]++
}

func countPath(path string, respTime uint64) {
	pathCounters[path]++
	pathDuration[path] += respTime
}

func getTopIps(max int) map[string]uint64 {
	retMap := make(map[string]uint64)
	maxArr := make([]string, max)
	for ip, n := range ipCounters {
		var i int
		for i = 0; i < len(maxArr) && n > ipCounters[maxArr[i]]; i++ {
		}
		if i > 0 {
			if i > 1 {
				copy(maxArr, maxArr[1:i])
			}
			maxArr[i-1] = ip
		}
	}
	for i, ip := range maxArr {
		if maxArr[i] != "" {
			retMap[ip] = ipCounters[ip]
		}
	}
	return retMap
}

func getTopPaths(max int) map[string]string {
	calculatePathAverages()
	retMap := make(map[string]string)
	maxArr := make([]string, max)
	for path, avg := range pathDurationAvg {
		var i int
		for i = 0; i < len(maxArr) && avg > pathDurationAvg[maxArr[i]]; i++ {
		}
		if i > 0 {
			if i > 1 {
				copy(maxArr, maxArr[1:i])
			}
			maxArr[i-1] = path
		}
	}
	for i, ip := range maxArr {
		if maxArr[i] != "" {
			retMap[ip] = fmt.Sprintf("%.2f", pathDurationAvg[ip])
		}
	}
	return retMap
}

func calculatePathAverages() {
	if len(pathDurationAvg) == len(pathDuration) {
		return
	}
	for ip, n := range pathCounters {
		pathDurationAvg[ip] = float64(pathDuration[ip]) / float64(n) / float64(OUTPUT_DIVISOR)
	}
}
