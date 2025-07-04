package main

import (
	"fmt"
	"github.com/friedHDD/Bedrock/server"
)

func main() {
	fmt.Printf(`
 ______  _______ ______   ______  _____  _______ _     _
 |_____] |______ |     \ |_____/ |     | |       |____/ 
 |_____] |______ |_____/ |    \_ |_____| |_____  |    \_
                                                        
A privacy-first file manager
===================
`)
	server.Start()

}
