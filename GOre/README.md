--GOre RAT--
GOre is a personal project to make a RAT/C2 in Golang with some C scripts if I want. Basically, it will be an attempt to learn go through a super fun
project. GOre is a portmanteu of "Go rev" or "go reverse".

--TODO--
to better faciliate an autonomous trojan, I want it to be equiped with everything it needs. This will bloat the binary, but it will assure me that I have
everything I need. Here are the things I want to be certain the RAT can do:

client
-C2 communication should be encrypted over TCP. The algo doesn't have to be unbreakable, just enough to throw off an analyst.
-port scanner, use blackhat GO's
-file curl capabilites; exfiltrate data over TCP
-Long sleeps; the RAT waits on server connections but should be asleep when not in use. maybe use polling/select, or lightweight threads? windows defender can see if a thread is
 opened by a process, but the false positive rates are absurdley high; ie when a web browser opens up a thread to do any kind of browsing.
-self-destruct; delete binary of self.
~process hollowing; this will be dependent on each system, so I'll need to be careful with implimentation.
-another idea, maybe try adding a Laplace client binary as an embedded file; I could see if that increases or reduces detection rates (probably an increase). I could also see if Laplace
 would be a good cover for GOre. 

server
-docker + nginx for communication. maybe I should do some UI design? Do a TUI if I have time.


