package helper

import "testing"

func TestComparePassword(t *testing.T) {
	type TestStruct struct {
		Password     string
		HashPassword string
		Expected     bool
	}

	// Valid password
	test := TestStruct{
		Password:     "password",
		HashPassword: "",
		Expected:     true,
	}

	helper := NewHelper(NewHelperOptions{
		JwtPrivateKeyPath: "../storage/key.pem",
		JwtPublicKeyPath:  "../storage/key.pem.pub",
	})

	hashPassword, err := helper.HashPassword(test.Password)
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}
	test.HashPassword = hashPassword
	if err := helper.ComparePassword(test.Password, test.HashPassword); err != nil {
		t.Errorf("Expected true, got false")
	}

	// Invalid password
	hashPassword, err = helper.HashPassword("wrongpassword")
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}
	test.HashPassword = hashPassword
	if err := helper.ComparePassword(test.Password, test.HashPassword); err == nil {
		t.Errorf("Expected false, got true")
	}
}

func TestGetToken(t *testing.T) {
	helper := NewHelper(NewHelperOptions{
		JwtPrivateKeyPath: "../storage/key.pem",
		JwtPublicKeyPath:  "../storage/key.pem.pub",
	})

	tests := []struct {
		caseName string
		Input    string
		Expected string
	}{
		{
			caseName: "Empty token",
			Input:    "test",
			Expected: "",
		},
		{
			caseName: "Valid token",
			Input:    "Bearer test",
			Expected: "test",
		},
		{
			caseName: "Valid token with space",
			Input:    "Bearer test test",
			Expected: "test test",
		},
		{
			caseName: "Invalid token",
			Input:    "Bearer102938jaskld",
			Expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			token := helper.GetToken(test.Input)
			if token != test.Expected {
				t.Errorf("Expected %s, got %s", test.Expected, token)
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	helper := NewHelper(NewHelperOptions{
		JwtPrivateKeyPath: "../storage/key.pem",
		JwtPublicKeyPath:  "../storage/key.pem.pub",
	})

	token := ""
	helper.GenerateAccessToken(&token, 1)
	if _, err := helper.VerifyToken(token); err != nil {
		t.Errorf("Expected error, got nil")
	}

	if _, err := helper.VerifyToken("invalid"); err == nil {
		t.Errorf("Expected error, got nil")
	}
}
