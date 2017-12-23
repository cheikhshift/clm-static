# clm-static
A simple load manager written in Go. CLM auto scales your distributes loads based on the flags supplied. It uses the norms of DHCP to predict the addresses of your replicated servers.

## Installing
Downloading and installing this package requires Go.

### CMD flags
	$ clm -h

	  -app string
	    	Run specified terminal command each time a new instance is needed.
	    	
	  -appPort string
	    	Port your instances will listen on. (default "8080")
	  -incby int
	    	Increase last octet of IP when picking the next server. (default 1)
	    	
	  -max int
	    	Maximum number of connections per instance. clm-static will divide tasks. (default 100)
	    	
	  -net string
	    	Router LAN IP address (DHCP subnet). (default "192.168.0.1") 
	    	
	  -port string
	    	Port clm-static should listen on. (default "9000") 
	    	
	  -start int
	    	First DHCP assigned IP of your instances (Initial value of last octet in IPv4 address). Example : with value 21, this tool will assume your first instance's ip will be 192.168.0.21 (default 1)

### Requirements
Replicated virtual or physical machines with the software you wish to scale.
