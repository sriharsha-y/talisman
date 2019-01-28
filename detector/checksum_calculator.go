package detector

import (
	"fmt"
	"os"
	"talisman/git_repo"
	"talisman/utilities"

	yaml "gopkg.in/yaml.v2"
)

//ChecksumCalculator represents a calculator for checksum
type ChecksumCalculator struct{}

//NewChecksumCalculator returns new instance of the ChecksumCompare
func NewChecksumCalculator() *ChecksumCalculator {
	cs := ChecksumCalculator{}
	return &cs
}

//FilterIgnoresBasedOnChecksums filters the file ignores from the TalismanRCIgnore which doesn't have any checksum value or having mismatched checksum value from the .talsimanrc
func (cc *ChecksumCalculator) FilterIgnoresBasedOnChecksums(additions []git_repo.Addition, ignoreConfig TalismanRCIgnore) TalismanRCIgnore {
	finalIgnores := []FileIgnoreConfig{}
	wd, _ := os.Getwd()
	repo := git_repo.RepoLocatedAt(wd)
	gitTrackedFilesAsAdditions := repo.TrackedFilesAsAdditions()
	gitTrackedFilesAsAdditions = append(gitTrackedFilesAsAdditions, additions...)
	for _, ignore := range ignoreConfig.FileIgnoreConfig {
		currentCollectiveChecksum := cc.calculateCollectiveChecksumForPattern(ignore.FileName, gitTrackedFilesAsAdditions)
		// Compare with previous checksum from FileIgnoreConfig
		if ignore.Checksum == currentCollectiveChecksum {
			finalIgnores = append(finalIgnores, ignore)
		}
	}
	rc := TalismanRCIgnore{}
	rc.FileIgnoreConfig = finalIgnores
	return rc
}

//SuggestTalismanRCForPatterns returns .talismanrc format suggestion with respective checksums for input patterns
func (cc *ChecksumCalculator) SuggestTalismanRCForPatterns(fileNamePatterns []string) string {
	wd, _ := os.Getwd()
	repo := git_repo.RepoLocatedAt(wd)
	gitTrackedFilesAsAdditions := repo.TrackedFilesAsAdditions()
	gitTrackedFilesAsAdditions = append(gitTrackedFilesAsAdditions, repo.StagedAdditions()...)
	var fileIgnoreConfigs []FileIgnoreConfig
	for _, pattern := range fileNamePatterns {
		collectiveChecksum := cc.calculateCollectiveChecksumForPattern(pattern, gitTrackedFilesAsAdditions)
		if collectiveChecksum != "" {
			fileIgnoreConfig := FileIgnoreConfig{pattern, collectiveChecksum, []string{}}
			fileIgnoreConfigs = append(fileIgnoreConfigs, fileIgnoreConfig)
		}
	}
	if len(fileIgnoreConfigs) == 0 {
		return ""
	}
	result := fmt.Sprintf("\n\x1b[33mFormat for .talismanrc for given file names\x1b[0m\n")
	talismanRcIgnoreConfig := TalismanRCIgnore{fileIgnoreConfigs}
	m, _ := yaml.Marshal(&talismanRcIgnoreConfig)
	return result + string(m)
}

func (cc *ChecksumCalculator) calculateCollectiveChecksumForPattern(fileNamePattern string, additions []git_repo.Addition) string {
	var patternpaths []string
	for _, addition := range additions {
		if addition.Matches(fileNamePattern) {
			patternpaths = append(patternpaths, string(addition.Path))
		}
	}
	// Calculate current collective checksum
	patternpaths = utilities.UniqueItems(patternpaths)
	if len(patternpaths) == 0 {
		return ""
	}
	currentCollectiveChecksum := utilities.CalculateCollectiveSHA256ForPaths(patternpaths)
	return currentCollectiveChecksum
}
