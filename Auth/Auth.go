package Auth

import (
	"context"
	Service "../LINE"
	"golang.org/x/xerrors"
)

type LINEClient struct {
	Profile           *Service.Profile
	
	TalkServiceClient *ThriftClient
	PollServiceClient *ThriftClient

	Talk              *Service.TalkServiceClient
	Poll              *Service.TalkServiceClient
}

type Config struct {
	Host                      string 
	TalkService               string 
	PollService               string 
	UserAgent                 string 
	LINEApp                   string 
	AccessToken               string 
}

func NewLINEClient(config Config) (*LINEClient, error) {
	client := new(LINEClient)
	{
		url := config.Host + config.TalkService
		tClient, err := NewThriftClient(url)
		if err != nil {
			return nil, xerrors.Errorf("failed to generate thrift client: %w", err)
		}
		serviceClient := Service.NewTalkServiceClient(tClient.StandardClient)
		client.TalkServiceClient = tClient
		client.Talk = serviceClient
	}
	{
		url := config.Host + config.PollService
		tClient, err := NewThriftClient(url)
		if err != nil {
			return nil, xerrors.Errorf("failed to generate thrift client: %w", err)
		}
		serviceClient := Service.NewTalkServiceClient(tClient.StandardClient)
		client.PollServiceClient = tClient
		client.Poll = serviceClient
	}
	client.TalkServiceClient.SetHeader("User-Agent", config.UserAgent)
	client.TalkServiceClient.SetHeader("X-Line-Application", config.LINEApp)
	client.TalkServiceClient.SetHeader("X-Line-Access", config.AccessToken)

	client.PollServiceClient.SetHeader("User-Agent", config.UserAgent)
	client.PollServiceClient.SetHeader("X-Line-Application", config.LINEApp)
	client.PollServiceClient.SetHeader("X-Line-Access", config.AccessToken)

	profile, err := client.Talk.GetProfile(context.Background(), Service.SyncReason_UNKNOWN)
	if err != nil {
		return nil, err
	}
	client.Profile = profile
	return client, nil
}