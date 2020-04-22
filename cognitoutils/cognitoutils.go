package cognitoutils

import (
	"fmt"

	"github.com/AmeDigital/aws-utils-go/sessionutils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func CreateUser(username string, userAttributes map[string]string, userPoolId string) error {

	fmt.Printf("CognitoCreateUser: Creating in userPool %s the username %s with attributes %+v\n", userPoolId, username, userAttributes)

	var attributes []*cognitoidentityprovider.AttributeType

	for k := range userAttributes {
		attributeType := &cognitoidentityprovider.AttributeType{
			Name:  aws.String(k),
			Value: aws.String(userAttributes[k]),
		}

		attributes = append(attributes, attributeType)
	}

	createUserInput := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:     aws.String(userPoolId),
		Username:       aws.String(username),
		UserAttributes: attributes,
	}

	var cognitoIdpClient *cognitoidentityprovider.CognitoIdentityProvider = cognitoidentityprovider.New(sessionutils.Session)
	_, err := cognitoIdpClient.AdminCreateUser(createUserInput)
	if err != nil {
		return err
	}

	fmt.Printf("CognitoCreateUser: User %s created successfully.\n", username)

	return nil
}

func GetUserIdentityId(username string, password string, appClientId string, identityPoolId string, cognitoTokenProvider string) (identityId string, err error) {
	fmt.Println("CognitoGetUserIdentityId start")
	var cognitoIdpClient *cognitoidentityprovider.CognitoIdentityProvider = cognitoidentityprovider.New(sessionutils.Session)
	cognitoIdentityClient := cognitoidentity.New(sessionutils.Session)

	initiateAuthInput := &cognitoidentityprovider.InitiateAuthInput{
		ClientId: &appClientId,
		AuthFlow: aws.String("CUSTOM_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": &username,
		},
	}

	initiateAuthOutput, err := cognitoIdpClient.InitiateAuth(initiateAuthInput)
	if err != nil {
		return "", err
	}

	respondToAuthChallengeInput := &cognitoidentityprovider.RespondToAuthChallengeInput{
		ClientId:      &appClientId,
		ChallengeName: aws.String("CUSTOM_CHALLENGE"),
		Session:       initiateAuthOutput.Session,
		ChallengeResponses: map[string]*string{
			"USERNAME": &username,
			"ANSWER":   &password,
		},
	}

	respondToAuthChallengeOutput, err := cognitoIdpClient.RespondToAuthChallenge(respondToAuthChallengeInput)
	if err != nil {
		return "", err
	}

	idToken := respondToAuthChallengeOutput.AuthenticationResult.IdToken

	getIdInput := &cognitoidentity.GetIdInput{
		IdentityPoolId: &identityPoolId,
		Logins: map[string]*string{
			cognitoTokenProvider: idToken,
		},
	}

	getIdOutput, err := cognitoIdentityClient.GetId(getIdInput)
	if err != nil {
		return "", err
	}

	fmt.Println("CognitoGetUserIdentityId end")

	return *getIdOutput.IdentityId, nil
}

func ListUsers(userPoolId string) (users []*cognitoidentityprovider.UserType, err error) {
	var cognitoIdpClient *cognitoidentityprovider.CognitoIdentityProvider = cognitoidentityprovider.New(sessionutils.Session)

	listUsersInput := &cognitoidentityprovider.ListUsersInput{
		UserPoolId: &userPoolId,
	}

	listUsersOutput, err := cognitoIdpClient.ListUsers(listUsersInput)

	return listUsersOutput.Users, err
}

func ListUserNames(userPoolId string) (usernames []*string, err error) {
	users, err := ListUsers(userPoolId)

	if err != nil {
		return usernames, err
	}

	for _, user := range users {
		usernames = append(usernames, user.Username)
	}

	return usernames, err
}

func ListUsersWithPrefixFilter(attributeName string, prefix string, userPoolId string) (usernames []string, err error) {
	var cognitoIdpClient *cognitoidentityprovider.CognitoIdentityProvider = cognitoidentityprovider.New(sessionutils.Session)

	listUsersInput := &cognitoidentityprovider.ListUsersInput{
		UserPoolId: &userPoolId,
		Filter:     aws.String(attributeName + " ^= " + prefix),
	}

	listUsersOutput, err := cognitoIdpClient.ListUsers(listUsersInput)

	if err != nil {
		return usernames, err
	}

	for _, user := range listUsersOutput.Users {
		for _, attribute := range user.Attributes {
			if *attribute.Name == "username" {
				usernames = append(usernames, *attribute.Value)
			}
		}
	}

	return usernames, nil
}
