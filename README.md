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

