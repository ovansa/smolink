package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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

func (suite *URLControllerTestSuite) TestShortenURL_Success() {
	payload := map[string]string{"url": "https://example.com"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.app.Router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var resp map[string]string
	err := json.NewDecoder(w.Body).Decode(&resp)
	suite.NoError(err)
	suite.Equal("https://example.com", resp["original_url"])
	suite.NotEmpty(resp["short_code"])
}

func (suite *URLControllerTestSuite) TestResolveURL_Success() {
	// Shorten first
	payload := map[string]string{"url": "https://golang.org"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.app.Router.ServeHTTP(w, req)

	var resp map[string]string
	_ = json.NewDecoder(w.Body).Decode(&resp)

	// Resolve
	code := resp["short_code"]
	req = httptest.NewRequest("GET", "/"+code, nil)
	req.Header.Set("User-Agent", "Go-Test")
	w = httptest.NewRecorder()
	suite.app.Router.ServeHTTP(w, req)

	suite.Equal(http.StatusFound, w.Code)
	suite.Equal("https://golang.org", w.Header().Get("Location"))
}

func (suite *URLControllerTestSuite) TestShortenURL_InvalidPayload() {
	req := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.app.Router.ServeHTTP(w, req)
	suite.Equal(http.StatusBadRequest, w.Code)

	body, _ := io.ReadAll(w.Body)
	suite.Contains(string(body), "Invalid request payload")
}

func TestURLControllerTestSuite(t *testing.T) {
	suite.Run(t, new(URLControllerTestSuite))
}
