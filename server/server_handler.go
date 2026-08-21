package server


type ServerOpts struct{
	Address string
	MaxSession int
	// these are fine for now ig
}


func StartServer(opts *ServerOpts){
	listener, err := net.Listen("tcp", opts.Address);
	
}
