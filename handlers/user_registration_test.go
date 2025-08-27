package handlers_test

/*
const testFromAddress = "testing@afterwill.life"

func setupEmailServerAndClient(t *testing.T) (*smtpmock.Server, *email.EmailService) {
	t.Helper()
	require := require.New(t)
	emailServer, address := mocksmtp.SetupEmailMock()
	require.NoError(emailServer.Start())
	client, err := mail.NewClient(address)
	require.NoError(err)
	emailService := email.NewEmailService(client, testFromAddress)
	return emailServer, emailService
}

func TestRegisterUser(t *testing.T) {
	const emailVerificationRoute = "https://afterwill.life/email/verify"
	t.Run("valid input", func(t *testing.T) {
		emailServer, emailService := setupEmailServerAndClient(t)

		c, w, assert := setupHandlerTest(t)
		mock := newMockDB()
		username := "mark"
		name := "Down"
		email := "test@google.com"
		password := "5UP3RS3CR37"
		input, err := sonic.Marshal(handlers.RegistrationInput{
			Username: username,
			Name:     &name,
			Email:    email,
			Password: password,
		})
		if err != nil {
			t.Fatalf("sonic failed to bind json, %v", err)
		}
		setGinHttpBody(c, input)

		handlers.RegisterUser(mock, createHash, emailService, emailVerificationRoute, JWT_SECRET)(c)

		assertHTTPStatus(t, c, http.StatusCreated, w, "http status code should indicate that the user was created")
		field := mock.Users[0]

		hash, _ := createHash(password, argon2id.DefaultParams)
		assert.Equal(hash, field.PasswordHash)
		assert.Equal(username, field.Username)
		assert.Equal(name, *field.Profile.Name)
		assert.Equal(email, field.Email)
	})
	t.Run("user already exists", func(t *testing.T) {
		c, w, assert := setupHandlerTest(t)
		mock := newMockDB()
		alreadyExistingUserName := "asdfasdf"
		alreadyExistinguser := models.User{
			ID: 1, Profile: &models.Profile{Name: &alreadyExistingUserName}, Username: "testuseralreadyexists", Email: "gregor@gregtech.eu",
		}
		mock.Users = append(mock.Users, alreadyExistinguser)

		input, err := sonic.Marshal(handlers.RegistrationInput{
			Username: alreadyExistinguser.Username, Email: alreadyExistinguser.Email, Password: "vverysecurepassword", Name: alreadyExistinguser.Profile.Name,
		})
		if err != nil {
			t.Fatalf("sonic failed to bind json: %v", err)
		}
		setGinHttpBody(c, input)

		handlers.RegisterUser(mock, createHash)(c)

		assertHTTPStatus(t, c, http.StatusConflict, w, "http status code should indicate that a user already exists")
		assert.Equal([]models.User{alreadyExistinguser}, mock.Users, "there should be just one user created")
	})
}
*/
