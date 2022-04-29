package pkg

import (
	"io/ioutil"
	"net/http"
)

func Client() {
	server := http.NewServeMux()
	server.HandleFunc("/transfer", func(w http.ResponseWriter, req *http.Request) {
		bye, _ := ioutil.ReadAll(req.Body)
		w.Write(bye)
	})
	err := http.ListenAndServe(":80", server)
	if err != nil {
		return
	}
}
