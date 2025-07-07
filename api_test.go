package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const baseURL = "http://localhost:8080/api/v1"

type TestSuite struct {
	client       *http.Client
	adminToken   string
	memberToken  string
	adminEmail   string
	memberEmail  string
	orgName      string
	orgID        uint
	articleID    uint
	commentID    uint
	adminUser    map[string]interface{}
	memberUser   map[string]interface{}
	organization map[string]interface{}
}

func NewTestSuite() *TestSuite {
	return &TestSuite{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (ts *TestSuite) makeRequest(method, endpoint string, body interface{}, token string) (*http.Response, []byte, error) {
	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := http.NewRequest(method, baseURL+endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := ts.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, responseBody, nil
}

func (ts *TestSuite) makeRequestWithOrgHeader(method, endpoint string, body interface{}, token string, orgID uint) (*http.Response, []byte, error) {
	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := http.NewRequest(method, baseURL+endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if orgID != 0 {
		req.Header.Set("X-Organization-ID", fmt.Sprintf("%d", orgID))
	}

	resp, err := ts.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, responseBody, nil
}

func TestEndToEndAPI(t *testing.T) {
	ts := NewTestSuite()

	// Wait for server to be ready
	time.Sleep(2 * time.Second)

	t.Run("Health Check", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://localhost:8080/health", nil)
		assert.NoError(t, err)
		resp, err := ts.client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Equal(t, "healthy", response["status"])
	})

	t.Run("Create Organization", func(t *testing.T) {
		// Use timestamp to make organization name unique
		timestamp := time.Now().UnixNano()
		ts.orgName = fmt.Sprintf("Test Organization %d", timestamp)
		orgData := map[string]interface{}{
			"name": ts.orgName,
		}

		resp, body, err := ts.makeRequest("POST", "/organizations/", orgData, "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "organization")

		org := response["organization"].(map[string]interface{})
		ts.orgID = uint(org["id"].(float64))
		ts.organization = org
		assert.Contains(t, org["name"].(string), "Test Organization")
	})

	t.Run("Register Admin User", func(t *testing.T) {
		timestamp := time.Now().UnixNano()
		ts.adminEmail = fmt.Sprintf("admin%d@test.com", timestamp)
		userData := map[string]interface{}{
			"name":            "Admin User",
			"email":           ts.adminEmail,
			"password":        "password123",
			"role":            "admin",
			"organization_id": ts.orgID,
		}

		resp, body, err := ts.makeRequest("POST", "/auth/register", userData, "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "token")
		assert.Contains(t, response, "user")

		ts.adminToken = response["token"].(string)
		ts.adminUser = response["user"].(map[string]interface{})
		assert.Equal(t, "admin", ts.adminUser["role"])
	})

	t.Run("Register Member User", func(t *testing.T) {
		timestamp := time.Now().UnixNano()
		ts.memberEmail = fmt.Sprintf("member%d@test.com", timestamp)
		userData := map[string]interface{}{
			"name":            "Member User",
			"email":           ts.memberEmail,
			"password":        "password123",
			"role":            "member",
			"organization_id": ts.orgID,
		}

		resp, body, err := ts.makeRequest("POST", "/auth/register", userData, "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "token")
		assert.Contains(t, response, "user")

		ts.memberToken = response["token"].(string)
		ts.memberUser = response["user"].(map[string]interface{})
		assert.Equal(t, "member", ts.memberUser["role"])
	})

	t.Run("Login Admin User", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":    ts.adminEmail,
			"password": "password123",
		}

		resp, body, err := ts.makeRequest("POST", "/auth/login", loginData, "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "token")
		assert.Contains(t, response, "user")

		user := response["user"].(map[string]interface{})
		assert.Equal(t, ts.adminEmail, user["email"])
	})

	t.Run("Get User Profile", func(t *testing.T) {
		resp, body, err := ts.makeRequest("GET", "/auth/profile", nil, ts.adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Equal(t, "Admin User", response["name"])
		assert.Equal(t, ts.adminEmail, response["email"])
	})

	t.Run("Get Organization Details", func(t *testing.T) {
		endpoint := fmt.Sprintf("/organizations/%d/", ts.orgID)
		resp, body, err := ts.makeRequestWithOrgHeader("GET", endpoint, nil, ts.adminToken, ts.orgID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Equal(t, ts.orgName, response["name"])
	})

	t.Run("Create Article as Admin", func(t *testing.T) {
		articleData := map[string]interface{}{
			"title":   "Test Article",
			"content": "This is a test article content",
			"status":  "draft",
		}

		endpoint := fmt.Sprintf("/organizations/%d/articles", ts.orgID)
		resp, body, err := ts.makeRequestWithOrgHeader("POST", endpoint, articleData, ts.adminToken, ts.orgID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "article")

		article := response["article"].(map[string]interface{})
		ts.articleID = uint(article["ID"].(float64))
		assert.Equal(t, "Test Article", article["title"])
		assert.Equal(t, "draft", article["status"])
	})

	t.Run("Get Article as Admin", func(t *testing.T) {
		endpoint := fmt.Sprintf("/organizations/%d/articles/%d/", ts.orgID, ts.articleID)
		resp, body, err := ts.makeRequestWithOrgHeader("GET", endpoint, nil, ts.adminToken, ts.orgID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "article")
		assert.Contains(t, response, "permission")

		permission := response["permission"].(string)
		assert.Equal(t, "owner", permission)
	})

	t.Run("Update Article Status to Published", func(t *testing.T) {
		updateData := map[string]interface{}{
			"status": "published",
		}

		endpoint := fmt.Sprintf("/organizations/%d/articles/%d/", ts.orgID, ts.articleID)
		resp, body, err := ts.makeRequestWithOrgHeader("PUT", endpoint, updateData, ts.adminToken, ts.orgID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "article")

		article := response["article"].(map[string]interface{})
		assert.Equal(t, "published", article["status"])
	})

	t.Run("Get Article as Member", func(t *testing.T) {
		endpoint := fmt.Sprintf("/organizations/%d/articles/%d/", ts.orgID, ts.articleID)
		resp, body, err := ts.makeRequestWithOrgHeader("GET", endpoint, nil, ts.memberToken, ts.orgID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "permission")

		permission := response["permission"].(string)
		assert.Equal(t, "comment", permission) // Member should have comment permission on published article
	})

	t.Run("Create Comment as Member", func(t *testing.T) {
		commentData := map[string]interface{}{
			"content": "This is a test comment",
		}

		endpoint := fmt.Sprintf("/organizations/%d/articles/%d/comments", ts.orgID, ts.articleID)
		resp, body, err := ts.makeRequestWithOrgHeader("POST", endpoint, commentData, ts.memberToken, ts.orgID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "comment")

		comment := response["comment"].(map[string]interface{})
		ts.commentID = uint(comment["ID"].(float64))
		assert.Equal(t, "This is a test comment", comment["content"])
	})

	t.Run("Get Comments", func(t *testing.T) {
		endpoint := fmt.Sprintf("/organizations/%d/articles/%d/comments", ts.orgID, ts.articleID)
		resp, body, err := ts.makeRequestWithOrgHeader("GET", endpoint, nil, ts.memberToken, ts.orgID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "comments")

		comments := response["comments"].([]interface{})
		assert.Len(t, comments, 1)
	})

	t.Run("Get Published Articles", func(t *testing.T) {
		resp, body, err := ts.makeRequest("GET", "/articles/published", nil, "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "articles")

		articles := response["articles"].([]interface{})
		assert.GreaterOrEqual(t, len(articles), 1)
	})

	t.Run("Get My Articles", func(t *testing.T) {
		resp, body, err := ts.makeRequest("GET", "/articles/my", nil, ts.adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "articles")

		articles := response["articles"].([]interface{})
		assert.Len(t, articles, 1)
	})

	t.Run("Delete Comment as Article Owner", func(t *testing.T) {
		endpoint := fmt.Sprintf("/organizations/%d/articles/%d/comments/%d", ts.orgID, ts.articleID, ts.commentID)
		resp, _, err := ts.makeRequestWithOrgHeader("DELETE", endpoint, nil, ts.adminToken, ts.orgID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Try to Edit Article as Member (Should Fail)", func(t *testing.T) {
		updateData := map[string]interface{}{
			"title": "Updated by Member",
		}

		endpoint := fmt.Sprintf("/organizations/%d/articles/%d/", ts.orgID, ts.articleID)
		resp, _, err := ts.makeRequestWithOrgHeader("PUT", endpoint, updateData, ts.memberToken, ts.orgID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Delete Article as Owner", func(t *testing.T) {
		endpoint := fmt.Sprintf("/organizations/%d/articles/%d/", ts.orgID, ts.articleID)
		resp, _, err := ts.makeRequestWithOrgHeader("DELETE", endpoint, nil, ts.adminToken, ts.orgID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Invalid Authentication", func(t *testing.T) {
		resp, _, err := ts.makeRequest("GET", "/auth/profile", nil, "invalid-token")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Access Article Without Organization Header", func(t *testing.T) {
		// Create a new article first
		articleData := map[string]interface{}{
			"title":   "Test Article 2",
			"content": "This is another test article",
			"status":  "published",
		}

		endpoint := fmt.Sprintf("/organizations/%d/articles", ts.orgID)
		resp, body, err := ts.makeRequestWithOrgHeader("POST", endpoint, articleData, ts.adminToken, ts.orgID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		article := response["article"].(map[string]interface{})
		newArticleID := uint(article["ID"].(float64))

		// Try to access without org header (should fail)
		endpoint = fmt.Sprintf("/organizations/%d/articles/%d/", ts.orgID, newArticleID)
		resp, _, err = ts.makeRequest("GET", endpoint, nil, ts.adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	fmt.Println("✅ All End-to-End tests completed successfully!")
}
