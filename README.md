# vlog
Basic Logger for GoLang. Has ability to wrap/rotate.

1) vlog provides two levels of log, Info and Debug. 
2) Only one level is active at a time.
3) The debug level provides additional information about the 
   caller like the package name and function name.
4) Log name provided by the user is suffixed with "_ nn.log"
5) The log files created have a size around 2GB.
6) New files are created with the number part incremented in the suffix.
7) The Debug level has the suffix "_ debug _ nn.log"
8) If the user process using this facility goes down, the logging
   restarts where it left.
9) Even though the user program can have multiple vlog.Info or vlog.Debug
   statements, only the logging level set passed is considered for logging.
   
   Please contact vishal.pandya@gmail.com if you see any issues or have comments.


