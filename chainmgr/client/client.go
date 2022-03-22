package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	hostUrl string
	timeOut string
}

type jsonRPCReq struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type jsonRPCResp struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *RPCError   `json:"error"`
	ID      int         `json:"id"`
}

var ErrtoSliceStr = errors.New("json: cannot unmarshal string into Go value of type []string")

func NewClient(url, timeOut string) *Client {
	return &Client{
		hostUrl: url,
		timeOut: timeOut,
	}
}

// CallMethod executes a JSON-RPC call with the given psrameters,which is important to the rpc server.
func (cli *Client) CallMethod(id int, methodname string, params interface{}, out interface{}) error {

	timeDur, err := time.ParseDuration(cli.timeOut)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: timeDur}

	req := &jsonRPCReq{
		JsonRPC: "2.0",
		ID:      id,
		Method:  methodname,
		Params:  params,
	}

	reqStr, err := json.Marshal(req)
	if err != nil {
		return err
	}

	result, err := client.Post(cli.hostUrl, "application/json;charset=utf-8", bytes.NewBuffer(reqStr))
	if err != nil {
		return err
	}
	defer result.Body.Close()
	// The result must be a pointer so that response json can unmarshal into it.

	resp := make(map[string]interface{})

	content, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, &resp); err != nil {
		return err
	}

	// logrus.Infof("methodname:%v params:%v resp:%v\n", methodname, params, string(content))

	bsErr, err := json.Marshal(resp["error"])
	if err != nil {
		return err
	}

	if string(bsErr) != "null" {
		return fmt.Errorf(string(bsErr))
	}

	bs, err := json.Marshal(resp["result"])
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, out)
	if err != nil {
		if err.Error() == ErrtoSliceStr.Error() {
			sentences := []string{}
			scanner := bufio.NewScanner(bytes.NewBuffer(bs))
			for scanner.Scan() {
				sentences = append(sentences, scanner.Text())
			}
			bs, err := json.Marshal(sentences)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(bs, out); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	return nil
}
