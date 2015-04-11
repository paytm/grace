package grace

import (
 "net"
 "net/http"
 "log"
 "strconv"
 "time"
 "os"
 graceful "gopkg.in/tylerb/graceful.v1"
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

  srv := &graceful.Server{ 
	  Timeout: 10*time.Second, 
          Server: &http.Server{
	    Handler: handler,
          },
        }

  log.Println("starting serve on fd ",fd)
  return srv.Serve(l)
}
