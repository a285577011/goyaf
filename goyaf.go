package goyaf

import (
	"flag"
	"fmt"
	"git.oschina.net/pbaapp/goyaf/lib"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type Goyaf struct {
	version string
}

//运行环境
var env string

//设置运行环境
func SetEnv(e string) {
	env = e
}

//获取运行环境
func GetEnv() string {
	return env
}

//监听的文件句柄
var listenFD int
var currentPath string

func Run() {
	env := flag.String("env", "product", "environment")
	lfd := flag.Int("listenFD", 0, "the already-open fd to listen on (internal use only)")
	flag.Parse()

	SetEnv(*env)
	listenFD = *lfd

	runServer()
}

func runServer() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	mux := &GoyafMux{}

	var err error
	var l net.Listener
	server := &http.Server{
		Addr:    ":" + GetConfigByKey("http-listen-port"),
		Handler: mux,
	}
	if listenFD != 0 {
		Log("Listening to existing fd ", listenFD)
		f := os.NewFile(uintptr(listenFD), "listen socket")
		l, err = net.FileListener(f)
	} else {
		Log("Listening on", server.Addr)
		l, err = net.Listen("tcp", server.Addr)
	}
	if err != nil {
		log.Fatal(err)
	}

	theStoppable = newStoppable(l)

	err = server.Serve(theStoppable)
	if theStoppable.stopped {
		for i := 0; i < 10; i++ {
			if connCount.get() == 0 {
				continue
			}
			log.Print("waiting for clients...")
			time.Sleep(1 * time.Second)
		}

		if connCount.get() == 0 {
			log.Print("server gracefully stopped.")
			os.Exit(0)
		} else {
			log.Fatalf("server stopped after 10 seconds with %d clients still connected.", connCount.get())
		}
	}
	if err != nil {
		log.Fatal(err)
	}
}

func UpgradeServer(w http.ResponseWriter, req *http.Request) {
	var sig signal

	tl := theStoppable.Listener.(*net.TCPListener)
	fl, err := tl.File()
	if err != nil {
		log.Fatal(err)
	}
	fd := fl.Fd()

	noCloseOnExec(fd)

	cmd := exec.Command(currentPath,
		fmt.Sprintf("-listenFD=%d", fd),
		fmt.Sprintf("-env=%s", env))

	stderr := GetConfigByKey("stderr")
	if len(stderr) > 0 {
		strderrfile, err := os.OpenFile(stderr, os.O_APPEND|os.O_WRONLY, 0600)
		if os.IsNotExist(err) {
			os.Create(stderr)
			strderrfile, _ = os.OpenFile(stderr, os.O_APPEND|os.O_WRONLY, 0600)
		}
		if err != nil {
			Log(err)
		}
		cmd.Stderr = strderrfile
	} else {
		cmd.Stderr = os.Stderr
	}
	stdout := GetConfigByKey("stdout")
	if len(stdout) > 0 {
		strdoutfile, err := os.OpenFile(stdout, os.O_APPEND|os.O_WRONLY, 0600)
		if os.IsNotExist(err) {
			os.Create(stdout)
			strdoutfile, _ = os.OpenFile(stdout, os.O_APPEND|os.O_WRONLY, 0600)
		}
		if err != nil {
			Log(err)
		}
		cmd.Stdout = strdoutfile
	} else {
		cmd.Stdout = os.Stdout
	}

	log.Print("starting cmd: ", cmd.Args)
	if err := cmd.Start(); err != nil {
		log.Print("error:", err)
		return
	}

	theStoppable.stop <- sig
}

//错误处理控制器
var panicHandleController interface{}

func SetPanicHandleController(c interface{}) {
	panicHandleController = c
}

type counter struct {
	m sync.Mutex
	c int
}

func (c counter) get() (ct int) {
	c.m.Lock()
	ct = c.c
	c.m.Unlock()
	return
}

var connCount counter

type watchedConn struct {
	net.Conn
}

func (w watchedConn) Close() error {
	connCount.m.Lock()
	connCount.c--
	connCount.m.Unlock()

	return w.Conn.Close()
}

type signal struct{}

type stoppableListener struct {
	net.Listener
	stop    chan signal
	stopped bool
}

var theStoppable *stoppableListener

func newStoppable(l net.Listener) (sl *stoppableListener) {
	sl = &stoppableListener{Listener: l, stop: make(chan signal, 1)}
	go func() {
		_ = <-sl.stop
		sl.stopped = true
		sl.Listener.Close()
	}()
	return
}

func (sl *stoppableListener) Accept() (c net.Conn, err error) {
	c, err = sl.Listener.Accept()
	if err != nil {
		return
	}

	c = watchedConn{Conn: c}

	// Count it
	connCount.m.Lock()
	connCount.c++
	connCount.m.Unlock()

	return
}

// These are here because there is no API in syscall for turning OFF
// close-on-exec (yet).

// from syscall/zsyscall_linux_386.go, but it seems like it might work
// for other platforms too.
func fcntl(fd int, cmd int, arg int) (val int, err error) {
	r0, _, e1 := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), uintptr(cmd), uintptr(arg))
	val = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

func noCloseOnExec(fd uintptr) {
	fcntl(int(fd), syscall.F_SETFD, ^syscall.FD_CLOEXEC)
}

func init() {
	currentPath = lib.GetCurrentPath()
}
