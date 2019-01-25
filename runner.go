package main

import (
	"fmt"
	"os"
	"talisman/detector"
	"talisman/git_repo"
)

const (
	//CompletedSuccessfully is an exit status that says that the current runners run completed without errors
	CompletedSuccessfully int = 0

	//CompletedWithErrors is an exit status that says that the current runners run completed with failures
	CompletedWithErrors int = 1
)

//Runner represents a single run of the validations for a given commit range
type Runner struct {
	additions []git_repo.Addition
	results   *detector.DetectionResults
}

//NewRunner returns a new Runner.
func NewRunner(additions []git_repo.Addition) *Runner {
	return &Runner{additions, detector.NewDetectionResults()}
}

//RunWithoutErrors will validate the commit range for errors and return either COMPLETED_SUCCESSFULLY or COMPLETED_WITH_ERRORS
func (r *Runner) RunWithoutErrors() int {
	r.doRun()
	r.printReport()
	return r.exitStatus()
}

//RunChecksumCalculator runs the checksum calculator against the patterns given as input
func (r *Runner) RunChecksumCalculator(fileNamePatterns []string) int {
	checksumCalculator := detector.NewChecksumCalculator()
	rcSuggestion := checksumCalculator.SuggestTalismanRCForPatterns(fileNamePatterns)
	if rcSuggestion != "" {
		fmt.Print(rcSuggestion)
	}
	return 0
}

func (r *Runner) doRun() {
	ignoresNew := detector.ReadConfigFromRCFile(readRepoFile())
	detector.DefaultChain().Test(r.additions, ignoresNew, r.results)
}

func (r *Runner) printReport() {
	if r.results.HasIgnores() || r.results.HasFailures() {
		fmt.Println(r.results.Report())
	}
}

func (r *Runner) exitStatus() int {
	if r.results.HasFailures() {
		return CompletedWithErrors
	}
	return CompletedSuccessfully
}

func readRepoFile() func(string) ([]byte, error) {
	wd, _ := os.Getwd()
	repo := git_repo.RepoLocatedAt(wd)
	return repo.ReadRepoFileOrNothing
}
