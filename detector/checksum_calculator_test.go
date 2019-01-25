package detector

import (
	"talisman/git_repo"
	"testing"

	"github.com/stretchr/testify/assert"
)

var talismanRCWithInCorrectChecksum = `
fileignoreconfig:
- filename : 'some_file.pem'
  checksum : '25bd31a28bf9d4e06327f1c4a5cab2260574ae508803f66adcc393350e994866'
  ignore_detectors : []
`

var talismanRCWithCorrectChecksum = `
fileignoreconfig:
- filename : 'some_file.pem'
  checksum : '7b9e45a4f384e9547940d51eecced4c1421dc10b581a15291b125ed0dff3b89c'
  ignore_detectors : []
`

var talismanRCWithWildCardAndCorrectChecksum = `
fileignoreconfig:
- filename : '*.pem'
  checksum : '7b9e45a4f384e9547940d51eecced4c1421dc10b581a15291b125ed0dff3b89c'
  ignore_detectors : []
`

var talismanRCWithOneCorrectChecksum = `
fileignoreconfig:
- filename : 'some_file.pem'
  checksum : '87139cc4d975333b25b6275f97680604add51b84eb8f4a3b9dcbbc652e6f27ac'
  ignore_detectors : []
- filename : 'test/some_other_file.pem'
  checksum : '87139cc4d975333b25b6275f97680604add51b84eb8f4a3b9dcbbc652e6f27ac'
  ignore_detectors : []
`

func TestShouldConsiderBothFilesForDetection(t *testing.T) {
	cc := NewChecksumCalculator()
	rc := NewTalismanRCIgnore([]byte(talismanRCWithInCorrectChecksum))
	addition1 := git_repo.NewAddition("some_file.pem", make([]byte, 0))
	addition2 := git_repo.NewAddition("test/some_file.pem", make([]byte, 0))

	filteredRC := cc.FilterIgnoresBasedOnChecksums([]git_repo.Addition{addition1, addition2}, rc)

	assert.Len(t, filteredRC.FileIgnoreConfig, 0, "Should return empty ignores and detectors should scan both files")
}

func TestShouldNotConsiderBothFilesForDetection(t *testing.T) {
	cc := NewChecksumCalculator()
	rc := NewTalismanRCIgnore([]byte(talismanRCWithCorrectChecksum))
	addition1 := git_repo.NewAddition("some_file.pem", make([]byte, 0))
	addition2 := git_repo.NewAddition("test/some_file.pem", make([]byte, 0))

	filteredRC := cc.FilterIgnoresBasedOnChecksums([]git_repo.Addition{addition1, addition2}, rc)

	assert.Len(t, filteredRC.FileIgnoreConfig, 1, "Should return 1 ignores which detectors should honor")
}

func TestShouldConsiderOneFileForDetection(t *testing.T) {
	cs := NewChecksumCalculator()
	rc := NewTalismanRCIgnore([]byte(talismanRCWithOneCorrectChecksum))
	addition1 := git_repo.NewAddition("some_file.pem", make([]byte, 0))
	addition2 := git_repo.NewAddition("test/some_other_file.pem", make([]byte, 0))

	filteredRC := cs.FilterIgnoresBasedOnChecksums([]git_repo.Addition{addition1, addition2}, rc)

	assert.Len(t, filteredRC.FileIgnoreConfig, 1, "Should return 1 ignore and detectors should scan that file")
}

func TestShouldConsiderBothFilesForDetectionIfTalismanRCIsEmpty(t *testing.T) {
	cs := NewChecksumCalculator()
	rc := NewTalismanRCIgnore([]byte{})
	addition1 := git_repo.NewAddition("some_file.pem", make([]byte, 0))
	addition2 := git_repo.NewAddition("test/some_file.pem", make([]byte, 0))

	filteredRC := cs.FilterIgnoresBasedOnChecksums([]git_repo.Addition{addition1, addition2}, rc)

	assert.Len(t, filteredRC.FileIgnoreConfig, 0, "Should return empty ignores and detectors should scan both files")
}

func TestShouldConsiderWildCardPatternForDetection(t *testing.T) {
	cs := NewChecksumCalculator()
	rc := NewTalismanRCIgnore([]byte(talismanRCWithWildCardAndCorrectChecksum))
	addition1 := git_repo.NewAddition("some_file.pem", make([]byte, 0))
	addition2 := git_repo.NewAddition("test/some_file.pem", make([]byte, 0))

	filteredRC := cs.FilterIgnoresBasedOnChecksums([]git_repo.Addition{addition1, addition2}, rc)

	assert.Len(t, filteredRC.FileIgnoreConfig, 1, "Should return 1 ignores which detectors should honor")
}
