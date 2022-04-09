package web

import (
	"distributeKV/config"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
)

type ProxyServer struct {
	PMap map[uint]string
	Host string
}

func CreatProxyServer(partition *config.Partitions, host string) *ProxyServer {
	pMap := make(map[uint]string)
	for _, s := range partition.Partitions {
		pMap[uint(s.Index)] = s.Host
	}
	return &ProxyServer{
		PMap: pMap,
		Host: host,
	}
}

func (p *ProxyServer) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	// calculate the hash of the key
	h := fnv.New64()
	h.Write([]byte(key))

	idx := uint(h.Sum64()) % uint(len(p.PMap))
	fmt.Printf("The hash is %d and the idex is %d", int(h.Sum64()), idx)
	if targetHost, found := p.PMap[idx]; found {
		redirect(targetHost, w, r)
	} else {
		fmt.Printf("Not found key %s, index = %d in partition map\n", key, idx)
	}
}

func (p *ProxyServer) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	// calculate the hash of the key
	h := fnv.New64()
	h.Write([]byte(key))
	idx := uint(h.Sum64()) % uint(len(p.PMap))
	fmt.Printf("The hash is %d and the idex is %d", int(h.Sum64()), idx)
	if targetHost, found := p.PMap[idx]; found {
		redirect(targetHost, w, r)
	} else {
		fmt.Printf("Not found key %s, index = %d in partition map\n", key, idx)
	}
}

func (p *ProxyServer) ListenAndServeProxy() error {
	return http.ListenAndServe("127.0.0.1:8888", nil)
}

func redirect(host string, w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)
	res, err := http.Get("http://" + host + r.RequestURI)
	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("Error routing request to partiton %s, %v\n", host, err)
		return
	}
	defer res.Body.Close()
	fmt.Printf("Routing to partition %s\n...", host)
	io.Copy(w, res.Body)
}
