package gcs

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UtilTestSuite struct {
	suite.Suite
}

func (s *UtilTestSuite) SetupTest() {
	omitArg := func(argName string, args []string) []string {
		for i := range args {
			splitArgs := strings.Split(args[i], "=")
			if len(splitArgs) != 2 {
				continue
			}

			name := splitArgs[0]
			if name == argName {
				return append(args[:i], args[i+1:]...)
			}
		}

		return args
	}

	os.Args = omitArg("--gcs_bucket", os.Args)
}

func (s *UtilTestSuite) TestParsePathWithGSPathSyntax() {
	bucket, key, err := parsePath("gs://bucket-34ff28a4-b823-4c05-8367-8d760927938b/test_file.pdf")
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), bucket, "bucket-34ff28a4-b823-4c05-8367-8d760927938b")
		assert.Equal(s.T(), key, "test_file.pdf")
	}
}

func (s *UtilTestSuite) TestParsePathWithBucketAndKey() {
	bucket, key, err := parsePath("bucket-34ff28a4-b823-4c05-8367-8d760927938b/test_file.pdf")
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), bucket, "bucket-34ff28a4-b823-4c05-8367-8d760927938b")
		assert.Equal(s.T(), key, "test_file.pdf")
	}
}

func (s *UtilTestSuite) TestParsePathWithNestedBucketAndKey() {
	bucket, key, err := parsePath("bucket-34ff28a4-b823-4c05-8367-8d760927938b/certificates/2024/test_file.pdf")
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), bucket, "bucket-34ff28a4-b823-4c05-8367-8d760927938b")
		assert.Equal(s.T(), key, "certificates/2024/test_file.pdf")
	}
}

func (s *UtilTestSuite) TestParsePathWithNestedBucketAndKey2() {
	os.Args = append(os.Args, "--gcs_bucket=bucket-34ff28a4-b823-4c05-8367-8d760927938b")

	bucket, key, err := parsePath("/bucket-34ff28a4-b823-4c05-8367-8d760927938b/tbebj9shg-u03a0kd3p44-bc1d7e42ae4d-512.jpeg")
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), bucket, "bucket-34ff28a4-b823-4c05-8367-8d760927938b")
		assert.Equal(s.T(), key, "tbebj9shg-u03a0kd3p44-bc1d7e42ae4d-512.jpeg")
	}
}

func (s *UtilTestSuite) TestParsePathThrowsErrorForInvalidPaths() {
	_, _, err := parsePath("test_file.pdf")
	assert.EqualError(s.T(), err, ErrInvalidPath.Error())
}

func (s *UtilTestSuite) TestParsePathAppendsGCSBucketIfSet() {
	os.Args = append(os.Args, "--gcs_bucket=my-bucket")

	bucket, path, err := parsePath("test_file.pdf")
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), bucket, "my-bucket")
		assert.Equal(s.T(), path, "test_file.pdf")
	}
}

func TestUtils(t *testing.T) {
	suite.Run(t, new(UtilTestSuite))
}
