package testservers

import (
	"sync"

	transport "github.com/erikh/go-transport"
	"github.com/sirupsen/logrus"
	assetsvc "github.com/tinyci/ci-agents/api/assetsvc/processors"
	datasvc "github.com/tinyci/ci-agents/api/datasvc/processors"
	logsvc "github.com/tinyci/ci-agents/api/logsvc/processors"
	queuesvc "github.com/tinyci/ci-agents/api/queuesvc/processors"
	d "github.com/tinyci/ci-agents/clients/data"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/clients/tinyci"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/gen/svc/uisvc/restapi"
	"github.com/tinyci/ci-agents/grpc/handler"
	"github.com/tinyci/ci-agents/grpc/services/asset"
	"github.com/tinyci/ci-agents/grpc/services/data"
	"github.com/tinyci/ci-agents/grpc/services/log"
	"github.com/tinyci/ci-agents/grpc/services/queue"
	"github.com/tinyci/ci-agents/handlers"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
	"google.golang.org/grpc"
)

var clients = config.ClientConfig{
	Data:  "localhost:6000",
	Queue: "localhost:6001",
	UI:    "http://localhost:6010",
	Asset: "localhost:6002",
	Log:   "localhost:6005",
}

// MakeUIServer makes a uisvc.
func MakeUIServer(client github.Client) (*handlers.H, chan struct{}, *tinyci.Client, *tinyci.Client, *errors.Error) {
	h := &handlers.H{
		Config: restapi.HandlerConfig{},
		Service: config.Service{
			Name: "uisvc",
		},
		UserConfig: config.UserConfig{
			OAuth: config.OAuthConfig{
				ClientID:     "client id",
				ClientSecret: "client secret",
				RedirectURL:  "http://localhost:6010/login",
			},
			ClientConfig: clients,
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
		},
	}

	d, err := d.New("localhost:6000", nil)
	if err != nil {
		return nil, nil, nil, nil, errors.New(err)
	}

	config.DefaultGithubClient = client
	doneChan, err := handlers.Boot(nil, h)
	if err != nil {
		return nil, nil, nil, nil, errors.New(err)
	}

	u, err := d.PutUser(&model.User{Username: "erikh", Token: &types.OAuthToken{Token: "dummy", Scopes: []string{"repo"}}})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	for _, cap := range model.AllCapabilities {
		if err := d.AddCapability(u, cap); err != nil {
			return nil, nil, nil, nil, err
		}
	}

	token, err := d.GetToken("erikh")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tc, err := tinyci.New("http://localhost:6010", token)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if _, err := tc.Errors(); err != nil { // connectivity check
		return nil, nil, nil, nil, err
	}

	_, err = d.PutUser(&model.User{Username: "erikh2", Token: &types.OAuthToken{Token: "dummy"}})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	token, err = d.GetToken("erikh2")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	utc, err := tinyci.New("http://localhost:6010", token)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return h, doneChan, tc, utc, nil
}

// MakeDataServer makes an instance of the datasvc on port 6000. It returns a
// chan which can be closed to terminate it, and any boot-time errors.
func MakeDataServer() (*handler.H, chan struct{}, *errors.Error) {
	h := &handler.H{
		Service: config.Service{
			UseDB: true,
			Name:  "datasvc",
		},
		UserConfig: config.UserConfig{
			ClientConfig: clients,
			DSN:          testutil.TestDBConfig,
			Port:         6000,
			URL:          "url",
			Auth: config.AuthConfig{
				TokenCryptKey: "1431d583a48a00243cc3d3d596ed362d77c50be4848dbf0d2f52bab841f072f9",
			},
		},
	}

	t, err := transport.Listen(nil, "tcp", "localhost:6000")
	if err != nil {
		return nil, nil, errors.New(err)
	}

	srv := grpc.NewServer()
	data.RegisterDataServer(srv, &datasvc.DataServer{H: h})

	doneChan, err := h.Boot(t, srv)
	return h, doneChan, errors.New(err)
}

// MakeAssetServer makes an instance of the assetsvc on port 6000. It returns a
// chan which can be closed to terminate it, and any boot-time errors.
func MakeAssetServer() (*handler.H, chan struct{}, *errors.Error) {
	t, err := transport.Listen(nil, "tcp", "localhost:6002")
	if err != nil {
		return nil, nil, errors.New(err)
	}

	h := &handler.H{
		Service: config.Service{Name: "assetsvc"},
		UserConfig: config.UserConfig{
			Auth: config.AuthConfig{
				TokenCryptKey: "1431d583a48a00243cc3d3d596ed362d77c50be4848dbf0d2f52bab841f072f9",
			},
		},
	}

	srv := grpc.NewServer()
	asset.RegisterAssetServer(srv, &assetsvc.AssetServer{})

	doneChan, err := h.Boot(t, srv)
	return h, doneChan, errors.New(err)
}

// MakeQueueServer makes an instance of the queuesvc on port 6001. It returns a
// chan which can be closed to terminate it, and any boot-time errors.
func MakeQueueServer() (*handler.H, chan struct{}, *errors.Error) {
	h := &handler.H{
		Service: config.Service{
			Name: "queuesvc",
		},
		UserConfig: config.UserConfig{
			DSN:          testutil.TestDBConfig,
			ClientConfig: clients,
			URL:          "url",
			Port:         6001,
			Auth: config.AuthConfig{
				TokenCryptKey: "1431d583a48a00243cc3d3d596ed362d77c50be4848dbf0d2f52bab841f072f9",
			},
		},
	}

	t, err := transport.Listen(nil, "tcp", "localhost:6001")
	if err != nil {
		return nil, nil, errors.New(err)
	}

	srv := grpc.NewServer()
	queue.RegisterQueueServer(srv, &queuesvc.QueueServer{H: h})

	doneChan, err := h.Boot(t, srv)
	return h, doneChan, errors.New(err)
}

// LogJournal is a journal of log entries intended to facilitate mocking the logsvc.
type LogJournal struct {
	Journal map[string][]*log.LogMessage
	mutex   sync.Mutex
}

// Tail echoes all journal entries to stdout
func (lj *LogJournal) Tail() {
	for {
		lj.mutex.Lock()
		for _, items := range lj.Journal {
			for _, item := range items {
				res := logrus.Fields{}
				for key, val := range item.Fields.Fields {
					res[key] = val.GetStringValue()
				}

				logrus.WithFields(res).Println(item.Message)
			}
		}
		// XXX manual version of Reset() to avoid deadlocking
		lj.Journal = map[string][]*log.LogMessage{}
		lj.mutex.Unlock()
	}
}

// Reset resets the log journal, erasing all recorded messages.
func (lj *LogJournal) Reset() {
	lj.mutex.Lock()
	defer lj.mutex.Unlock()
	lj.Journal = map[string][]*log.LogMessage{}
}

// Append appends a message.
func (lj *LogJournal) Append(level string, msg *log.LogMessage) {
	lj.mutex.Lock()
	defer lj.mutex.Unlock()

	if _, ok := lj.Journal[level]; !ok {
		lj.Journal[level] = []*log.LogMessage{}
	}

	lj.Journal[level] = append(lj.Journal[level], msg)
}

// MakeLogServer makes a logsvc.
func MakeLogServer() (*handler.H, chan struct{}, *LogJournal, *errors.Error) {
	journal := &LogJournal{Journal: map[string][]*log.LogMessage{}}

	logDispatch := logsvc.DispatchTable{
		logsvc.LevelDebug: func(wf logsvc.Dispatcher, msg *log.LogMessage) {
			journal.Append(logsvc.LevelDebug, msg)
		},
		logsvc.LevelError: func(wf logsvc.Dispatcher, msg *log.LogMessage) {
			journal.Append(logsvc.LevelError, msg)
		},
		logsvc.LevelInfo: func(wf logsvc.Dispatcher, msg *log.LogMessage) {
			journal.Append(logsvc.LevelInfo, msg)
		},
	}

	h := &handler.H{
		Service: config.Service{
			Name: "logsvc",
		},
		UserConfig: config.UserConfig{
			Port: 6005,
			// FIXME this is really dumb and should be unnecessary
			Auth: config.AuthConfig{
				TokenCryptKey: "1431d583a48a00243cc3d3d596ed362d77c50be4848dbf0d2f52bab841f072f9",
			},
		},
	}

	t, err := transport.Listen(nil, "tcp", "localhost:6005")
	if err != nil {
		return nil, nil, nil, errors.New(err)
	}

	srv := grpc.NewServer()
	log.RegisterLogServer(srv, logsvc.New(logDispatch))

	doneChan, err := h.Boot(t, srv)
	return h, doneChan, journal, errors.New(err)
}
