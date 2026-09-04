SF - SSH Forwarder in Go


This is a simple program for self use to control my other laptop which
is in non-static IP conditions from my other devices.


The architecture is pretty simple, it utilises a central passthrough
server to pass commands and results back and forth accross the server.


This is the initial version and so it only allos for non-blocking commands
to be ran without any error.

TBH, full ssh capability is not too hard from here, but that will depend
upon my mood and how much do i end up using this.



For now, here is rough architecture:



<client>--------------------<central>---------------------<server>
	|                          |                            |
My main PC                    VPS                        side-PC


The central job is a continuously running process which also
holds the current command to be executed and waiting for return,
it also manages whether the server connection is ok or not via health
checks (not added completely yet).


About Testing:
I tested this system on common commands as well as non-blocking ones (made using
appending & to regular commands) -- it all appeared to work ig



