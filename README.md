# pryct1-ST0257
This is a small, quick project made for my operating systems class. Its supposed to simulate a kernel talking to processes through messages. This simulation is achieved by passing messages through TCP sockets between certain processes.

###WARNING
This project really was more about a demonstration and learning Go than it was about making a robust program, it lacks a lot of polish. I'm not gonna work on it any further, but its still worth keeping around as reference material for future, more polished projects.


## HOW TO RUN IT
This project lacks a GUI, but it does programatically open terminal windows. However, you need to specify what terminal you'll be using, this program doesn't even try to use the `$SHELL` variable or anything else. To specify it, simply go to `src/simk.go`, there are two `const`s called `appClientCMD` and `fmClientCMD`. Leave the `go run X` part alone, but change the `st -e` part to match your terminal. 

Now you can run the project. To do so, open two terminals, in both of them cd into `src`. Simply type `go run simkmain.go` in one and `go run tcmain.go` in the other. 


## HOW TO USE IT
A `process` is simply one of the `*main.go` files running. They will only run if the "kernel" process (`simkmain.go`) is active. There are 3 types of `process`
1. USER - an instance of `tcmain.go` running, it's supposed to simulate a user
2. FM - an instance of `fmmain.go` running, it's supposed to be a barebones "file manager" that can only create and delete new folders. These are put into `fakefs/`. You can only have one of these open at a time.
3. PROC - an instance of `acmain.go` running, it's supposed to be a generic application that can recieve messages

Each process has a unique number. Process 0 is always the "kernel". Process 1 is always the FM and it will not be running at startup. The user interface is pretty self explanatory, except for the `4. send message to proc` option. After typing your destination, you can type anything. If you direct it at a PROC, nothing will happen other than it possibly printing on that terminal and getting back an `OK` or an `ERR`. If you direct it at the FM you get some commands:
- `cr [dirname]` will create a folder `dirname`
- `rm [dirname]` will delete the folder `dirname`
- `lg [num]` (where num is a number) will return the last `num` logged actions

Anything else will give you a parse error. 

To avoid clutter, the system does NOT try to log OK, ERR, or IDEN messages. It stores the logs in `fakefs/logFile`. 
