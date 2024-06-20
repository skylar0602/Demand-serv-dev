package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	log "github.com/cihub/seelog"
	"github.com/smarterwallet/demand-abstraction-serv/config"
	"github.com/smarterwallet/demand-abstraction-serv/model"
)

var (
	Base *BaseService
)

type BaseService struct {
	cfg *config.Config
}

func NewBaseService(cfg *config.Config) {
	Base = &BaseService{cfg: cfg}
}

func (s *BaseService) LoadTokens() (*model.AssetConfigResp, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/package", s.cfg.ConfigEndpoint))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("get cross chain config error: %v", err)
		return nil, err
	}
	var res model.AssetConfigResp
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	if res.Code != 200 {
		return nil, errors.New("load tokens error")
	}
	return &res, nil
}

func (s *BaseService) CheckCross(sourceChainId, targetChainId int, token string) (bool, []byte) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/cross-chain-config?sourceChainId=%d&destChainId=%d&crossChainTokenName=%s", s.cfg.CrossEndpoint, sourceChainId, targetChainId, token))
	if err != nil {
		log.Errorf("get cross chain config error: %v", err)
		return false, nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("get cross chain config error: %v", err)
		return false, nil
	}
	var crossChainConfig model.CrossChainResp
	if err := json.Unmarshal(body, &crossChainConfig); err != nil {
		log.Errorf("get cross chain config error: %v", err)
		return false, nil
	}
	return crossChainConfig.Code == 200, body
}

func (s *BaseService) CheckSwap(req model.SwapReq) (string, []byte, error) {
	buf, _ := json.Marshal(req)
	resp, err := http.Post(fmt.Sprintf("%s/api/v1/swap-path/min-in-amount", s.cfg.SwapEndpoint), "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Errorf("uniswap request error: %s", err)
		return "", nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("get cross chain config error: %v", err)
		return "", nil, err
	}
	var res model.SwapResp
	if err = json.Unmarshal(body, &res); err != nil {
		log.Errorf("unmarshal swap result error: %v", err)
		return "", nil, err
	}
	if res.Code != 200 || res.Result.MinInAmount == "" {
		log.Warnf("uniswap unable to swap: %d", res.Code)
		return "", nil, err
	}
	return res.Result.MinInAmount, body, nil
}
