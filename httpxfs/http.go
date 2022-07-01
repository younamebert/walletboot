// package httpxfs

// import (
// 	"bufio"
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"io/ioutil"
// 	"net/http"
// 	"time"
// )

// type Client struct {
// 	hostUrl string
// 	timeOut string
// }

// type jsonRPCReq struct {
// 	JsonRPC string      `json:"jsonrpc"`
// 	ID      int         `json:"id"`
// 	Method  string      `json:"method"`
// 	Params  interface{} `json:"params"`
// }

// type jsonRPCResp struct {
// 	JSONRPC string      `json:"jsonrpc"`
// 	Result  interface{} `json:"result"`
// 	Error   *RPCError   `json:"error"`
// 	ID      int         `json:"id"`
// }

// var ErrtoSliceStr = errors.New("json: cannot unmarshal string into Go value of type []string")

// func NewClient(url, timeOut string) *Client {
// 	return &Client{
// 		hostUrl: url,
// 		timeOut: timeOut,
// 	}
// }

// // CallMethod executes a JSON-RPC call with the given psrameters,which is important to the rpc server.
// func (cli *Client) CallMethod(id int, methodname string, params interface{}, out interface{}) error {

// 	timeDur, err := time.ParseDuration(cli.timeOut)
// 	if err != nil {
// 		return err
// 	}

// 	client := &http.Client{Timeout: timeDur}

// 	req := &jsonRPCReq{
// 		JsonRPC: "2.0",
// 		ID:      id,
// 		Method:  methodname,
// 		Params:  params,
// 	}

// 	reqStr, err := json.Marshal(req)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := client.Post(cli.hostUrl, "application/json;charset=utf-8", bytes.NewBuffer(reqStr))
// 	if err != nil {
// 		return err
// 	}
// 	defer result.Body.Close()
// 	// The result must be a pointer so that response json can unmarshal into it.

// 	content, err := ioutil.ReadAll(result.Body)
// 	if err != nil {
// 		return err
// 	}

// 	bs, err := json.Marshal(content)
// 	if err != nil {
// 		return err
// 	}
// 	err = json.Unmarshal(bs, out)
// 	if err != nil {
// 		if err.Error() == ErrtoSliceStr.Error() {
// 			sentences := []string{}
// 			scanner := bufio.NewScanner(bytes.NewBuffer(bs))
// 			for scanner.Scan() {
// 				sentences = append(sentences, scanner.Text())
// 			}
// 			bs, err := json.Marshal(sentences)
// 			if err != nil {
// 				return err
// 			}
// 			if err := json.Unmarshal(bs, out); err != nil {
// 				return err
// 			}
// 			return nil
// 		}
// 		return err
// 	}

// 	return nil
// }
// Copyright 2018 The xfsgo Authors
// This file is part of the xfsgo library.
//
// The xfsgo library is free software: you can redistribute it and/or modify
// it under the terms of the MIT Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The xfsgo library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// MIT Lesser General Public License for more details.
//
// You should have received a copy of the MIT Lesser General Public License
// along with the xfsgo library. If not, see <https://mit-license.org/>.

package httpxfs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
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
	Error   rpcError    `json:"error"`
	ID      int         `json:"id"`
}

func NewClient(url, timeOut string) *Client {
	return &Client{
		hostUrl: url,
		timeOut: timeOut,
	}
}

// CallMethod executes a JSON-RPC call with the given psrameters,which is important to the rpc server.
func (cli *Client) CallMethod(id int, methodname string, params interface{}, out interface{}) error {
	client := resty.New()

	timeDur, err := time.ParseDuration(cli.timeOut)
	if err != nil {
		return err
	}
	client = client.SetTimeout(timeDur)
	req := &jsonRPCReq{
		JsonRPC: "2.0",
		ID:      id,
		Method:  methodname,
		Params:  params,
	}
	// The result must be a pointer so that response json can unmarshal into it.
	var resp *jsonRPCResp = nil
	r, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&resp). // or SetResult(AuthSuccess{}).
		Post(cli.hostUrl)
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("resp null")
	}

	e := resp.Error.Message
	if e != "" {
		return fmt.Errorf(e)
	}

	js, err := json.Marshal(resp.Result)
	if err != nil {
		return err
	}
	err = json.Unmarshal(js, out)
	if err != nil {
		return err
	}
	_ = r
	return nil
}
