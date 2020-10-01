package main


import "testing"

func TestSignup(t *testing.T) {
	creds := &Credentials{
		Password: "testuser",
		Username: "testpassword",
	}
	err := Signup(creds)
	if err != nil {
		t.Errorf("Sign up is incorrect")
	}
}


func TestSignin(t *testing.T) {
	creds := &Credentials{
		Password: "testuser",
		Username: "testpassword",
	}
	_, err := Signin(creds)
	if err != nil {
		if err.Error() == "Password is not valid" {
			t.Errorf("Unauthorized")
		} else {
			t.Errorf("Sign in is incorrect")
		}
	}
}

func TestWelcomeS(t *testing.T) {
	sessionToken := "some_session_token"
	_, err := WelcomeS(sessionToken)
	if err != nil {
		if err.Error() == "No sessions found for user" {
			t.Errorf("Unauthorized")
		} else {
			t.Errorf("Incorrect session token was provided")
		}
	}
}

func TestRefreshS(t *testing.T) {
	sessionToken := "some_session_token"
	_, err := RefreshS(sessionToken)
	if err != nil {
		if err.Error() == "No sessions found for user" {
			t.Errorf("Unauthorized")
		} else {
			t.Errorf("Incorrect session token was provided")
		}
	}
}