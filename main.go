/*
 * Copyright 2018 Johannes Donath <johannesd@torchmind.com>
 * and other copyright owners as documented in the project's IP log.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
  "flag"
  "fmt"
  log "github.com/hashicorp/go-hclog"
  "golang.org/x/sync/errgroup"
  "io"
  "io/ioutil"
  "net"
  "os"
  "strconv"
  "strings"
  "time"
)

const versionNumber = "0.2.0"

func main() {
  var help bool
  var version bool

  var introductionFile string
  var serverName string
  var target string

  flag.BoolVar(&help, "help", false, "displays this help message")
  flag.BoolVar(&version, "version", false, "displays the application version number")

  flag.StringVar(&introductionFile, "introduction-file", "", "specifies a file to read an introduction (or explanation) message from")
  flag.StringVar(&serverName, "server-name", "irc.example.org", "specifies the server's display host name")
  flag.StringVar(&target, "target", "", "specifies the target server")
  flag.Parse()

  if help {
    printUsage(os.Stdout)
    return
  }
  if version {
    fmt.Printf("irc-redirect v%s\n", versionNumber)
    fmt.Printf("Licensed under the Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0.txt>\n")
    return
  }

  if flag.NArg() == 0 {
    fmt.Fprintf(os.Stderr, "error: at least one listener is required\n\n")
    printUsage(os.Stderr)
    os.Exit(1)
  }

  var targetHost string
  var targetPort uint16
  if target != "" {
    i := strings.IndexRune(target, ':')
    if i != -1 {
      targetHost = target[:i]
      if port, err := strconv.ParseUint(target[i+1:], 10, 16); err == nil {
        targetPort = uint16(port)
      } else {
        fmt.Fprintf(os.Stderr, "error: illegal target server port: %s", err)
        os.Exit(1)
      }
    } else {
      targetHost = target
      targetPort = 6667
    }
  }

  var introduction []string
  if introductionFile != "" {
    f, err := ioutil.ReadFile(introductionFile)
    if err != nil {
      fmt.Fprintf(os.Stderr, "error: cannot read introduction file: %s\n", err)
      os.Exit(1)
    }

    txt := string(f)
    txt = strings.Replace(txt, "\r", "", -1)
    introduction = strings.Split(txt, "\n")
  }

  var g errgroup.Group
  for _, addr := range flag.Args() {
    g.Go(func() error {
      listener, err := net.Listen("tcp", addr)
      if err != nil {
        return err
      }
      log.Default().Info("Listener started successfully", "addr", addr)

      for {
        conn, err := listener.Accept()
        if err != nil {
          log.Default().Warn("failed to accept incoming connection", "error", err)
          continue
        }

        go func() {
          defer conn.Close()

          log.Default().Debug("sending introduction", "client", conn.RemoteAddr())
          for _, l := range introduction {
            conn.Write([]byte(fmt.Sprintf(":%s NOTICE * :%s\r\n", serverName, l)))
          }

          if target != "" {
            log.Default().Debug("redirecting client", "client", conn.RemoteAddr(), "target", fmt.Sprintf("%s:%d", targetHost, targetPort))
            conn.Write([]byte(fmt.Sprintf(":%s 010 * %s %d :Port redirect\r\n", serverName, targetHost, targetPort)))
          }

          conn.Write([]byte("ERROR :This service is unavailable\r\n"))
          time.Sleep(time.Second)
        }()
      }
    })
  }

  err := g.Wait()
  if err != nil {
    fmt.Fprintf(os.Stderr, "error: failed to listen on one or more addresses: %s\n", err)
    os.Exit(2)
  }
}

func printUsage(w io.Writer) {
  fmt.Fprintf(w, "Usage: %s [flags] <listener1> [listener2] [listener3] ...\n\n", os.Args[0])
  fmt.Fprintf(w, "Available flags:\n")
  flag.CommandLine.VisitAll(func(flag *flag.Flag) {
    fmt.Fprintf(w, "  -%s\t%s", flag.Name, flag.Usage)
    if flag.DefValue != "" {
      fmt.Fprintf(w, " (default: %s)\n", flag.DefValue)
    } else {
      fmt.Fprint(w, "\n")
    }
  })
  fmt.Fprint(w, "For instance:\n\n", )
  fmt.Fprintf(w, "  $ %s -introduction-file=intro.txt -server-name=irc.example.org -target=irc.freenode.net:6669 0.0.0.0:6667\n\n", os.Args[0])
  fmt.Fprint(w, "Which will redirect all clients which connect via port 6667 (on all interfaces) to irc.freenode.net on port 6669.\n")
}
