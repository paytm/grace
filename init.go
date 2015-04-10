package grace

import (
 "net"
 "net/http"
 "log"
 "strconv"
 "os"
 "os/signal"
 "syscall"
)

func Serve(hport string, handler http.Handler) (error) {
  var l net.Listener

  fd := os.Getenv("EINHORN_FDS")
  if (fd != "") {
    sock,err := strconv.Atoi(fd)
    if (err == nil) {
      log.Println("detected socketmaster, listening on",fd)
      file := os.NewFile(uintptr(sock), "listener")
      fl, err := net.FileListener(file)
      if (err == nil) {
        l = fl
      }
    }
  }

  if l == nil {
    var err error
    l,err = net.Listen("tcp4",hport)
    if err != nil {
      log.Fatal(err)
    }
  }

  c:= make(chan os.Signal, 1)
  signal.Notify(c, syscall.SIGTERM) // listen for term
  go sigHandler(c,l)
  // TODO: Create a new server and call serve on it instead of 
  // http.Serve, setting connstate to enable really graceful shutdown
  return http.Serve(l,handler)
}

func sigHandler(c chan os.Signal, l net.Listener) {
  // Block until a signal is received.
  for _ = range c {
    log.Println("Terminating on SIGTERM", os.Getpid())
    l.Close() // FIXME: do via connstate
  }
}
