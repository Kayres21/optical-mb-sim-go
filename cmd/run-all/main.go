// Command run-all fans out multiple simulador config files across parallel
// OS processes (one process per config), since the simulador binary itself
// is not safe to run concurrently in-process (it mutates global stdout/stderr).
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"
)

type fragOptions struct {
	Dir    string
	Stride int
}

type jobResult struct {
	Config   string
	Duration time.Duration
	Err      error
	LogPath  string
}

func main() {
	configsDir := flag.String("configs-dir", "configs", "Directory to recursively search for config JSON files")
	binPath := flag.String("bin", "bin/simulador", "Path to the simulador binary")
	jobs := flag.Int("jobs", runtime.NumCPU(), "Maximum number of concurrent simulation processes")
	build := flag.Bool("build", true, "Build the simulador binary before running")
	runLogsDir := flag.String("run-logs-dir", "logs/run-all", "Directory to store per-process stdout/stderr captures")
	logs := flag.Bool("logs", true, "Pass -logs flag through to each simulation")
	extraArgs := flag.String("args", "", "Extra raw flags appended to every simulador invocation, e.g. \"-defrag-mode=before_arrival\"")
	fragDir := flag.String("frag-dir", "", "Directory for full-resolution fragmentation streams (one file set per config); disabled when empty")
	fragStride := flag.Int("frag-stride", 1, "Record one fragmentation sample every N arrivals (used with -frag-dir)")
	flag.Parse()

	if *jobs < 1 {
		*jobs = 1
	}

	if *build {
		if err := buildBinary(*binPath); err != nil {
			log.Fatalf("Failed to build simulador: %v", err)
		}
	}

	configPaths, err := findConfigs(*configsDir)
	if err != nil {
		log.Fatalf("Failed to find config files: %v", err)
	}
	if len(configPaths) == 0 {
		log.Fatalf("No config files found under %s", *configsDir)
	}

	if err := os.MkdirAll(*runLogsDir, 0o755); err != nil {
		log.Fatalf("Failed to create run-logs directory: %v", err)
	}

	fmt.Printf("Running %d configs with up to %d concurrent processes...\n", len(configPaths), *jobs)

	if *fragDir != "" {
		if err := os.MkdirAll(*fragDir, 0o755); err != nil {
			log.Fatalf("Failed to create frag-dir: %v", err)
		}
	}
	frag := fragOptions{Dir: *fragDir, Stride: *fragStride}

	results := runAll(configPaths, *binPath, *runLogsDir, *logs, *extraArgs, frag, *jobs)

	failed := printSummary(results)
	if failed > 0 {
		os.Exit(1)
	}
}

func buildBinary(binPath string) error {
	fmt.Println("Building simulador...")
	if err := os.MkdirAll(filepath.Dir(binPath), 0o755); err != nil {
		return err
	}
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func findConfigs(dir string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".json" {
			paths = append(paths, path)
		}
		return nil
	})
	sort.Strings(paths)
	return paths, err
}

func runAll(configPaths []string, binPath, runLogsDir string, logsEnabled bool, extraArgs string, frag fragOptions, jobs int) []jobResult {
	results := make([]jobResult, len(configPaths))
	sem := make(chan struct{}, jobs)
	var wg sync.WaitGroup

	for i, configPath := range configPaths {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, configPath string) {
			defer wg.Done()
			defer func() { <-sem }()
			results[i] = runOne(binPath, configPath, runLogsDir, logsEnabled, extraArgs, frag)
		}(i, configPath)
	}

	wg.Wait()
	return results
}

func runOne(binPath, configPath, runLogsDir string, logsEnabled bool, extraArgs string, frag fragOptions) jobResult {
	start := time.Now()

	name := sanitizeName(configPath)
	logPath := filepath.Join(runLogsDir, fmt.Sprintf("%s_%s.log", name, start.Format("20060102_150405")))

	logFile, err := os.Create(logPath)
	if err != nil {
		return jobResult{Config: configPath, Duration: time.Since(start), Err: err}
	}
	defer logFile.Close()

	args := []string{"-config", configPath, fmt.Sprintf("-logs=%t", logsEnabled)}
	if frag.Dir != "" {
		fragPrefix := filepath.Join(frag.Dir, fmt.Sprintf("%s_%s", name, start.Format("20060102_150405")))
		args = append(args, "-frag-out="+fragPrefix, fmt.Sprintf("-frag-stride=%d", frag.Stride))
	}
	if extraArgs != "" {
		args = append(args, extraArgs)
	}

	cmd := exec.Command(binPath, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	fmt.Printf("[start] %s -> %s\n", configPath, logPath)
	err = cmd.Run()
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("[fail]  %s (%s) — see %s\n", configPath, duration.Round(time.Millisecond), logPath)
	} else {
		fmt.Printf("[done]  %s (%s)\n", configPath, duration.Round(time.Millisecond))
	}

	return jobResult{Config: configPath, Duration: duration, Err: err, LogPath: logPath}
}

func sanitizeName(configPath string) string {
	trimmed := configPath
	if abs, err := filepath.Abs(configPath); err == nil {
		if rel, err := filepath.Rel(".", abs); err == nil {
			trimmed = rel
		}
	}
	trimmed = trimmed[:len(trimmed)-len(filepath.Ext(trimmed))]
	out := make([]rune, 0, len(trimmed))
	for _, r := range trimmed {
		if r == os.PathSeparator || r == '/' {
			out = append(out, '_')
		} else {
			out = append(out, r)
		}
	}
	return string(out)
}

func printSummary(results []jobResult) int {
	failed := 0
	var total time.Duration
	fmt.Println("\n=== Summary ===")
	for _, r := range results {
		total += r.Duration
		if r.Err != nil {
			failed++
			fmt.Printf("FAIL  %-60s %-10s %v\n", r.Config, r.Duration.Round(time.Millisecond), r.Err)
		} else {
			fmt.Printf("OK    %-60s %-10s\n", r.Config, r.Duration.Round(time.Millisecond))
		}
	}
	fmt.Printf("\n%d/%d succeeded, %d failed. Total process time: %s\n",
		len(results)-failed, len(results), failed, total.Round(time.Millisecond))
	return failed
}
