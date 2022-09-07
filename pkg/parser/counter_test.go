package parser

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCountIp(t *testing.T) {
	for _, tc := range []struct {
		name string
		ips  []string
	}{
		{"one ip address", []string{"1.1.1.1"}},
		{"different ip addresses", []string{"1.1.1.1", "2.2.2.2"}},
		{"repeating ip addresses", []string{"1.1.1.1", "2.2.2.2", "1.1.1.1", "3.3.3.3", "4.4.4.4", "3.3.3.3", "1.1.1.1"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cleanGlobalCounters()
			for _, ip := range tc.ips {
				countIp(ip)
			}
			for _, ip := range tc.ips {
				expectedN := occurrencesInSlice(ip, tc.ips)
				if n, found := ipCounters[ip]; !found || n != expectedN {
					t.Fatalf("ip counters incorrect: found=%v, expectedN=%v, n=%v", found, expectedN, n)
				}
			}
		})
	}
}

func TestGetTopIps(t *testing.T) {
	for _, tc := range []struct {
		name string
		ips  []string
		max  int
	}{
		{"empty", nil, 0},
		{"non-empty, zero max", []string{"1.1.1.1", "2.2.2.2"}, 0},
		{"one ip address, max 1", []string{"1.1.1.1"}, 1},
		{"different ip addresses, max 1", []string{"1.1.1.1", "2.2.2.2"}, 1},
		{"different ip addresses, max same as slice size", []string{"1.1.1.1", "2.2.2.2"}, 2},
		{"different ip addresses, max bigger than size", []string{"1.1.1.1", "2.2.2.2"}, 5},
		{"repeating ip addresses, max 1", []string{"1.1.1.1", "2.2.2.2", "1.1.1.1", "3.3.3.3", "4.4.4.4", "3.3.3.3", "1.1.1.1"}, 1},
		{"repeating ip addresses, max less than distinct ips", []string{"1.1.1.1", "2.2.2.2", "1.1.1.1", "3.3.3.3", "4.4.4.4", "3.3.3.3", "1.1.1.1"}, 2},
		{"repeating ip addresses, max same as distinct ips", []string{"1.1.1.1", "2.2.2.2", "1.1.1.1", "3.3.3.3", "4.4.4.4", "3.3.3.3", "1.1.1.1"}, 4},
		{"repeating ip addresses, bigger than distinct ips", []string{"1.1.1.1", "2.2.2.2", "1.1.1.1", "3.3.3.3", "4.4.4.4", "3.3.3.3", "1.1.1.1"}, 5},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cleanGlobalCounters()
			for _, ip := range tc.ips {
				countIp(ip)
			}
			ipsMap := getTopIps(tc.max)
			expectedOccurMap := make(map[string]uint64)
			for _, ip := range tc.ips {
				expectedOccurMap[ip] = occurrencesInSlice(ip, tc.ips)
			}
			expectedSize := minInts(tc.max, len(expectedOccurMap))
			if len(ipsMap) < expectedSize {
				t.Fatalf("result ip map is of different size than requested/possible %v != %v", len(ipsMap), expectedSize)
			}
			for ip, n := range ipsMap {
				if n != expectedOccurMap[ip] {
					t.Fatalf("unexpected number of occurrences of an ip address: %v != %v", n, expectedOccurMap[ip])
				}
				ipsLarger := 0
				ipsSmaller := 0
				for ip2, n2 := range expectedOccurMap {
					if ip2 == ip {
						continue
					}
					if n2 > n {
						ipsLarger++
					} else {
						ipsSmaller++
					}
				}
				if ipsLarger > tc.max-1 {
					t.Fatalf("Too many ips with bigger occurrence: %v > %v", ipsLarger, tc.max-1)
				}
				if ipsSmaller < len(expectedOccurMap)-tc.max {
					t.Fatalf("Too few ips with smaller occurrence: %v < %v", ipsSmaller, len(expectedOccurMap)-tc.max)
				}
			}
		})
	}
}

type PathTest struct {
	path     string
	duration uint64
}

func TestCountPaths(t *testing.T) {
	for _, tc := range []struct {
		name  string
		paths []PathTest
	}{
		{"one path", []PathTest{{"/admin.php", 1000}}},
		{"distinct paths", []PathTest{{"/admin.php", 1300}, {"/api/load.php", 1500}}},
		{"repeating paths", []PathTest{{"/admin.php", 1000}, {"/api/load.php", 2000}, {"/admin.php", 1500}, {"/admin.php", 1800}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cleanGlobalCounters()
			for _, p := range tc.paths {
				countPath(p.path, p.duration)
			}
			for _, p := range tc.paths {
				expectedCount := occurrencesInPathTest(p.path, tc.paths)
				expectedDuration := totalPathTestDuration(p.path, tc.paths)
				if n, found := pathCounters[p.path]; !found || n != expectedCount {
					t.Fatalf("path counters incorrect: found=%v, expectedCount=%v, n=%v", found, expectedCount, n)
				}
				if d, found := pathDuration[p.path]; !found || d != expectedDuration {
					t.Fatalf("path counters incorrect: found=%v, expectedCount=%v, n=%v", found, expectedCount, d)
				}
			}
		})
	}
}

func TestGetTopPaths(t *testing.T) {
	for _, tc := range []struct {
		name  string
		paths []PathTest
		max   int
	}{
		{"empty", nil, 0},
		{"empty, non-zero max", nil, 2},
		{"non-empty, zero max", []PathTest{{"/obi/wan/kenobi", 1000}}, 0},
		{"distinct paths, max 1", []PathTest{{"/obi/wan/kenobi", 1900}, {"/master/yoda", 1300}, {"/han/solo", 1500}, {"/luke/skywalker", 1000}}, 1},
		{"distinct paths, max greater than size", []PathTest{{"/obi/wan/kenobi", 1900}, {"/master/yoda", 1300}, {"/han/solo", 1500}, {"/luke/skywalker", 1000}}, 5},
		{"distinct paths, max less than size", []PathTest{{"/obi/wan/kenobi", 1900}, {"/master/yoda", 1300}, {"/han/solo", 1500}, {"/luke/skywalker", 1000}}, 2},
		{"repeating paths, max 1", []PathTest{{"/han/solo", 1900}, {"/master/yoda", 1300}, {"/han/solo", 1500}, {"/luke/skywalker", 1000}, {"/luke/skywalker", 1800}, {"/luke/skywalker", 1600}}, 1},
		{"repeating paths, max greater than size", []PathTest{{"/han/solo", 1900}, {"/master/yoda", 1300}, {"/han/solo", 1500}, {"/luke/skywalker", 1000}, {"/luke/skywalker", 1800}, {"/luke/skywalker", 1600}}, 4},
		{"repeating paths, max less than size", []PathTest{{"/han/solo", 1900}, {"/master/yoda", 1300}, {"/han/solo", 1500}, {"/luke/skywalker", 1000}, {"/luke/skywalker", 1800}, {"/luke/skywalker", 1600}}, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cleanGlobalCounters()
			for _, p := range tc.paths {
				countPath(p.path, p.duration)
			}
			pathsMap := getTopPaths(tc.max)
			expectedOccurMap := make(map[string]uint64)
			expectedTotalDurationMap := make(map[string]uint64)
			expectedAvgDurationMap := make(map[string]string)
			for _, p := range tc.paths {
				expectedOccurMap[p.path] = occurrencesInPathTest(p.path, tc.paths)
				expectedTotalDurationMap[p.path] = totalPathTestDuration(p.path, tc.paths)
				expectedAvgDurationMap[p.path] = fmt.Sprintf("%0.2f", float64(expectedTotalDurationMap[p.path])/float64(expectedOccurMap[p.path])/float64(OUTPUT_DIVISOR))
			}
			expectedSize := minInts(tc.max, len(expectedOccurMap))
			if len(pathsMap) < expectedSize {
				t.Fatalf("result path map is of different size than requested/possible %v != %v", len(pathsMap), expectedSize)
			}
			for p, avg := range pathsMap {
				if avg != expectedAvgDurationMap[p] {
					t.Fatalf("unexpected number of occurrences of a path: %v != %v", avg, expectedAvgDurationMap[p])
				}
				avgF, _ := strconv.ParseFloat(avg, 64)
				pathsLarger := 0
				pathsSmaller := 0
				expectedAvgS := fmt.Sprintf("%.2f", pathDurationAvg[p])
				if avg != expectedAvgS {
					t.Fatalf("returned avg for path '%v' unequal to the computed one: '%v' != '%v'", p, avg, expectedAvgS)
				}
				for p2, a2 := range expectedAvgDurationMap {
					if p == p2 {
						continue
					}
					a2F, _ := strconv.ParseFloat(a2, 64)
					if a2F > avgF {
						pathsLarger++
					} else {
						pathsSmaller++
					}
				}
				if pathsLarger > tc.max-1 {
					t.Fatalf("Too many paths with bigger avg: %v > %v", pathsLarger, tc.max-1)
				}
				if pathsSmaller < len(expectedOccurMap)-tc.max {
					t.Fatalf("Too few paths with smaller avg: %v < %v", pathsSmaller, len(expectedOccurMap)-tc.max)
				}
			}
		})
	}
}

func occurrencesInSlice(k string, arr []string) uint64 {
	retval := uint64(0)
	for _, v := range arr {
		if k == v {
			retval++
		}
	}
	return retval
}

func cleanGlobalCounters() {
	ipCounters = make(map[string]uint64)
	pathCounters = make(map[string]uint64)
	pathDuration = make(map[string]uint64)
	pathDurationAvg = make(map[string]float64)
}
func occurrencesInPathTest(k string, pta []PathTest) uint64 {
	retval := uint64(0)
	for _, pt := range pta {
		if pt.path == k {
			retval++
		}
	}
	return retval
}

func totalPathTestDuration(k string, pta []PathTest) uint64 {
	retval := uint64(0)
	for _, pt := range pta {
		if pt.path == k {
			retval += pt.duration
		}
	}
	return retval
}

func minInts(a, b int) int {
	if a < b {
		return a
	}
	return b
}
