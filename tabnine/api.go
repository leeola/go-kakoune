package tabnine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	// Version seems to indicate a version of the TabNine API, of which this code was
	// tested with and so it hardcodes the version in use.
	Version = "0.4.0"
)

type request struct {
	Version string `json:"version"`
	Request struct {
		Prefetch     *PrefetchRequest     `json:"Prefetch,omitempty"`
		Autocomplete *AutocompleteRequest `json:"Autocomplete,omitempty"`
	} `json:"request"`
}

type PrefetchRequest struct {
	Filename string `json:"filename"`
}

type AutocompleteRequest struct {
	RegionIncludesEnd       bool   `json:"region_includes_end"`
	After                   string `json:"after"`
	Before                  string `json:"before"`
	RegionIncludesBeginning bool   `json:"region_includes_beginning"`
	Filename                string `json:"filename"`
	MaxNumResults           uint8  `json:"max_num_results"`
}

type AutocompleteResponse struct {
	SuffixToSubstitute string               `json:"suffix_to_substitute"`
	Results            []AutocompleteResult `json:"results"`
	IsActive           bool                 `json:"is_active"`
	// not seen a promo yet, not sure on the format, leaving it as bytes
	// so the raw output can be recorded.
	PromotionalMessage json.RawMessage `json:"promotional_message"`
}

type AutocompleteResult struct {
	Result             string `json:"result"`
	PrefixToSubstitute string `json:"prefix_to_substitute"`
}

// reference request:
// {
//   "version": "0.4.0",
//   "request": {
//     "Prefetch": {
//       "filename": "/foo/bar/baz.go"
//     }
//   }
// }
//
// no response for prefetches.
func (t Tabnine) Prefetch(req PrefetchRequest) error {
	treq := request{Version: Version}
	treq.Request.Prefetch = &req

	b, err := json.Marshal(treq)
	if err != nil {
		return fmt.Errorf("marshal: %v", err)
	}

	rc, err := t.io.SendRecv(bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("sendrecv: %v", err)
	}

	// no data to care about, close it immediately.
	//
	// Perhaps draining it first is better?
	return rc.Close()
}

// reference request:
// {
//   "version": "0.4.0",
//   "request": {
//     "Autocomplete": {
//       "region_includes_end": true,
//       "after": "\n}",
//       "before": "package main\n\nimport \"fmt\"\n\nfunc main() {\n fmt.Println(\"hello\")\n fmt.d",
//       "region_includes_beginning": true,
//       "filename": "/foo/bar/baz.go",
//       "max_num_results": 5
//     }
//   }
// }
//
// reference response:
// {
//   "suffix_to_substitute": "d",
//   "results": [
//     {
//       "result": "desc",
//       "prefix_to_substitute": ""
//     },
//     {
//       "result": "define",
//       "prefix_to_substitute": ""
//     },
//     {
//       "result": "define-command",
//       "prefix_to_substitute": ""
//     },
//     {
//       "result": "do",
//       "prefix_to_substitute": ""
//     },
//     {
//       "result": "doesn't",
//       "prefix_to_substitute": ""
//     }
//   ],
//   "is_active": false,
//   "promotional_message": []
// }
func (t Tabnine) Autocomplete(req AutocompleteRequest) (AutocompleteResponse, error) {
	treq := request{Version: Version}
	treq.Request.Autocomplete = &req

	b, err := json.Marshal(treq)
	if err != nil {
		return AutocompleteResponse{}, fmt.Errorf("marshal: %v", err)
	}

	rc, err := t.io.SendRecv(bytes.NewReader(b))
	if err != nil {
		return AutocompleteResponse{}, fmt.Errorf("sendrecv: %v", err)
	}
	defer rc.Close()

	b, err = ioutil.ReadAll(rc)
	if err != nil {
		return AutocompleteResponse{}, fmt.Errorf("readall: %v", err)
	}

	var res AutocompleteResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return AutocompleteResponse{}, fmt.Errorf("unmarshal: %v", err)
	}

	return res, nil
}
