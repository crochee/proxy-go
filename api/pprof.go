// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/8

package api

import (
	"html/template"
	"net/http"
	"proxy-go/logger"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"proxy-go/api/response"
)

var profileDescriptions = map[string]string{
	"allocs":       "A sampling of all past memory allocations",
	"block":        "Stack traces that led to blocking on synchronization primitives",
	"cmdline":      "The command line invocation of the current program",
	"goroutine":    "Stack traces of all current goroutines",
	"heap":         "A sampling of memory allocations of live objects. You can specify the gc GET parameter to run GC before taking the heap sample.",
	"mutex":        "Stack traces of holders of contended mutexes",
	"profile":      "CPU profile. You can specify the duration in the seconds GET parameter. After you get the profile file, use the go tool pprof command to investigate the profile.",
	"threadcreate": "Stack traces that led to the creation of new OS threads",
	"trace":        "A trace of execution of the current program. You can specify the duration in the seconds GET parameter. After you get the trace file, use the go tool trace command to investigate the trace.",
}

var indexTmpl = template.Must(template.New("index").Parse(`<html>
<head>
<title>/debug/pprof/</title>
<style>
.profile-name{
	display:inline-block;
	width:6rem;
}
</style>
</head>
<body>
/debug/pprof/<br>
<br>
Types of profiles available:
<table>
<thead><td>Count</td><td>Profile</td></thead>
{{range .}}
	<tr>
	<td>{{.Count}}</td><td><a href={{.Href}}>{{.Name}}</a></td>
	</tr>
{{end}}
</table>
<a href="goroutine?debug=2">full goroutine stack dump</a>
<br/>
<p>
Profile Descriptions:
<ul>
{{range .}}
<li><div class=profile-name>{{.Name}}:</div> {{.Desc}}</li>
{{end}}
</ul>
</p>
</body>
</html>
`))

// Profile godoc
// @Summary pprof profile
// @Description get pprof profile
// @Tags pprof
// @Accept application/json
// @Produce application/octet-stream
// @Param seconds query string false "second default(30)"
// @Success 200
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /debug/pprof/index [get]
func Index(ctx *gin.Context) {
	ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")

	logger.FromContext(ctx.Request.Context()).Debug("start index")
	type profile struct {
		Name  string
		Href  string
		Desc  string
		Count int
	}
	var profiles []profile
	for _, p := range pprof.Profiles() {
		profiles = append(profiles, profile{
			Name:  p.Name(),
			Href:  p.Name() + "?debug=1",
			Desc:  profileDescriptions[p.Name()],
			Count: p.Count(),
		})
	}

	// Adding other profiles exposed from within this package
	for _, p := range []string{"cmdline", "profile", "trace"} {
		profiles = append(profiles, profile{
			Name: p,
			Href: p,
			Desc: profileDescriptions[p],
		})
	}

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := indexTmpl.Execute(ctx.Writer, profiles); err != nil {
		ctx.Writer.Header().Del("Content-Type")
		ctx.Writer.Header().Set("X-Go-Pprof", "1")
		ctx.Writer.Header().Del("Content-Disposition")
		response.GinError(ctx, err)
	}
}

// Profile godoc
// @Summary pprof profile
// @Description get pprof profile
// @Tags pprof
// @Accept application/json
// @Produce application/octet-stream
// @Param seconds query string false "second default(30)"
// @Success 200
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /debug/pprof/profile [get]
func Profile(ctx *gin.Context) {
	ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	sec := ctx.GetFloat64("seconds")
	if sec <= 0 {
		sec = 30
	}

	if durationExceedsWriteTimeout(ctx.Request, sec) {
		response.GinError(ctx, response.Error(http.StatusBadRequest,
			"profile duration exceeds server's WriteTimeout"))
		return
	}

	// Set Content Type assuming StartCPUProfile will work,
	// because if it does it starts writing.
	ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
	ctx.Writer.Header().Set("Content-Disposition", `attachment; filename="cpu.prof"`)
	if err := pprof.StartCPUProfile(ctx.Writer); err != nil {
		// StartCPUProfile failed, so no writes yet.
		ctx.Writer.Header().Del("Content-Type")
		ctx.Writer.Header().Del("Content-Disposition")
		response.GinError(ctx, response.ErrorWiths(http.StatusInternalServerError,
			"Could not enable CPU profiling", err))
		return
	}
	sleep(ctx.Writer, time.Duration(sec)*time.Second)
	pprof.StopCPUProfile()
}

// Trace godoc
// @Summary pprof trace
// @Description get pprof trace
// @Tags pprof
// @Accept application/json
// @Produce application/octet-stream
// @Param seconds query string false "second default(30)"
// @Success 200
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /debug/pprof/trace [get]
func Trace(ctx *gin.Context) {
	ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	sec := ctx.GetFloat64("seconds")
	if sec <= 0 {
		sec = 30
	}

	if durationExceedsWriteTimeout(ctx.Request, sec) {
		response.GinError(ctx, response.Error(http.StatusBadRequest,
			"profile duration exceeds server's WriteTimeout"))
		return
	}

	// Set Content Type assuming trace.Start will work,
	// because if it does it starts writing.
	ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
	ctx.Writer.Header().Set("Content-Disposition", `attachment; filename="trace.prof"`)
	if err := trace.Start(ctx.Writer); err != nil {
		// trace.Start failed, so no writes yet.
		ctx.Writer.Header().Del("Content-Type")
		ctx.Writer.Header().Del("Content-Disposition")
		response.GinError(ctx, response.ErrorWiths(http.StatusInternalServerError,
			"Could not enable tracing", err))

		return
	}
	sleep(ctx.Writer, time.Duration(sec*float64(time.Second)))
	trace.Stop()
}

// Heap godoc
// @Summary pprof heap
// @Description get pprof heap
// @Tags pprof
// @Accept application/json
// @Produce application/octet-stream
// @Success 200
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /debug/pprof/heap [get]
func Heap(ctx *gin.Context) {
	ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	// Set Content Type assuming trace.Start will work,
	// because if it does it starts writing.
	debug := ctx.GetInt("debug")
	if debug != 0 {
		ctx.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	} else {
		ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
		ctx.Writer.Header().Set("Content-Disposition", `attachment; filename="mem.prof"`)
	}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(ctx.Writer); err != nil {
		ctx.Writer.Header().Del("Content-Type")
		ctx.Writer.Header().Del("Content-Disposition")
		response.GinError(ctx, response.ErrorWiths(http.StatusInternalServerError,
			"could not write memory profile", err))
	}
}

func sleep(w http.ResponseWriter, d time.Duration) {
	clientGone := make(<-chan bool)
	if cn, ok := w.(http.CloseNotifier); ok {
		clientGone = cn.CloseNotify()
	}
	select {
	case <-time.After(d):
	case <-clientGone:
	}
}

func durationExceedsWriteTimeout(r *http.Request, seconds float64) bool {
	srv, ok := r.Context().Value(http.ServerContextKey).(*http.Server)
	return ok && srv.WriteTimeout != 0 && seconds >= srv.WriteTimeout.Seconds()
}
