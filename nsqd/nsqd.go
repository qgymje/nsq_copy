package nsqd

import (
	"crypto/tls"
	"net"
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

	opts atomic.Value

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
