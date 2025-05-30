package main

import (
	"net/http"
	"net/url"
	"snippetbox.samedarslan28.net/internal/assert"
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, body := ts.get(t, "/ping")
	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestHome(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, _ := ts.get(t, "/")
	assert.Equal(t, statusCode, http.StatusOK)
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, tt.wantCode, code)

			if tt.wantBody != "" {
				assert.Contains(t, body, tt.wantBody)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")

	csrfToken := extractCSRFToken(t, body)

	t.Logf("CSRF token is: %q", csrfToken)
}

func TestUserSignupPost(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid email",
			userName:     validName,
			userEmail:    "bob@example.",
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Short password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pa$$",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Duplicate email",
			userName:     validName,
			userEmail:    "dupe@example.com",
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, `/user/signup`, form)

			assert.Equal(t, tt.wantCode, code)

			if tt.wantFormTag != "" {
				assert.Contains(t, body, tt.wantFormTag)
			}
		})
	}
}

func TestUserLogin(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/login")
	validCSRFToken := extractCSRFToken(t, body)

	type testCase struct {
		name         string
		email        string
		password     string
		csrfToken    string
		expectedCode int
	}

	tests := []testCase{
		{
			name:         "Valid credentials",
			email:        "alice@example.com",
			password:     "pa$$word",
			csrfToken:    validCSRFToken,
			expectedCode: http.StatusSeeOther,
		},
		{
			name:         "Blank email",
			email:        "",
			password:     "pa$$word",
			csrfToken:    validCSRFToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid email format",
			email:        "bademail",
			password:     "pa$$word",
			csrfToken:    validCSRFToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Blank password",
			email:        "alice@example.com",
			password:     "",
			csrfToken:    validCSRFToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Short password",
			email:        "alice@example.com",
			password:     "short",
			csrfToken:    validCSRFToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid credentials",
			email:        "alice@example.com",
			password:     "wrongpass",
			csrfToken:    validCSRFToken,
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "Invalid CSRF token",
			email:        "alice@example.com",
			password:     "pa$$word",
			csrfToken:    "invalid-token",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid User",
			email:        "test@gmail.com",
			password:     "test_password",
			csrfToken:    validCSRFToken,
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			form.Add("csrf_token", tt.csrfToken)

			statusCode, _, _ := ts.postForm(t, "/user/login", form)
			assert.Equal(t, tt.expectedCode, statusCode)
		})
	}
}

func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/snippet/create")
		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, headers.Get("Location"),
			"/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		authenticateTestUser(t, ts)

		code, _, body := ts.get(t, "/snippet/create")
		assert.Equal(t, code, http.StatusOK)
		assert.Contains(t, body, "<form action='/snippet/create' method= 'POST'>")
	})
}

func TestAccount(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	authenticateTestUser(t, ts)

	statusCode, _, _ := ts.get(t, "/user/account")
	assert.Equal(t, statusCode, http.StatusOK)
}

func TestLogoutPost(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("email", "alice@example.com")
	form.Add("password", "pa$$word")
	form.Add("csrf_token", csrfToken)
	_, _, _ = ts.postForm(t, "/user/login", form)

	logoutForm := url.Values{}

	logoutForm.Add("csrf_token", csrfToken)

	statusCode, _, _ := ts.postForm(t, "/user/logout", logoutForm)
	assert.Equal(t, statusCode, http.StatusSeeOther)

}

func TestAbout(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, _ := ts.get(t, "/about")
	assert.Equal(t, statusCode, http.StatusOK)
}

func TestChangePassword(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	authenticateTestUser(t, ts)

	statusCode, _, _ := ts.get(t, "/user/account/password/update")
	assert.Equal(t, statusCode, http.StatusOK)
}

func TestChangePasswordPost(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	csrfToken := authenticateTestUser(t, ts)
	tests := []struct {
		name            string
		currentPassword string
		newPassword     string
		confirmPassword string
		csrfToken       string
		wantCode        int
	}{
		{
			name:            "Empty form submission",
			currentPassword: "",
			newPassword:     "",
			confirmPassword: "",
			csrfToken:       csrfToken,
			wantCode:        http.StatusUnprocessableEntity,
		},
		{
			name:            "Passwords do not match",
			currentPassword: "pa$$word",
			newPassword:     "newPassword1",
			confirmPassword: "newPassword2",
			csrfToken:       csrfToken,
			wantCode:        http.StatusUnprocessableEntity,
		},
		{
			name:            "Invalid current password",
			currentPassword: "wrongpassword",
			newPassword:     "newPassword1",
			confirmPassword: "newPassword1",
			csrfToken:       csrfToken,
			wantCode:        http.StatusUnprocessableEntity,
		},
		{
			name:            "Valid password change",
			currentPassword: "pa$$word",
			newPassword:     "newPassword1",
			confirmPassword: "newPassword1",
			csrfToken:       csrfToken,
			wantCode:        http.StatusSeeOther,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("current_password", tt.currentPassword)
			form.Add("new_password", tt.newPassword)
			form.Add("new_password_confirm", tt.confirmPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, _ := ts.postForm(t, "/user/account/password/update", form)

			if code != tt.wantCode {
				t.Errorf("%s: expected status code %d; got %d", tt.name, tt.wantCode, code)
			}
		})
	}
}
func TestSnippetCreatePost(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	csrfToken := authenticateTestUser(t, ts) // Logs in the test user and returns CSRF token

	tests := []struct {
		name      string
		title     string
		content   string
		expires   string
		csrfToken string
		wantCode  int
	}{
		{
			name:      "Valid submission",
			title:     "Test Snippet",
			content:   "This is a test snippet content.",
			expires:   "7",
			csrfToken: csrfToken,
			wantCode:  http.StatusSeeOther,
		},
		{
			name:      "Empty title and content",
			title:     "",
			content:   "",
			expires:   "7",
			csrfToken: csrfToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Invalid expires value",
			title:     "Test Snippet",
			content:   "Some content here",
			expires:   "999",
			csrfToken: csrfToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Missing CSRF token",
			title:     "Test Snippet",
			content:   "Some content here",
			expires:   "7",
			csrfToken: "", // no CSRF
			wantCode:  http.StatusBadRequest,
		},
		{
			name:      "Title exceeds character limit",
			title:     strings.Repeat("a", 101),
			content:   "Some content here",
			expires:   "1",
			csrfToken: csrfToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("title", tt.title)
			form.Add("content", tt.content)
			form.Add("expires", tt.expires)
			form.Add("csrf_token", tt.csrfToken)

			code, _, _ := ts.postForm(t, "/snippet/create", form)

			if code != tt.wantCode {
				t.Errorf("%s: expected status code %d, got %d", tt.name, tt.wantCode, code)
			}
		})
	}
}
