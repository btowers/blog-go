package main

import (
	"github.com/btowers/blog-go/pkg/adder"
	"github.com/btowers/blog-go/pkg/auth"
	"github.com/btowers/blog-go/pkg/lister"
	"github.com/btowers/blog-go/pkg/remover"
	"github.com/btowers/blog-go/pkg/rest"
	"github.com/btowers/blog-go/pkg/storage/mongodb"
	"github.com/btowers/blog-go/pkg/updater"
)

func main() {
	// New Storage
	s := mongodb.NewStorage()

	// New Services
	auth := auth.NewService(s)
	add := adder.NewService(s)
	list := lister.NewService(s)
	remove := remover.NewService(s)
	update := updater.NewService(s)
	//	retrieve := retriever.NewService(s)

	// New Router
	router := rest.NewRouter(auth, add, list, remove, update)

	// Start Server on "localhost:8080"
	router.Run(":8080")

}
