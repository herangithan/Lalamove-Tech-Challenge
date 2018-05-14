package main

import (
	"context"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"

	"io/ioutil"
	"strings"
	"os"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	semver.Sort(releases)
	for index, release := range releases{
		if release.PreRelease != "" || release.Major<minVersion.Major || release.Minor<minVersion.Minor || release.Patch<minVersion.Patch {
			continue
		}
		if (index+1) <= len(releases)-1 && releases[index].Major==releases[index+1].Major && releases[index].Minor==releases[index+1].Minor{
			continue
		}
		versionSlice = append([]*semver.Version{release},versionSlice...)
	}
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
	return versionSlice
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}

	//File Reader
	if(len(os.Args)<2){
		fmt.Printf("Please Pass an argument for the text file")
		os.Exit(-1)
	}
	filename := os.Args[1]
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error Reading File %q", filename)
		os.Exit(-1)
	}

	//ReadText Files
	textfile:= string(dat)
	lines := strings.Split(textfile,"\r\n")
	for _,line := range lines {
		text := strings.Split(line, ",")
		if len(text) != 2{
			fmt.Printf("Error Reading Line %s, does not have the format owner/repository,min_version",line)
			continue
		}
		repository, min_version := text[0], text[1]
		//fmt.Printf("repo: %q and min version: %q", repository, min_version)
		owner, repo:= strings.Split(repository,"/")[0],strings.Split(repository,"/")[1]

		releases, _, err := client.Repositories.ListReleases(ctx, owner, repo, opt)
		if err != nil {
			fmt.Printf("Error while	Listing Releases: %q", err)
			os.Exit(-1)
			//panic(err) // is this really a good way?
		}
		minVersion := semver.New(min_version)
		allReleases := make([]*semver.Version, len(releases))
		for i, release := range releases {
			versionString := *release.TagName
			if versionString[0] == 'v' {
				versionString = versionString[1:]
			}
			allReleases[i] = semver.New(versionString)
		}
		versionSlice := LatestVersions(allReleases, minVersion)

		fmt.Printf("latest versions of %s/%s: %s\n",owner,repo, versionSlice)
	}
}
