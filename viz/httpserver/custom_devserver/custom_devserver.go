package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/bazelbuild/rules_typescript/devserver/concatjs"
	"github.com/bazelbuild/rules_typescript/devserver/devserver"
	"github.com/bazelbuild/rules_typescript/devserver/runfiles"
	"github.com/golang/glog"
)

var (
	port = flag.Int("port", 5432, "server port to listen on")
	// The "base" CLI flag is only kept because within Google3 because removing would be a breaking change due to
	// ConcatJS and "devserver/devserver.go" still respecting the specified base flag.
	base            = flag.String("base", "", "server base (required, runfiles of the binary)")
	pkgs            = flag.String("packages", "", "root package(s) to serve, comma-separated")
	manifest        = flag.String("manifest", "", "sources manifest (.MF)")
	scriptsManifest = flag.String("scripts_manifest", "", "preScripts manifest (.MF)")
	servingPath     = flag.String("serving_path", "/_/ts_scripts.js", "path to serve the combined sources at")
	entryModule     = flag.String("entry_module", "", "entry module name")
	strictLoading   = flag.Bool("strict_loading", false, "fail if manifest specifies a bad path")
)

func main() {
	flag.Parse()

	//correctedBase := filepath.Join(*base, "..")
	correctedBase := *base //filepath.Join(*base)

	if len(*pkgs) == 0 || (*manifest == "") || (*scriptsManifest == "") {
		fmt.Fprintf(os.Stderr, "Required argument not set\n")
		os.Exit(1)
	}

	{
		baseAbs, err := filepath.Abs(*base)
		if err != nil {
			glog.Fatalf("failed to get absolute path to --base %q: %v", *base, err)
		}

		glog.Infof("Using runfiles at \n  base %q (abs = %q),\n  scripts_manifest %q\n  manifest %q ", *base, baseAbs, *scriptsManifest, *manifest)
		glog.Infof("RUNFILES_MANIFEST_FILE: %q", os.Getenv("RUNFILES_MANIFEST_FILE"))
	}
	manifestPath, err := runfiles.Runfile(correctedBase, *scriptsManifest)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find scripts_manifest in runfiles: %v\n", err)
		os.Exit(1)
	}

	scriptFiles, err := manifestFiles(manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read scripts_manifest: %v\n", err)
		os.Exit(1)
	}

	if !strings.HasPrefix(*servingPath, "/") {
		fmt.Fprintf(os.Stderr, "The specified serving_path does not start with a slash. "+
			"This causes the serving path to not have any effect.\n")
		os.Exit(1)
	}

	preScripts := make([]string, 0, 100)
	postScripts := make([]string, 0, 1)

	// Include the livereload script if IBAZEL_LIVERELOAD_URL is set.
	livereloadUrl := os.Getenv("IBAZEL_LIVERELOAD_URL")
	if livereloadUrl != "" {
		fmt.Printf("Serving livereload script from %s\n", livereloadUrl)
		livereloadLoaderSnippet := fmt.Sprintf(`(function(){
	const script = document.createElement('script');
	script.src = "%s";
	document.head.appendChild(script);
})();`, livereloadUrl)
		preScripts = append(preScripts, livereloadLoaderSnippet)
	}

	// Include the profiler script if IBAZEL_PROFILER_URL is set.
	profilerScriptURL := os.Getenv("IBAZEL_PROFILER_URL")
	if profilerScriptURL != "" {
		fmt.Printf("Serving profiler script from %s\n", profilerScriptURL)
		profilerLoaderSnippet := fmt.Sprintf(`(function(){
	const script = document.createElement('script');
	script.src = "%s";
	document.head.appendChild(script);
})();`, profilerScriptURL)
		preScripts = append(preScripts, profilerLoaderSnippet)
	}

	correctScriptFilePath := func(scriptPath string) string {
		return filepath.Join("..", scriptPath)
	}

	// Include all user scripts in preScripts. This should always include
	// the requirejs script which is added to scriptFiles by the devserver
	// skylark rule.
	for _, scriptFile := range scriptFiles {
		scriptFile = correctScriptFilePath(scriptFile)
		runfile, err := runfiles.Runfile(correctedBase, scriptFile)
		if err != nil {
			glog.Errorf("Could not find runfile %s, got error %s", scriptFile, err)
			if *strictLoading {
				glog.Exitf("terminiting due to inability to read runfile %q", scriptFile)
			}
		}

		js, err := loadScript(runfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read script %s: %v\n", scriptFile, err)
		} else {
			preScripts = append(preScripts, js)
		}
	}

	// If the entryModule is set then add a snippet to load
	// the application to postScripts to be outputted after the sources
	if *entryModule != "" {
		postScripts = append(postScripts, fmt.Sprintf("require([\"%s\"]);", *entryModule))
	}

	http.Handle(*servingPath, concatjs.ServeConcatenatedJS(*manifest, correctedBase, preScripts, postScripts,
		&RunfileFileSystem{}))
	pkgList := strings.Split(*pkgs, ",")
	http.HandleFunc("/", devserver.CreateFileHandler(*servingPath, *manifest, pkgList, correctedBase))

	h, err := os.Hostname()
	if err != nil {
		h = "localhost"
	}
	// Detect if we are running in a linux container inside ChromeOS
	// If so, we assume you want to use the native browser (outside the container)
	// so you'll need to modify the hostname to access the server
	if _, err := os.Stat("/etc/apt/sources.list.d/cros.list"); err == nil {
		h = h + ".linux.test"
	}

	fmt.Printf("Server listening on http://%s:%d/\n", h, *port)
	fmt.Fprintln(os.Stderr, http.ListenAndServe(fmt.Sprintf(":%d", *port), nil).Error())
	os.Exit(1)
}

func loadScript(path string) (string, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("// %s\n%s", path, buf), nil
}

// manifestFiles parses a manifest, returning a list of the files in the manifest.
func manifestFiles(manifest string) ([]string, error) {
	f, err := os.Open(manifest)
	if err != nil {
		return nil, fmt.Errorf("could not read manifest %s: %s", manifest, err)
	}
	defer f.Close()
	return manifestFilesFromReader(f)
}

// manifestFilesFromReader is a helper for manifestFiles, split out for testing.
func manifestFilesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		path := s.Text()
		if path == "" {
			continue
		}
		lines = append(lines, path)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// RunfileFileSystem implements FileSystem type from concatjs.
type RunfileFileSystem struct{}

// StatMtime gets the filestamp for the last file modification.
func (fs *RunfileFileSystem) StatMtime(filename string) (time.Time, error) {
	s, err := os.Stat(filename)
	if err != nil {
		return time.Time{}, err
	}
	return s.ModTime(), nil
}

// ReadFile reads a file given its file name
func (fs *RunfileFileSystem) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// ResolvePath resolves the specified path within a given root using Bazel's runfile resolution.
// This is necessary because on Windows, runfiles are not symlinked and need to be
// resolved using the runfile manifest file.
func (fs *RunfileFileSystem) ResolvePath(root string, manifestFilePath string) (string, error) {
	runfileEntries, err := bazel.ListRunfiles()
	if err != nil {
		return "", err
	}
	// See https://github.com/bazelbuild/rules_nodejs/blob/master/packages/typescript/src/internal/devserver/ts_devserver.bzl#L59
	// All manifest files are relative to a workspace, so parse out the workspace name
	// and use this to look up the corresponding runfile.
	workspace, shortName, err := splitWorkspaceShortName(manifestFilePath)
	if err != nil {
		return "", fs.suggestIntendedRunfile(err, manifestFilePath)
	}

	knownWorkspaces := make(map[string]bool)
	for _, entry := range runfileEntries {
		knownWorkspaces[entry.Workspace] = true
		if entry.Workspace != workspace {
			continue
		}
		if entry.ShortPath != shortName {
			continue
		}
		return entry.Path, nil
	}
	if !knownWorkspaces[workspace] {
		return "", fs.suggestIntendedRunfile(fmt.Errorf("workspace %q not found in workspaces of known runfiles %v (short name %s)", workspace, knownWorkspaces, shortName), manifestFilePath)
	}
	return "", fs.suggestIntendedRunfile(fmt.Errorf("could not find runfile (%q, %q)", workspace, shortName), manifestFilePath)
}

func (fs *RunfileFileSystem) suggestIntendedRunfile(originalErr error, manifestFilePath string) error {
	return originalErr
}

func splitWorkspaceShortName(manifestPath string) (workspace, shortName string, err error) {
	manifestPath = path.Clean(manifestPath)
	parts := strings.SplitN(manifestPath, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("manifest file path should have form <workspace>/<short path>, found no slash: %q", manifestPath)
	}
	// if trimmed := strings.TrimPrefix(parts[1], "../"); trimmed != parts[1] {
	// 	return splitWorkspaceShortName(trimmed)
	// }
	return parts[0], parts[1], nil
}

func findClosestRunfileEntry(manifestPath string, runfileEntries []bazel.RunfileEntry) *bazel.RunfileEntry {
	var best *bazel.RunfileEntry
	longestSharedSuffix := ""
	for _, rfe := range runfileEntries {
		if suf := sharedSuffix(manifestPath, rfe.Path); len(suf) > len(longestSharedSuffix) {
			longestSharedSuffix = suf
			rfe := rfe
			best = &rfe
			if len(suf) == len(manifestPath) {
				break
			}
		}
	}
	if len(manifestPath) > 0 && (float64(len(longestSharedSuffix))/float64(len(manifestPath))) > .40 {
		return best
	}
	return nil
}

func sharedSuffix(a, b string) string {
	posA, posB := len(a), len(b)
	for {
		if posA == 0 || posB == 0 {
			break
		}
		if a[posA-1] != b[posB-1] {
			break
		}
		posA--
		posB--
	}
	return a[posA:]
}
