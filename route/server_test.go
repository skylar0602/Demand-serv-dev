package route_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/smarterwallet/demand-abstraction-serv/model"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	client := http.Client{}
	t.Run("ping", func(t *testing.T) {
		resp, err := client.Get("http://127.0.0.1:8080/v1")
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("strategy not support", func(t *testing.T) {
		req := &model.DemandRequest{
			Model:  "v1",
			Demand: "reset my password",
		}
		buf, err := json.Marshal(req)
		assert.Nil(t, err)
		resp, err := client.Post("http://127.0.0.1:8080/v1/demand", "application/json", bytes.NewReader(buf))
		assert.Nil(t, err)
		assert.Equal(t, 500, resp.StatusCode)
	})
	t.Run("demand", func(t *testing.T) {
		req := &model.DemandRequest{
			Model:  "v1",
			Demand: "I want High return and low risk",
		}
		buf, err := json.Marshal(req)
		assert.Nil(t, err)
		resp, err := client.Post("http://127.0.0.1:8080/v1/demand", "application/json", bytes.NewReader(buf))
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("demand crosschain", func(t *testing.T) {
		req := &model.DemandRequest{
			Model:  "v1",
			Demand: "I want to transfer 1000 USDC to 0xd8da6bf26964af9d7eed9e03e53415d37aa96045 from Ethereum Goerli to Avalance",
		}
		buf, err := json.Marshal(req)
		assert.Nil(t, err)
		resp, err := client.Post("http://127.0.0.1:8080/v1/demand", "application/json", bytes.NewReader(buf))
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		res := &model.DemandResponse{}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(body, &res); err != nil {
			t.Fatal(err)
		}
		t.Logf("%+v", res.Detail)
	})
	t.Run("chat demand", func(t *testing.T) {
		cid, err := startChat()
		assert.Nil(t, err)
		err = chat(cid, "Current Goerli balance: 20USDC, Fuji balance: 90USDC. I want to transfer 100USDC to 0x5134F00C95b8e794db38E1eE39397d8086cee7Ed on Fuji")
		assert.Nil(t, err)
		err = chat(cid, "Current Goerli balance: 10USDC, Fuji balance: 100USDC. I want to transfer 90USDC to 0x5134F00C95b8e794db38E1eE39397d8086cee7Ed on Fuji")
		assert.Nil(t, err)
		err = chat(cid, "I want to transfer 5 more USDC")
		assert.Nil(t, err)
	})
	t.Run("no swap + no crosschain", func(t *testing.T) {
		cid, err := startChat()
		assert.Nil(t, err)
		err = initBalance(cid)
		assert.Nil(t, err)
		err = chat(cid, "I want to transfer 80 USDC to 0x5134F00C95b8e794db38E1eE39397d8086cee7Ed on mumbai")
		assert.Nil(t, err)
	})
	t.Run("no swap + no crosschain + stable", func(t *testing.T) {
		cid, err := startChat()
		assert.Nil(t, err)
		err = initBalance(cid)
		assert.Nil(t, err)
		err = chat(cid, "I want to transfer 80 dollar to 0x5134F00C95b8e794db38E1eE39397d8086cee7Ed on mumbai")
		assert.Nil(t, err)
	})
	t.Run("no swap + crosschain", func(t *testing.T) {
		cid, err := startChat()
		assert.Nil(t, err)
		err = initBalance(cid)
		assert.Nil(t, err)
		// todo target chain detect incorrect
		err = chat(cid, "I want to transfer 120USDC to 0x5134F00C95b8e794db38E1eE39397d8086cee7Ed on target chain fuji")
		// current usdc: 100
		// target usdc: 25
		// crosschain 95 usdc to fuji
		assert.Nil(t, err)
	})
	t.Run("no swap + stable + no crosschain", func(t *testing.T) {
		cid, err := startChat()
		assert.Nil(t, err)
		err = initBalance(cid)
		assert.Nil(t, err)
		err = chat(cid, "I want to transfer 120 dollars to 0x5134F00C95b8e794db38E1eE39397d8086cee7Ed on target chain fuji")
		assert.Nil(t, err)
	})

	// todo test
	t.Run("no swap + stable + no crosschain", func(t *testing.T) {
		cid, err := startChat()
		assert.Nil(t, err)
		err = initBalance(cid)
		assert.Nil(t, err)
		err = chat(cid, "I want to transfer 120 dollars to 0x5134F00C95b8e794db38E1eE39397d8086cee7Ed on target chain fuji")
		assert.Nil(t, err)
	})
}

func initBalance(cid string) error {
	client := http.Client{}
	// balance:
	// -- mumbai: 100USDC 80USDt 170DAI
	// -- fuji: 25USDC 60USDt 90DAI
	req := &model.CtxRequest{
		Address:   "0x5134F00C95b8e794db38E1eE39397d8086cee7Ed",
		BaseChain: "mumbai",
		Balances: map[string][]model.Reserve{
			"mumbai": {{Symbol: "USDC", Balance: 100}, {Symbol: "USDT", Balance: 80}, {Symbol: "DAI", Balance: 170}},
			"fuji":   {{Symbol: "USDC", Balance: 25}, {Symbol: "USDT", Balance: 60}, {Symbol: "DAI", Balance: 90}},
		},
	}
	buf, err := json.Marshal(req)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest("POST", "http://127.0.0.1:8080/v1/ctx", bytes.NewReader(buf))
	if err != nil {
		return err
	}
	httpReq.Header.Set(model.CIDHeader, cid)
	_, err = client.Do(httpReq)
	if err != nil {
		return err
	}
	fmt.Printf("init ctx\n")
	return nil
}

func startChat() (string, error) {
	client := http.Client{}
	req := &model.DemandRequest{}
	buf, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	httpReq, err := http.NewRequest("POST", "http://127.0.0.1:8080/v1/chat", bytes.NewReader(buf))
	if err != nil {
		return "", err
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	cid := resp.Header.Get(model.CIDHeader)
	fmt.Printf("start cid: %v\n", cid)
	return cid, nil
}

func chat(cid, demand string) error {
	client := http.Client{}
	req := &model.DemandRequest{
		Model:  "v1",
		Demand: demand,
	}
	buf, err := json.Marshal(req)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest("POST", "http://127.0.0.1:8080/v1/chat", bytes.NewReader(buf))
	if err != nil {
		return err
	}
	httpReq.Header.Set(model.CIDHeader, cid)
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	res := &model.DemandResponse{}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, &res); err != nil {
		return err
	}
	fmt.Printf("resp: %+v\n", res)
	return nil
}
