package nsqd

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qgymje/nsq_copy/internal/clusterinfo"
	"github.com/qgymje/nsq_copy/internal/dirlock"
	"github.com/qgymje/nsq_copy/internal/statsd"
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
	errValue  atomic.Value // 用于保存errStore
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
	dataPath := opts.DataPath
	if opts.DataPath == "" {
		cwd, _ := os.Getwd()
		dataPath = cwd
	}

	n := &NSQD{
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
	n.errValue.Store(errStore{})

	// 上来就锁目录
	err := n.dl.Lock()
	if err != nil {
		n.logf("FATAL: --data-path=%s in use (possibly by another instance of nsqd)", dataPath)
		// 参数/配置错误直接退出
		os.Exit(1)
	}

	// default is 6
	if opts.MaxDeflateLevel < 1 || opts.MaxDeflateLevel > 9 {
		n.logf("FATAL: --max-deflate-level must be [1,9]")
		os.Exit(1)
	}

	if opts.ID < 0 || opts.ID >= 1024 {
		n.logf("FATAL: --worker-id must be [0,1024]")
		os.Exit(1)
	}

	// 开始射击到一些statics的数据, 这是设计一个服务端软件应该有的模块
	if opts.StatsdPrefix != "" {
		var port string
		_, port, err = net.SplitHostPort(opts.HTTPAddress)
		if err != nil {
			n.logf("ERROR: failed to parse HTTP address (%s) - %s", opts.HTTPAddress, err)
			os.Exit(1)
		}
		// BroadcastAddress 默认为hostname, 也就是本机
		statsdHostKey := statsd.HostKey(net.JoinHostPort(opts.BroadcastAddress, port))
		// 默认为nsq.%s, 如果指定了,则变成 任意一个hostname.%s, 这样适合做多个不同机器的状态数据分类汇总
		// 这个需要在写一些服务统计时候需要注意的
		// 哪怕是一个简单的worker, 最好也要提供http端口, 用于查看状态
		prefixWithHost := strings.Replace(opts.StatsdPrefix, "%s", statsdHostKey, -1)
		if prefixWithHost[len(prefixWithHost)-1] != '.' {
			prefixWithHost += "."
		}
		opts.StatsdPrefix = prefixWithHost
	}

	return n
}

func (n *NSQD) logf(f string, args ...interface{}) {
	// 默认的logger在NewOptions里设置了,为std log
	// log.New(os.Stderr, "[nsqd]", log.Ldate|log.Ltime|log.Lmicroseconds)
	if n.getOpts().Logger == nil {
		return
	}
	n.getOpts().Logger.Output(2, fmt.Sprintf(f, args...))
}

func (n *NSQD) getOpts() *Options {
	return n.opts.Load().(*Options)
}

func (n *NSQD) swapOpts(opts *Options) {
	// Store 是atomic.Value的一个方法
	// 设置在n.opts
	// 对应的方法是Load
	n.opts.Store(opts)
}

func (n *NSQD) Main() {

}
