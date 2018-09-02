This little go program permits the redirection of IRC users to another server (typically
to make users aware of the SSL compatible ports).

Usage
-----

```
Usage: irc-redirect [flags] <listener1> [listener2] [listener3] ...

Available Flags:

  -introduction-file        specifies a file to read an introduction (or explanation) message from
  -server-name              specifies the server's display host name
  -target-host              specifies the hostname to redirect to
  -target-port              specifies the port to redirect to
  
For instance:

  $ irc-redirect -introduction-file=intro.txt -server-name=irc.example.org -target-host=irc.freenode.net -target-port=6669 0.0.0.0:6667
  
Which will redirect all clients which connect via port 6667 (on all interfaces) to
irc.freenode.net on port 6669.
```

Installing from Source
----------------------

1. Make sure go is installed within your environment
2. Execute `go get github.com/dotStart/irc-redirect`

Building a full distribution Package
------------------------------------

1. Make sure go is installed within your environment
2. Execute `go get -d github.com/dotStart/irc-redirect`
3. Switch to the project dir (`cd $GOROOT/src/github.com/dotstart/irc-redirect`)
4. Execute `make`

License
-------

```
Copyright [year] [name] <[email]>
and other copyright owners as documented in the project's IP log.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
