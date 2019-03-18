package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version

	// Sort releases in an increasing order
	semver.Sort(releases)
	// Reverse releases
	for l, r := 0, len(releases)-1; l < r; l, r = l+1, r-1 {
		releases[l], releases[r] = releases[r], releases[l]
	}

	for _, release := range releases {
		if !minVersion.LessThan(*release) {
			break
		}
		// If a version whose major and minor version is the same as release is already appended to versionSlice, skip this release
		if len(versionSlice) != 0 &&
			versionSlice[len(versionSlice)-1].Major == release.Major &&
			versionSlice[len(versionSlice)-1].Minor == release.Minor {
			// Skip this version
			continue
		}
		// If a release is not a pre-release, append the release
		if release.PreRelease == "" {
			versionSlice = append(versionSlice, release)
		}
	}
	return versionSlice
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Read a file
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}

	isFirstLine := true
	for scanner.Scan() {
		// Skip the first line
		if isFirstLine {
			isFirstLine = false
			continue
		}

		line := strings.Split(scanner.Text(), ",")
		repository := strings.Split(line[0], "/")
		releases, _, err := client.Repositories.ListReleases(ctx, repository[0], repository[1], opt)
		// If error, skip this repository
		if err != nil {
			fmt.Println(err)
			continue
		}
		minVersion := semver.New(line[1])
		allReleases := make([]*semver.Version, len(releases))
		for i, release := range releases {
			versionString := *release.TagName
			if versionString[0] == 'v' {
				versionString = versionString[1:]
			}
			allReleases[i] = semver.New(versionString)
		}
		versionSlice := LatestVersions(allReleases, minVersion)

		fmt.Printf("latest versions of %s/%s: %s", repository[0], repository[1], versionSlice)
	}
}
