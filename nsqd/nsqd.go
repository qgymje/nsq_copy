package nsqd

import (
	"crypto/tls"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qgymje/nsq_copy/internal/clusterinfo"
	"github.com/qgymje/nsq_copy/internal/dirlock"
	"github.com/qgymje/nsq_copy/internal/util"
)

const (
	TLSNotRequired = iota
	TLSRequiredExceptHTTP
	TLSRequired
)

type errStore struct {
	err error
}

// NSQD 表示一个nsqd实例对象，在其它对象里以ctx被引用
type NSQD struct {
	// 如果是一个计算器，请放在第一位, 为了兼容32位系统
	clientIDSequence int64

	sync.RWMutex

	opts atomic.Value // Store

	dl        *dirlock.DirLock
	isLoading int32
	errValue  atomic.Value
	startTime time.Time

	topicMap map[string]*Topic

	lookupPeers atomic.Value

	tcpListener   net.Listener
	httpListener  net.Listener
	httpsListener net.Listener
	tlsConfig     *tls.Config

	poolSize int

	idChan               chan MessageID
	notifyChan           chan interface{}
	optsNotificationChan chan struct{}
	exitChan             chan int
	waitGroup            util.WaitGroupWrapper

	ci *clusterinfo.ClusterInfo
}

// New 在nsqd启动的时候被调用，生成一个新的nsqd实例
// Options 是启动时设置的参数，nsqd的默认参数就在这个对象里设置
func New(opts *Options) *NSQD {
	// 设计disk_queue的地址，如果没有设置，则为启动实例的目录
	dataPath := opt.DataPath
	if opts.DataPath == "" {
		cwd, _ := os.Getwd()
		dataPath = cwd
	}

	n := &NSQDD{
		startTime:            time.Now(),
		topicMap:             make(map[string]*Topic),
		idChan:               make(chan MessageID, 4096), // messageid 列表只有4096个值，通过uuid生成, 是一个id生成器
		exitChan:             make(chan int),
		notifyChan:           make(chan interface{}),
		optsNotificationChan: make(chan struct{}, 1), //此处为何要带一个buffer?
		//ci: clusterinfo.New(opts.Logger, http_api.NewClient(nil, opts.HTTPClientConnectTimeout, opts.HTTPClientRequestTimeout)),
		dl: dirlock.New(dataPath),
	}
	n.swapOpts(opts)
}

func (n *NSQD) getOpts() *Options {
	return n.opts.Load().(*Options)
}

func (n *NSQD) swapOpts(opt *Options) {
	// Store 是atomic.Value的一个方法
	// 设置在n.opts
	// 对应的方法是Load
	n.opts.Store(opts)
}
