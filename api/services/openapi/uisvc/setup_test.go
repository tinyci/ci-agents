package uisvc

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	check "github.com/erikh/check"
	grpcHandler "github.com/tinyci/ci-agents/api/handlers/grpc"
	"github.com/tinyci/ci-agents/api/services/grpc/assetsvc"
	"github.com/tinyci/ci-agents/api/services/grpc/datasvc"
	"github.com/tinyci/ci-agents/api/services/grpc/logsvc"
	"github.com/tinyci/ci-agents/api/services/grpc/queuesvc"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/clients/asset"
	"github.com/tinyci/ci-agents/clients/data"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/clients/tinyci"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/testutil/testclients"
	topTypes "github.com/tinyci/ci-agents/types"
)

type uisvcSuite struct {
	datasvcClient  *testclients.DataClient
	queuesvcClient *testclients.QueueClient
	assetsvcClient *asset.Client

	logDoneChan   chan struct{}
	queueDoneChan chan struct{}
	dataDoneChan  chan struct{}
	assetDoneChan chan struct{}
	logHandler    *grpcHandler.H
	dataHandler   *grpcHandler.H
	queueHandler  *grpcHandler.H
	assetHandler  *grpcHandler.H

	logJournal *logsvc.LogJournal
}

var _ = check.Suite(&uisvcSuite{})

func TestUISvc(t *testing.T) {
	check.TestingT(t)
}

func (us *uisvcSuite) SetUpTest(c *check.C) {
	testutil.WipeDB()

	var err error
	us.dataHandler, us.dataDoneChan, err = datasvc.MakeDataServer()
	c.Assert(err, check.IsNil)

	us.queueHandler, us.queueDoneChan, err = queuesvc.MakeQueueServer()
	c.Assert(err, check.IsNil)

	us.assetHandler, us.assetDoneChan, err = assetsvc.MakeAssetServer()
	c.Assert(err, check.IsNil)

	us.logHandler, _, us.logDoneChan, us.logJournal, err = logsvc.MakeLogServer()
	c.Assert(err, check.IsNil)

	go us.logJournal.Tail()

	us.datasvcClient, err = testclients.NewDataClient()
	c.Assert(err, check.IsNil)

	us.queuesvcClient, err = testclients.NewQueueClient(us.datasvcClient)
	c.Assert(err, check.IsNil)

	us.assetsvcClient, err = asset.NewClient(config.DefaultServices.Asset.String(), nil, false)
	c.Assert(err, check.IsNil)
}

func (us *uisvcSuite) TearDownTest(c *check.C) {
	close(us.dataDoneChan)
	close(us.queueDoneChan)
	close(us.logDoneChan)
	close(us.assetDoneChan)
	time.Sleep(100 * time.Millisecond)
}

// MakeUIServer makes a uisvc.
func MakeUIServer(client github.Client) (*H, chan struct{}, *tinyci.Client, *tinyci.Client, error) {
	conf := config.UserConfig{
		OAuth: config.OAuthConfig{
			ClientID:     "client id",
			ClientSecret: "client secret",
			RedirectURL:  "http://localhost:6010/login",
		},
		ClientConfig: config.TestClientConfig,
		URL:          "http://localhost",
		Auth: config.AuthConfig{
			SessionCryptKey: "0431d583a48a00243cc3d3d596ed362d77c50be4848dbf0d2f52bab841f072f9",
			TokenCryptKey:   "1431d583a48a00243cc3d3d596ed362d77c50be4848dbf0d2f52bab841f072f9",
		},
		Websockets: config.Websockets{
			InsecureWebSockets: true,
		},
		RequestLogging: true,
		Port:           6010,
	}
	if err := conf.Auth.Validate(true); err != nil {
		return nil, nil, nil, nil, err
	}

	config.SetDefaultGithubClient(client, "")
	finished := make(chan struct{})

	handler := &H{
		ServiceName: "uisvc",
		Port:        6010,
		Config:      conf,
	}

	var err error
	doneChan, err := handler.Boot(finished)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	d, err := data.New(config.DefaultServices.Data.String(), nil, false)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	content, err := json.Marshal(&topTypes.OAuthToken{Username: "erikh", Token: "dummy", Scopes: []string{"repo"}})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	u, err := d.PutUser(context.Background(), &types.User{Username: "erikh", TokenJSON: content})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	for _, cap := range topTypes.AllCapabilities {
		if err := d.AddCapability(context.Background(), u, cap); err != nil {
			return nil, nil, nil, nil, err
		}
	}

	token, err := d.GetToken(context.Background(), "erikh")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tc, err := tinyci.New(config.DefaultServices.UI.String(), token, nil)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if _, err := tc.Errors(context.Background()); err != nil { // connectivity check
		return nil, nil, nil, nil, err
	}

	content, err = json.Marshal(&topTypes.OAuthToken{Token: "dummy"})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	_, err = d.PutUser(context.Background(), &types.User{Username: "erikh2", TokenJSON: content})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	token, err = d.GetToken(context.Background(), "erikh2")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	utc, err := tinyci.New(config.DefaultServices.UI.String(), token, nil)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return handler, doneChan, tc, utc, nil
}
