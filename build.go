package main

import "encoding/json"

var gitTag string    // populated with -ldflags
var gitBranch string // populated with -ldflags
var gitHash string   // populated with -ldflags
var gitDirty string  // populated with -ldflags

type BuildInfo struct {
	GitTag    string `json:"gitTag"`
	GitBranch string `json:"gitBranch"`
	GitHash   string `json:"gitHash"`
	GitDirty  bool   `json:"gitDirty"`
}

func (b *BuildInfo) toJson() string {
	jsonBytes, err := json.Marshal(b)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func GetBuildInfo() BuildInfo {

	info := BuildInfo{
		GitTag:    "unknown",
		GitBranch: "unknown",
		GitHash:   "unknown",
		GitDirty:  false,
	}

	if gitTag != "" {
		info.GitTag = gitTag
	}

	if gitBranch != "" {
		info.GitBranch = gitBranch
	}

	if gitHash != "" {
		info.GitHash = gitHash
	}

	if gitDirty == "true" {
		info.GitDirty = true
	}

	return info
}
