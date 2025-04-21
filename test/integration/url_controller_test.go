package integration

import (
	"io"
	"net/http"
	"smolink/internal/errors"
	"smolink/internal/routes"
	"smolink/test"
	"testing"

	"github.com/stretchr/testify/suite"
)

type URLControllerTestSuite struct {
	suite.Suite
	app *test.TestApp
}

func (suite *URLControllerTestSuite) SetupSuite() {
	suite.app = test.SetupTestApp()
}

func (suite *URLControllerTestSuite) TearDownSuite() {
	suite.app.Cleanup()
}

func (suite *URLControllerTestSuite) SetupTest() {
	suite.app.ResetState()
}

const shortenURLEndpoint = routes.APIPrefix + routes.ShortenURLPath

func (suite *URLControllerTestSuite) TestShortenURLWithoutCustomCode_Success() {
	originalURL := "https://golang.org"
	payload := map[string]string{"url": originalURL}

	w := test.CreateTestRequest(suite.T(), suite.app.Router, http.MethodPost, shortenURLEndpoint, payload, "")

	suite.Equal(http.StatusCreated, w.Code)

	var resp map[string]string
	test.ParseResponse(suite.T(), w, &resp)
	suite.Equal(originalURL, resp["original_url"])
	suite.NotEmpty(resp["short_code"])
}

func (suite *URLControllerTestSuite) TestShortenURLWithCustomCode_Success() {
	shortCode, originalURL := "golang", "https://golang.org"
	payload := map[string]string{"url": originalURL, "customCode": shortCode}
	w := test.CreateTestRequest(suite.T(), suite.app.Router, http.MethodPost, shortenURLEndpoint, payload, "")

	suite.Equal(http.StatusCreated, w.Code)

	var resp map[string]string
	test.ParseResponse(suite.T(), w, &resp)
	suite.Equal(originalURL, resp["original_url"])
	suite.Equal(shortCode, resp["short_code"])
}

func (suite *URLControllerTestSuite) TestShortenURLWithExistingCustomCode_Failure() {
	shortCode, originalURL := "golang", "https://golang.org"
	err := suite.app.SeedShortURL(shortCode, originalURL)
	suite.Require().NoError(err)

	payload := map[string]string{"url": "https://example.com", "customCode": shortCode}
	w := test.CreateTestRequest(suite.T(), suite.app.Router, http.MethodPost, shortenURLEndpoint, payload, "")

	suite.Equal(http.StatusConflict, w.Code)

	var resp map[string]string
	test.ParseResponse(suite.T(), w, &resp)
	suite.Equal(errors.ErrCodeInUse.Code, resp["code"])
	suite.Equal(errors.ErrCodeInUse.Message, resp["message"])
}

func (suite *URLControllerTestSuite) TestShortenURLWithInvalidURL_Failure() {
	shortCode, originalURL := "newShortCode", "{{randomURL}}"

	payload := map[string]string{"url": originalURL, "customCode": shortCode}
	w := test.CreateTestRequest(suite.T(), suite.app.Router, http.MethodPost, shortenURLEndpoint, payload, "")

	suite.Equal(http.StatusBadRequest, w.Code)

	var resp map[string]string
	test.ParseResponse(suite.T(), w, &resp)
	suite.Equal(errors.ErrInvalidURL.Code, resp["code"])
	suite.Equal(errors.ErrInvalidURL.Message, resp["message"])
}

func (suite *URLControllerTestSuite) TestShortenURL_InvalidPayload() {
	w := test.CreateTestRequest(suite.T(), suite.app.Router, http.MethodPost, shortenURLEndpoint, nil, "")
	suite.Equal(http.StatusBadRequest, w.Code)

	body, _ := io.ReadAll(w.Body)
	suite.Contains(string(body), "Invalid request payload")
}

func (suite *URLControllerTestSuite) TestResolveURL_Success() {
	shortCode, originalURL := "golang", "https://golang.org"
	err := suite.app.SeedShortURL(shortCode, originalURL)
	suite.Require().NoError(err)

	code := shortCode
	w := test.CreateTestRequest(suite.T(), suite.app.Router, http.MethodGet, shortenURLEndpoint+"/"+code, nil, "")

	suite.Equal(http.StatusFound, w.Code)
	suite.Equal(originalURL, w.Header().Get("Location"))
}

func (suite *URLControllerTestSuite) TestResolveURL_ShortCodeDoesNotExist_Fail() {
	code := "shortCodeThatDoesNotExist"
	w := test.CreateTestRequest(suite.T(), suite.app.Router, http.MethodGet, shortenURLEndpoint+"/"+code, nil, "")

	suite.Equal(http.StatusNotFound, w.Code)
	var resp map[string]string
	test.ParseResponse(suite.T(), w, &resp)
	suite.Equal(errors.ErrShortCodeNotFound.Code, resp["code"])
	suite.Equal(errors.ErrShortCodeNotFound.Message, resp["message"])
}

func TestURLControllerTestSuite(t *testing.T) {
	suite.Run(t, new(URLControllerTestSuite))
}
