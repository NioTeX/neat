package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "log"
    "net/http"
)

type Message struct {
    Body string
}

func main() {
    api := rest.NewApi()
    api.Use(rest.DefaultDevStack...)
    router, err := rest.MakeRouter(
        rest.Get("/run", func(w rest.ResponseWriter, req *rest.Request) {
			asd, _ := activate()
            w.WriteJson("asd")
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
    api.SetApp(router)
    log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

func activate() (exp expos, err error) {
	ctx := starter.NewContext(&Evaluator{})
	if exp, err := starter.NewExperiment(ctx, ctx, i); err != nil {
		return nil, err
	} else {
		return exp, nil
	}
}
