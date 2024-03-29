package httppub

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/tiennam886/manager/pkg/messaging"
)

var _eventPubs map[string][]Publisher
var _mu sync.RWMutex

type Publisher struct {
	Name     string
	Endpoint url.URL
	Header   http.Header
}

func NewPublisher(name string, endpoint url.URL, header http.Header) *Publisher {
	return &Publisher{name, endpoint, header}
}

func (p *Publisher) DoRequest(e messaging.Event) error {
	if e.Name() == "" || e.JSON() == nil || len(e.JSON()) == 0 {
		return errors.New("publisher cannot do request with invalid event")
	}

	req, err := http.NewRequest(http.MethodPost, p.Endpoint.String(), bytes.NewBuffer(e.JSON()))
	if err != nil {
		return err
	}

	req.Header = p.Header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp == nil {
		return fmt.Errorf("published event [%s] received nil response", e.Name())
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("published event [%s] received invalid status code (%v)", e.Name(), resp.StatusCode)
	}

	return nil
}

func ConnectPub(p Publisher, eventName string) error {
	_mu.Lock()
	defer _mu.Unlock()

	pubs, ok := _eventPubs[eventName]
	if !ok {
		pubs = make([]Publisher, 0)
	}
	fmt.Println(_eventPubs)
	for _, pub := range pubs {
		if pub.Name == p.Name {
			return errors.New("publisher already registered")
		}
	}

	pubs = append(pubs, p)
	_eventPubs[eventName] = pubs

	return nil
}

func DisconnectPub(p Publisher, eventName string) error {
	_mu.Lock()
	defer _mu.Unlock()

	pubs, ok := _eventPubs[eventName]
	if !ok {
		pubs = make([]Publisher, 0)
	}

	tmp := make([]Publisher, 0)
	for _, pub := range pubs {
		if pub.Name == p.Name {
			continue
		}
		tmp = append(tmp, pub)
	}

	_eventPubs[eventName] = tmp

	return nil
}

// Publish @TODO: handling error
func Publish(event messaging.Event) {
	_mu.RLock()
	pubs := _eventPubs[event.Name()]
	_mu.RUnlock()

	for _, pub := range pubs {
		go func() {
			_ = pub.DoRequest(event)
		}()
	}
}

func init() {
	_eventPubs = make(map[string][]Publisher)
}
