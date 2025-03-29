package main

import (
	"apisuario/core"
	"apisuario/cmd/infraestructure/routes"
	"log"
)

func main() {
    db := core.InitDB()
    router := routes.SetupRouter(db)

    log.Println("Server running on port 5000")
    log.Fatal(router.Run(":5000"))
}
