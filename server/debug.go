package server

import "log"

var Debug = debug{}

type debug struct{}

func (debug) PrintTypMap(r *Router) {
	for typid := range r.m {
		log.Printf("typ id: %x, typ name: %v\n", typid, md5Typ[typid])
	}
}
