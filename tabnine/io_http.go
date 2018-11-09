package tabnine

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
)

const (
	httpProxyPortFilename = "httpProxy.port"
)

type HTTPClient struct {
	addr string
}

type HTTPServerConfig struct {
	TabnineBin string
	ConfigDir  string
}

// TODO(leeola): improve ionfigocess pid / port recording,
// plain files offer no concurrent safety and are a bit meh.
type HTTPServer struct {
	tabnineBin    string
	configDir     string
	tabnineStdin  io.Writer
	tabnineStdout *bufio.Reader
}

func NewHTTPClient(configDir string) (HTTPClient, error) {
	p := filepath.Join(configDir, httpProxyPortFilename)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return HTTPClient{}, fmt.Errorf("writefile: %v", err)
	}

	addr := fmt.Sprintf("localhost:%s", string(b))

	return HTTPClient{addr: addr}, nil
}

func NewHTTPServer(c HTTPServerConfig) (HTTPServer, error) {
	return HTTPServer{
		tabnineBin: c.TabnineBin,
		configDir:  c.ConfigDir,
	}, nil
}

func (c HTTPClient) SendRecv(req io.Reader) (io.ReadCloser, error) {
	res, err := http.Post(c.addr, "", req)
	if err != nil {
		return nil, fmt.Errorf("post: %v", err)
	}
	return res.Body, nil
}

func (h *HTTPServer) ListenAndServe(addr string) error {
	cmd := exec.Command(h.tabnineBin)

	// TODO(leeola): defer closure of the pipe.
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdinpipe: %v", err)
	}

	// TODO(leeola): defer closure of the pipe.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdoutpipe: %v", err)
	}

	h.tabnineStdin = stdin
	h.tabnineStdout = bufio.NewReader(stdout)

	go func() {
		log.Printf("%s starting..", h.tabnineBin)
		if err := cmd.Run(); err != nil {
			log.Printf("%s error: %v", h.tabnineBin, err)
		}
		log.Printf("%s exited", h.tabnineBin)
	}()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("create listener: %v", err)
	}
	defer listener.Close()

	if err := h.writePort(listener); err != nil {
		return fmt.Errorf("writeport: %v", err)
	}

	log.Printf("listening on %s..", listener.Addr().String())
	return http.Serve(listener, http.HandlerFunc(h.handler))
}

func (h HTTPServer) writePort(l net.Listener) error {
	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return fmt.Errorf("splithostport: %v", err)
	}

	p := filepath.Join(h.configDir, httpProxyPortFilename)
	err = ioutil.WriteFile(p, []byte(port), 0644)
	if err != nil {
		return fmt.Errorf("writefile: %v", err)
	}

	return nil
}

func (h HTTPServer) handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v %v", r.Method, r.URL, r.Proto)
	defer r.Body.Close()

	// ignore non-post/put methods.
	if r.Method != "POST" && r.Method != "PUT" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if _, err := io.Copy(h.tabnineStdin, r.Body); err != nil {
		log.Printf("copy http body: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err := h.tabnineStdin.Write([]byte{10}); err != nil {
		log.Printf("write delim byte error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := h.tabnineStdout.ReadBytes(10)
	if err != nil {
		log.Printf("read stdout: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, bytes.NewReader(b)); err != nil {
		log.Printf("copy response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
