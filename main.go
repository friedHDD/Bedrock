package main

import (
	"fmt"
	"github.com/friedHDD/Bedrock/core"
	"github.com/friedHDD/Bedrock/server"
	"log"
)

func main() {
	fmt.Printf(`
 ______  _______ ______   ______  _____  _______ _     _
 |_____] |______ |     \ |_____/ |     | |       |____/ 
 |_____] |______ |_____/ |    \_ |_____| |_____  |    \_
                                                        
A privacy-first file manager
===================
`)
	err := core.InitAll()
	if err != nil {
		log.Fatal(err)
	}
	server.Start()

}
