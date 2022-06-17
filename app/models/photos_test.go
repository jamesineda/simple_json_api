package models

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type PhotosTestSuite struct {
	suite.Suite
	Photos []string
}

func (suite *PhotosTestSuite) SetupTest() {
	suite.Photos = []string{
		"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",
	}
}

func (suite *PhotosTestSuite) Test_Unmarshall() {
	var photos Photos = suite.Photos
	mimes, _ := photos.GetFileTypes()
	suite.Equal([]string{"image/png", "image/jpeg"}, mimes)
}

func TestPhotosTestSuite(t *testing.T) {
	suite.Run(t, new(PhotosTestSuite))
}
