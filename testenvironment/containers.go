package testenvironment

import (
	"sync"
	"github.com/byrnedo/prefab"
	"time"
)

type startContainerFunc func(...prefab.ConfOverrideFunc) (string, string)
type waitFunc func(string, time.Duration) (error)

type startAndWaitFuncs struct {
	name string
	startFunc startContainerFunc
	waitFunc waitFunc
}

type urlTuple struct {
	name string
	url string
}
// Uses prefab (github.com/byrnedo/prefab) to create docker containers for test environments
type TestEnvironment struct {
	maxTimeout time.Duration
	ids []string
	funcList []startAndWaitFuncs
	idChan chan string
	urlChan chan urlTuple
	urls map[string]string
	doneChan chan bool
}

func NewTestEnvironment() *TestEnvironment {
	return &TestEnvironment{}
}

func (this *TestEnvironment) GetUrl(name string) (string, bool) {
	val, found := this.urls[name]
	return val, found
}

// Set the max timeout for waiting on individual containers
func (this *TestEnvironment) WithMaxTimeout(timeout time.Duration) *TestEnvironment {
	this.maxTimeout = timeout
	return this
}

// Queue a mongo container
func (this *TestEnvironment) WithMongo() *TestEnvironment {
	this.funcList = append(this.funcList, startAndWaitFuncs{
		name: "mongo",
		startFunc: prefab.StartMongoContainer,
		waitFunc: prefab.WaitForMongo,
	})
	return this
}

// Queue a mysql container
func (this *TestEnvironment) WithMysql() *TestEnvironment {
	this.funcList = append(this.funcList, startAndWaitFuncs{
		name: "mysql",
		startFunc: prefab.StartMysqlContainer,
		waitFunc: prefab.WaitForMysql,
	})
	return this
}

// Queue a mongo tmpfs container
func (this *TestEnvironment) WithMongoTmpfs() *TestEnvironment {
	this.funcList = append(this.funcList, startAndWaitFuncs{
		name: "mongo",
		startFunc: prefab.StartMongoTmpfsContainer,
		waitFunc: prefab.WaitForMongo,
	})
	return this
}

// Queue a mysql container on tmpfs
func (this *TestEnvironment) WithMysqlTmpfs() *TestEnvironment {
	this.funcList = append(this.funcList, startAndWaitFuncs{
		name: "mysql",
		startFunc: prefab.StartMysqlTmpfsContainer,
		waitFunc: prefab.WaitForMysql,
	})
	return this
}

// Queue a nats container
func (this *TestEnvironment) WithNats() *TestEnvironment {
	this.funcList = append(this.funcList, startAndWaitFuncs{
		name: "nats",
		startFunc: prefab.StartNatsContainer,
		waitFunc: prefab.WaitForNats,
	})
	return this
}

// Queue a postgres container
func (this *TestEnvironment) WithPostgres() *TestEnvironment {
	this.funcList = append(this.funcList, startAndWaitFuncs{
		name: "postgres",
		startFunc: prefab.StartPostgresContainer,
		waitFunc: prefab.WaitForPostgres,
	})
	return this
}

// Launch the queued containers
func (this *TestEnvironment) Launch() {
	wg :=  sync.WaitGroup{}
	wg.Add(len(this.funcList))

	if this.maxTimeout == 0 {
		this.maxTimeout = 25 * time.Second
	}
	this.urls = make(map[string]string)

	this.doneChan = make(chan bool)
	this.idChan = make(chan string)
	this.urlChan = make(chan urlTuple)

	go func() {
		for {
			select {
			case id := <-this.idChan:
				if len(id) > 0 {
					this.ids = append(this.ids, id)
				}
			case  <-this.doneChan:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case tuple := <-this.urlChan:
				this.urls[tuple.name] = tuple.url
			case <-this.doneChan:
				return
			}
		}
	}()


	for _, group := range this.funcList {
		go func(grpCpy startAndWaitFuncs){
			id, url := grpCpy.startFunc()
			this.idChan <- id
			grpCpy.waitFunc(url, this.maxTimeout)
			this.urlChan <- urlTuple{grpCpy.name, url}
			wg.Done()
		}(group)
	}

	wg.Wait()
	this.doneChan <- true
	this.doneChan <- true
	close(this.idChan)
	close(this.urlChan)
	return
}

func (this *TestEnvironment) Cleanup() {
	wg := sync.WaitGroup{}

	wg.Add(len(this.ids))
	for _, id := range this.ids {
		go func(){
			prefab.Remove(id)
			wg.Done()
		}()
		time.Sleep(50 * time.Millisecond)
	}
	wg.Wait()
}

