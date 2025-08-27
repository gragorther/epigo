package mocksmtp

import smtpmock "github.com/mocktools/go-smtp-mock/v2"

func SetupEmailMock() (server *smtpmock.Server, address string) {
	address = "127.0.0.1"
	return smtpmock.New(smtpmock.ConfigurationAttr{HostAddress: address}), address
}
