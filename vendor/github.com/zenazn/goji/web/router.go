package web

import (
	"net/http"
	"sort"
	"strings"
	"sync"
)

type method int

const (
	mCONNECT method = 1 << iota
	mDELETE
	mGET
	mHEAD
	mOPTIONS
	mPATCH
	mPOST
	mPUT
	mTRACE
	// We only natively support the methods above, but we pass through other
	// methods. This constant pretty much only exists for the sake of mALL.
	mIDK

	mALL method = mCONNECT | mDELETE | mGET | mHEAD | mOPTIONS | mPATCH |
		mPOST | mPUT | mTRACE | mIDK
)

// The key used to communicate to the NotFound handler what methods would have
// been allowed if they'd been provided.
const ValidMethodsKey = "goji.web.ValidMethods"

var validMethodsMap = map[string]method{
	"CONNECT": mCONNECT,
	"DELETE":  mDELETE,
	"GET":     mGET,
	"HEAD":    mHEAD,
	"OPTIONS": mOPTIONS,
	"PATCH":   mPATCH,
	"POST":    mPOST,
	"PUT":     mPUT,
	"TRACE":   mTRACE,
}

type route struct {
	prefix  string
	method  method
	pattern Pattern
	handler Handler
}

type router struct {
	lock     sync.Mutex
	routes   []route
	notFound Handler
	machine  *routeMachine
}

func httpMethod(mname string) method {
	if method, ok := validMethodsMap[mname]; ok {
		return method
	}
	return mIDK
}

func (rt *router) compile() *routeMachine {
	rt.lock.Lock()
	defer rt.lock.Unlock()
	sm := routeMachine{
		sm:     compile(rt.routes),
		routes: rt.routes,
	}
	rt.setMachine(&sm)
	return &sm
}

func (rt *router) getMatch(c *C, w http.ResponseWriter, r *http.Request) Match {
	rm := rt.getMachine()
	if rm == nil {
		rm = rt.compile()
	}

	methods, route := rm.route(c, w, r)
	if route != nil {
		return Match{
			Pattern: route.pattern,
			Handler: route.handler,
		}
	}

	if methods == 0 {
		return Match{Handler: rt.notFound}
	}

	var methodsList = make([]string, 0)
	for mname, meth := range validMethodsMap {
		if methods&meth != 0 {
			methodsList = append(methodsList, mname)
		}
	}
	sort.Strings(methodsList)

	if c.Env == nil {
		c.En