package model

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const (
	ChainInternalTransfer = "chain-internal-transfer"
	CrossChainTransfer    = "cross-chain-transfer"
)

type (
	Reserve struct {
		Symbol  string  `json:"symbol"`
		Balance float64 `json:"balance"`
		Address string  `json:"address"`
	}
	CtxRequest struct {
		Address   string               `json:"address"`
		BaseChain string               `json:"baseChain"`
		Balances  map[string][]Reserve `json:"balances"`
	}
	DemandRequest struct {
		Model  string `json:"model"`
		Demand string `json:"demand"`
	}
	DemandResponse struct {
		Category string     `json:"category"`
		Summary  string     `json:"summary"`
		Detail   DetailResp `json:"detail"`
	}
	DetailResp struct {
		Reply string        `json:"reply"`
		OPs   []interface{} `json:"ops"`
	}
	CrossChainResponse struct {
		Type            string          `json:"type"`
		RawResponse     json.RawMessage `json:"raw_response"`
		SourceChainId   int             `json:"source_chain_id"`
		SourceChainName string          `json:"source_chain_name"`
		Token           string          `json:"token"`
		Amount          string          `json:"amount"`
		Receiver        string          `json:"receiver"`
		TargetChainId   int             `json:"target_chain_id"`
		TargetChainName string          `json:"target_chain_name"`
	}
	SwapResponse struct {
		Type        string          `json:"type"`
		RawResponse json.RawMessage `json:"raw_response"`
		ChainId     int             `json:"chain_id"`
		ChainName   string          `json:"chain_name"`
		SourceToken string          `json:"source_token"`
		TargetToken string          `json:"target_token"`
		Dex         string          `json:"dex"`
		SwapIn      string          `json:"swap_in"`
		SwapOut     string          `json:"swap_out"`
	}
	TradeStrategyResponse struct {
		BotName    string      `json:"bot_name"`
		Strategy   string      `json:"strategy"`
		MinReturn  string      `json:"min_return"`
		MaxReturn  string      `json:"max_return"`
		Operations []Operation `json:"operations"`
	}
	Operation struct {
		Seq   int            `json:"seq"`
		Type  string         `json:"type"`
		Param OperationParam `json:"param"`
	}
	OperationParam struct {
		From             string      `json:"from"`
		To               string      `json:"to"`
		GasFee           string      `json:"gas_fee"`
		FeeUint          string      `json:"fee_uint"`
		ConditionsSymbol []string    `json:"conditions_symbol"`
		Conditions       []Condition `json:"conditions"`
	}
	Condition struct {
		TokenName  string `json:"tokenName"`
		Trend      string `json:"trend"`
		Percentage string `json:"percentage"`
	}
)

const (
	ConversationID   = "conversationID"
	CIDHeader        = "X-SmartWallet-CID"
	DialogueRoleUser = "User"
	DialogueRoleAI   = "AI"
	ModelV1          = "v1"
)

type (
	ConversationCtx struct {
		Cid       string     `json:"cid"`
		Dialogues []Dialogue `json:"dialogues"`
	}
	Dialogue struct {
		Type      string `json:"type"`
		Role      string `json:"role"`
		Content   string `json:"content"`
		Timestamp int64  `json:"timestamp"`
	}
)

type CrossChainResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  []struct {
		ID                  int             `json:"ID"`
		CreatedAt           time.Time       `json:"CreatedAt"`
		UpdatedAt           time.Time       `json:"UpdatedAt"`
		DeletedAt           interface{}     `json:"DeletedAt"`
		SourceChainId       int             `json:"sourceChainId"`
		DestChainId         int             `json:"destChainId"`
		Priority            int             `json:"priority"`
		CrossChainTokenName string          `json:"crossChainTokenName"`
		ProtocolName        string          `json:"protocolName"`
		Config              json.RawMessage `json:"config"`
	} `json:"result"`
}

type SwapResult struct {
	Quote struct {
		Numerator   []interface{} `json:"numerator"`
		Denominator []int         `json:"denominator"`
		Currency    struct {
			ChainId  int    `json:"chainId"`
			Decimals int    `json:"decimals"`
			IsNative bool   `json:"isNative"`
			IsToken  bool   `json:"isToken"`
			Address  string `json:"address"`
		} `json:"currency"`
		DecimalScale []int `json:"decimalScale"`
	} `json:"quote"`
	QuoteGasAdjusted struct {
		Numerator   []int `json:"numerator"`
		Denominator []int `json:"denominator"`
		Currency    struct {
			ChainId  int    `json:"chainId"`
			Decimals int    `json:"decimals"`
			IsNative bool   `json:"isNative"`
			IsToken  bool   `json:"isToken"`
			Address  string `json:"address"`
		} `json:"currency"`
		DecimalScale []int `json:"decimalScale"`
	} `json:"quoteGasAdjusted"`
	EstimatedGasUsed struct {
		Type string `json:"type"`
		Hex  string `json:"hex"`
	} `json:"estimatedGasUsed"`
	EstimatedGasUsedQuoteToken struct {
		Numerator   []int `json:"numerator"`
		Denominator []int `json:"denominator"`
		Currency    struct {
			ChainId  int    `json:"chainId"`
			Decimals int    `json:"decimals"`
			Symbol   string `json:"symbol"`
			Name     string `json:"name"`
			IsNative bool   `json:"isNative"`
			IsToken  bool   `json:"isToken"`
			Address  string `json:"address"`
		} `json:"currency"`
		DecimalScale []int `json:"decimalScale"`
	} `json:"estimatedGasUsedQuoteToken"`
	EstimatedGasUsedUSD struct {
		Numerator   []int `json:"numerator"`
		Denominator []int `json:"denominator"`
		Currency    struct {
			ChainId  int    `json:"chainId"`
			Decimals int    `json:"decimals"`
			Symbol   string `json:"symbol"`
			Name     string `json:"name"`
			IsNative bool   `json:"isNative"`
			IsToken  bool   `json:"isToken"`
			Address  string `json:"address"`
		} `json:"currency"`
		DecimalScale []int `json:"decimalScale"`
	} `json:"estimatedGasUsedUSD"`
	GasPriceWei struct {
		Type string `json:"type"`
		Hex  string `json:"hex"`
	} `json:"gasPriceWei"`
	Route []struct {
		Protocol string `json:"protocol"`
		Amount   struct {
			Numerator   []int `json:"numerator"`
			Denominator []int `json:"denominator"`
			Currency    struct {
				ChainId  int    `json:"chainId"`
				Decimals int    `json:"decimals"`
				IsNative bool   `json:"isNative"`
				IsToken  bool   `json:"isToken"`
				Address  string `json:"address"`
			} `json:"currency"`
			DecimalScale []int `json:"decimalScale"`
		} `json:"amount"`
		RawQuote struct {
			Type string `json:"type"`
			Hex  string `json:"hex"`
		} `json:"rawQuote"`
		SqrtPriceX96AfterList []struct {
			Type string `json:"type"`
			Hex  string `json:"hex"`
		} `json:"sqrtPriceX96AfterList"`
		InitializedTicksCrossedList []int `json:"initializedTicksCrossedList"`
		QuoterGasEstimate           struct {
			Type string `json:"type"`
			Hex  string `json:"hex"`
		} `json:"quoterGasEstimate"`
		Quote struct {
			Numerator   []interface{} `json:"numerator"`
			Denominator []int         `json:"denominator"`
			Currency    struct {
				ChainId  int    `json:"chainId"`
				Decimals int    `json:"decimals"`
				IsNative bool   `json:"isNative"`
				IsToken  bool   `json:"isToken"`
				Address  string `json:"address"`
			} `json:"currency"`
			DecimalScale []int `json:"decimalScale"`
		} `json:"quote"`
		Percent int `json:"percent"`
		Route   struct {
			MidPrice interface{} `json:"_midPrice"`
			Pools    []struct {
				Token0 struct {
					ChainId  int    `json:"chainId"`
					Decimals int    `json:"decimals"`
					IsNative bool   `json:"isNative"`
					IsToken  bool   `json:"isToken"`
					Address  string `json:"address"`
				} `json:"token0"`
				Token1 struct {
					ChainId  int    `json:"chainId"`
					Decimals int    `json:"decimals"`
					Symbol   string `json:"symbol"`
					Name     string `json:"name"`
					IsNative bool   `json:"isNative"`
					IsToken  bool   `json:"isToken"`
					Address  string `json:"address"`
				} `json:"token1"`
				Fee              int   `json:"fee"`
				SqrtRatioX96     []int `json:"sqrtRatioX96"`
				Liquidity        []int `json:"liquidity"`
				TickCurrent      int   `json:"tickCurrent"`
				TickDataProvider struct {
				} `json:"tickDataProvider"`
			} `json:"pools"`
			TokenPath []struct {
				ChainId  int    `json:"chainId"`
				Decimals int    `json:"decimals"`
				IsNative bool   `json:"isNative"`
				IsToken  bool   `json:"isToken"`
				Address  string `json:"address"`
				Symbol   string `json:"symbol,omitempty"`
				Name     string `json:"name,omitempty"`
			} `json:"tokenPath"`
			Input struct {
				ChainId  int    `json:"chainId"`
				Decimals int    `json:"decimals"`
				IsNative bool   `json:"isNative"`
				IsToken  bool   `json:"isToken"`
				Address  string `json:"address"`
			} `json:"input"`
			Output struct {
				ChainId  int    `json:"chainId"`
				Decimals int    `json:"decimals"`
				IsNative bool   `json:"isNative"`
				IsToken  bool   `json:"isToken"`
				Address  string `json:"address"`
			} `json:"output"`
			Protocol string `json:"protocol"`
		} `json:"route"`
		GasModel struct {
		} `json:"gasModel"`
		QuoteToken struct {
			ChainId  int    `json:"chainId"`
			Decimals int    `json:"decimals"`
			IsNative bool   `json:"isNative"`
			IsToken  bool   `json:"isToken"`
			Address  string `json:"address"`
		} `json:"quoteToken"`
		TradeType      int `json:"tradeType"`
		GasCostInToken struct {
			Numerator   []int `json:"numerator"`
			Denominator []int `json:"denominator"`
			Currency    struct {
				ChainId  int    `json:"chainId"`
				Decimals int    `json:"decimals"`
				Symbol   string `json:"symbol"`
				Name     string `json:"name"`
				IsNative bool   `json:"isNative"`
				IsToken  bool   `json:"isToken"`
				Address  string `json:"address"`
			} `json:"currency"`
			DecimalScale []int `json:"decimalScale"`
		} `json:"gasCostInToken"`
		GasCostInUSD struct {
			Numerator   []int `json:"numerator"`
			Denominator []int `json:"denominator"`
			Currency    struct {
				ChainId  int    `json:"chainId"`
				Decimals int    `json:"decimals"`
				Symbol   string `json:"symbol"`
				Name     string `json:"name"`
				IsNative bool   `json:"isNative"`
				IsToken  bool   `json:"isToken"`
				Address  string `json:"address"`
			} `json:"currency"`
			DecimalScale []int `json:"decimalScale"`
		} `json:"gasCostInUSD"`
		GasEstimate struct {
			Type string `json:"type"`
			Hex  string `json:"hex"`
		} `json:"gasEstimate"`
		QuoteAdjustedForGas struct {
			Numerator   []int `json:"numerator"`
			Denominator []int `json:"denominator"`
			Currency    struct {
				ChainId  int    `json:"chainId"`
				Decimals int    `json:"decimals"`
				IsNative bool   `json:"isNative"`
				IsToken  bool   `json:"isToken"`
				Address  string `json:"address"`
			} `json:"currency"`
			DecimalScale []int `json:"decimalScale"`
		} `json:"quoteAdjustedForGas"`
		PoolAddresses []string `json:"poolAddresses"`
		TokenPath     []struct {
			ChainId  int    `json:"chainId"`
			Decimals int    `json:"decimals"`
			IsNative bool   `json:"isNative"`
			IsToken  bool   `json:"isToken"`
			Address  string `json:"address"`
			Symbol   string `json:"symbol,omitempty"`
			Name     string `json:"name,omitempty"`
		} `json:"tokenPath"`
	} `json:"route"`
	Trade struct {
		Swaps []struct {
			Route struct {
				MidPrice interface{} `json:"_midPrice"`
				Pools    []struct {
					Token0 struct {
						ChainId  int    `json:"chainId"`
						Decimals int    `json:"decimals"`
						IsNative bool   `json:"isNative"`
						IsToken  bool   `json:"isToken"`
						Address  string `json:"address"`
					} `json:"token0"`
					Token1 struct {
						ChainId  int    `json:"chainId"`
						Decimals int    `json:"decimals"`
						Symbol   string `json:"symbol"`
						Name     string `json:"name"`
						IsNative bool   `json:"isNative"`
						IsToken  bool   `json:"isToken"`
						Address  string `json:"address"`
					} `json:"token1"`
					Fee              int   `json:"fee"`
					SqrtRatioX96     []int `json:"sqrtRatioX96"`
					Liquidity        []int `json:"liquidity"`
					TickCurrent      int   `json:"tickCurrent"`
					TickDataProvider struct {
					} `json:"tickDataProvider"`
				} `json:"pools"`
				TokenPath []struct {
					ChainId  int    `json:"chainId"`
					Decimals int    `json:"decimals"`
					IsNative bool   `json:"isNative"`
					IsToken  bool   `json:"isToken"`
					Address  string `json:"address"`
					Symbol   string `json:"symbol,omitempty"`
					Name     string `json:"name,omitempty"`
				} `json:"tokenPath"`
				Input struct {
					ChainId  int    `json:"chainId"`
					Decimals int    `json:"decimals"`
					IsNative bool   `json:"isNative"`
					IsToken  bool   `json:"isToken"`
					Address  string `json:"address"`
				} `json:"input"`
				Output struct {
					ChainId  int    `json:"chainId"`
					Decimals int    `json:"decimals"`
					IsNative bool   `json:"isNative"`
					IsToken  bool   `json:"isToken"`
					Address  string `json:"address"`
				} `json:"output"`
				Protocol string `json:"protocol"`
				Path     []struct {
					ChainId  int    `json:"chainId"`
					Decimals int    `json:"decimals"`
					IsNative bool   `json:"isNative"`
					IsToken  bool   `json:"isToken"`
					Address  string `json:"address"`
					Symbol   string `json:"symbol,omitempty"`
					Name     string `json:"name,omitempty"`
				} `json:"path"`
			} `json:"route"`
			InputAmount struct {
				Numerator   []int `json:"numerator"`
				Denominator []int `json:"denominator"`
				Currency    struct {
					ChainId  int    `json:"chainId"`
					Decimals int    `json:"decimals"`
					IsNative bool   `json:"isNative"`
					IsToken  bool   `json:"isToken"`
					Address  string `json:"address"`
				} `json:"currency"`
				DecimalScale []int `json:"decimalScale"`
			} `json:"inputAmount"`
			OutputAmount struct {
				Numerator   []interface{} `json:"numerator"`
				Denominator []int         `json:"denominator"`
				Currency    struct {
					ChainId  int    `json:"chainId"`
					Decimals int    `json:"decimals"`
					IsNative bool   `json:"isNative"`
					IsToken  bool   `json:"isToken"`
					Address  string `json:"address"`
				} `json:"currency"`
				DecimalScale []int `json:"decimalScale"`
			} `json:"outputAmount"`
		} `json:"swaps"`
		Routes []struct {
			MidPrice interface{} `json:"_midPrice"`
			Pools    []struct {
				Token0 struct {
					ChainId  int    `json:"chainId"`
					Decimals int    `json:"decimals"`
					IsNative bool   `json:"isNative"`
					IsToken  bool   `json:"isToken"`
					Address  string `json:"address"`
				} `json:"token0"`
				Token1 struct {
					ChainId  int    `json:"chainId"`
					Decimals int    `json:"decimals"`
					Symbol   string `json:"symbol"`
					Name     string `json:"name"`
					IsNative bool   `json:"isNative"`
					IsToken  bool   `json:"isToken"`
					Address  string `json:"address"`
				} `json:"token1"`
				Fee              int   `json:"fee"`
				SqrtRatioX96     []int `json:"sqrtRatioX96"`
				Liquidity        []int `json:"liquidity"`
				TickCurrent      int   `json:"tickCurrent"`
				TickDataProvider struct {
				} `json:"tickDataProvider"`
			} `json:"pools"`
			TokenPath []struct {
				ChainId  int    `json:"chainId"`
				Decimals int    `json:"decimals"`
				IsNative bool   `json:"isNative"`
				IsToken  bool   `json:"isToken"`
				Address  string `json:"address"`
				Symbol   string `json:"symbol,omitempty"`
				Name     string `json:"name,omitempty"`
			} `json:"tokenPath"`
			Input struct {
				ChainId  int    `json:"chainId"`
				Decimals int    `json:"decimals"`
				IsNative bool   `json:"isNative"`
				IsToken  bool   `json:"isToken"`
				Address  string `json:"address"`
			} `json:"input"`
			Output struct {
				ChainId  int    `json:"chainId"`
				Decimals int    `json:"decimals"`
				IsNative bool   `json:"isNative"`
				IsToken  bool   `json:"isToken"`
				Address  string `json:"address"`
			} `json:"output"`
			Protocol string `json:"protocol"`
			Path     []struct {
				ChainId  int    `json:"chainId"`
				Decimals int    `json:"decimals"`
				IsNative bool   `json:"isNative"`
				IsToken  bool   `json:"isToken"`
				Address  string `json:"address"`
				Symbol   string `json:"symbol,omitempty"`
				Name     string `json:"name,omitempty"`
			} `json:"path"`
		} `json:"routes"`
		TradeType int `json:"tradeType"`
	} `json:"trade"`
	MethodParameters struct {
		Calldata string `json:"calldata"`
		Value    string `json:"value"`
		To       string `json:"to"`
	} `json:"methodParameters"`
	BlockNumber struct {
		Type string `json:"type"`
		Hex  string `json:"hex"`
	} `json:"blockNumber"`
	HitsCachedRoute bool `json:"hitsCachedRoute"`
}

type SwapReq struct {
	ChainId         int     `json:"chainId"`
	TokenInAddress  string  `json:"tokenIn"`
	TokenOutAddress string  `json:"tokenOut"`
	AmountOut       float64 `json:"amountOut"`
}

type SwapResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		MinInAmount string `json:"minInAmount"`
	} `json:"result"`
}

type AssetConfigResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Common struct {
			ID        int         `json:"ID"`
			CreatedAt time.Time   `json:"CreatedAt"`
			UpdatedAt time.Time   `json:"UpdatedAt"`
			DeletedAt interface{} `json:"DeletedAt"`
			Name      string      `json:"name"`
			Version   string      `json:"version"`
			Config    struct {
				Url struct {
					Mpc struct {
						Api  string `json:"api"`
						Wasm string `json:"wasm"`
					} `json:"mpc"`
					AutoTrading struct {
						Mumbai string `json:"mumbai"`
					} `json:"autoTrading"`
					Asset   string `json:"asset"`
					Storage string `json:"storage"`
				} `json:"url"`
				ContractAddress struct {
					AutoTrading string `json:"autoTrading"`
				} `json:"contractAddress"`
			} `json:"config"`
		} `json:"common"`
		Chain []struct {
			ID        int         `json:"ID"`
			CreatedAt time.Time   `json:"CreatedAt"`
			UpdatedAt time.Time   `json:"UpdatedAt"`
			DeletedAt interface{} `json:"DeletedAt"`
			NetWorkId int         `json:"netWorkId"`
			Name      string      `json:"name"`
			Icon      string      `json:"icon"`
			Tokens    []struct {
				TokenId int    `json:"tokenId"`
				Name    string `json:"name"`
				Fee     int    `json:"fee,omitempty"`
				Address string `json:"address"`
				Decimal int    `json:"decimal"`
				Icon    string `json:"icon"`
				Type    int    `json:"type,omitempty"`
			} `json:"tokens"`
			Erc4337ContractAddress *struct {
				SimpleAccountFactory string `json:"simpleAccountFactory"`
				TokenPaymaster       struct {
					Swt string `json:"swt"`
				} `json:"tokenPaymaster"`
				Entrypoint string `json:"entrypoint"`
			} `json:"erc4337ContractAddress"`
			RpcApi          string `json:"rpcApi,omitempty"`
			BundlerApi      string `json:"bundlerApi,omitempty"`
			BlockScanUrl    string `json:"blockScanUrl,omitempty"`
			CreateWalletApi string `json:"createWalletApi,omitempty"`
			ApiType         int    `json:"apiType,omitempty"`
			ProduceBlock24H int    `json:"produceBlock24h,omitempty"`
		} `json:"chain"`
	} `json:"result"`
}

func (d *Dialogue) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (c *CtxRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CtxRequest) GetTokenAddress(chain, symbol string) string {
	for _, b := range c.Balances[chain] {
		if b.Symbol == symbol {
			return b.Address
		}
	}
	return ""
}

func (c *CtxRequest) GetTokenBalance(chain, symbol string) decimal.Decimal {
	for _, b := range c.Balances[chain] {
		if b.Symbol == symbol {
			return decimal.NewFromFloat(b.Balance)
		}
	}
	return decimal.Zero
}

func (c *CtxRequest) Format() {
	c.BaseChain = strings.ToLower(c.BaseChain)
	b := make(map[string][]Reserve)
	for chain, reserves := range c.Balances {
		lc := strings.ToLower(chain)
		b[lc] = make([]Reserve, 0)
		for _, reserve := range reserves {
			b[lc] = append(b[lc], Reserve{
				Symbol:  strings.ToUpper(reserve.Symbol),
				Balance: reserve.Balance,
			})
		}
	}
	c.Balances = b
}

func (c *CtxRequest) GetTokens(chain string) []Reserve {
	tokens := make([]Reserve, 0)
	for _, reserve := range c.Balances[chain] {
		tokens = append(tokens, reserve)
	}
	return tokens
}
