package redt

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethertypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	lrucache "github.com/hashicorp/golang-lru"
)

func Test_sigHash(t *testing.T) {
	type args struct {
		header *ethertypes.Header
	}
	tests := []struct {
		name     string
		args     args
		wantHash common.Hash
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHash := sigHash(tt.args.header); !reflect.DeepEqual(gotHash, tt.wantHash) {
				t.Errorf("sigHash() = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}

func Test_toBlockNumArg(t *testing.T) {
	type args struct {
		number int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toBlockNumArg(tt.args.number); got != tt.want {
				t.Errorf("toBlockNumArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestUrl() string {
	urlB, err := os.ReadFile("privatedata/test_url_http")
	if err != nil {
		return ""
	}
	return string(urlB)
}

func TestNewRedTNode(t *testing.T) {

	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "firstTest",
			args: args{
				url: getTestUrl(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRedTNode(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRedTNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRedTNode_Close(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			rt.Close()
		})
	}
}

func TestRedTNode_EthClient(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   *ethclient.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			if got := rt.EthClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedTNode.EthClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedTNode_RpcClient(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   *rpc.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			if got := rt.RpcClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedTNode.RpcClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedTNode_BlockByHash(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	type args struct {
		blockhash common.Hash
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ethertypes.Block
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			got, err := rt.BlockByHash(tt.args.blockhash)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedTNode.BlockByHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedTNode.BlockByHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedTNode_HeaderByNumber(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	type args struct {
		number int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ethertypes.Header
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			got, err := rt.HeaderByNumber(tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedTNode.HeaderByNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedTNode.HeaderByNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedTNode_CurrentBlockNumber(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			got, err := rt.CurrentBlockNumber()
			if (err != nil) != tt.wantErr {
				t.Errorf("RedTNode.CurrentBlockNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedTNode.CurrentBlockNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedTNode_NodeInfo(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		want    *p2p.NodeInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			got, err := rt.NodeInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("RedTNode.NodeInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedTNode.NodeInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedTNode_Peers(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*p2p.PeerInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			got, err := rt.Peers()
			if (err != nil) != tt.wantErr {
				t.Errorf("RedTNode.Peers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedTNode.Peers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedTNode_RefreshValidators(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			if err := rt.RefreshValidators(); (err != nil) != tt.wantErr {
				t.Errorf("RedTNode.RefreshValidators() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedTNode_Validators(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   []common.Address
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			if got := rt.Validators(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedTNode.Validators() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedTNode_AllValidators(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   map[common.Address]*NodeInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			if got := rt.AllValidators(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedTNode.AllValidators() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedTNode_ValidatorInfo(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	type args struct {
		validator common.Address
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *NodeInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			if got := rt.ValidatorInfo(tt.args.validator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedTNode.ValidatorInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedTNode_DisplayMyInfo(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			rt.DisplayMyInfo()
		})
	}
}

func TestRedTNode_DisplayPeersInfo(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			rt.DisplayPeersInfo()
		})
	}
}

func TestRedTNode_getValSet(t *testing.T) {
	type fields struct {
		cli           *ethclient.Client
		headerCache   *lrucache.Cache
		rpccli        *rpc.Client
		valSet        []common.Address
		allValidators map[common.Address]*NodeInfo
		timeout       time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		want    []common.Address
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RedTNode{
				cli:           tt.fields.cli,
				headerCache:   tt.fields.headerCache,
				rpccli:        tt.fields.rpccli,
				valSet:        tt.fields.valSet,
				allValidators: tt.fields.allValidators,
				timeout:       tt.fields.timeout,
			}
			got, err := rt.getValSet()
			if (err != nil) != tt.wantErr {
				t.Errorf("RedTNode.getValSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedTNode.getValSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSignersFromBlock(t *testing.T) {
	type args struct {
		header *ethertypes.Header
	}
	tests := []struct {
		name        string
		args        args
		wantAuthor  common.Address
		wantSigners []common.Address
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAuthor, gotSigners, err := SignersFromBlock(tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignersFromBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAuthor, tt.wantAuthor) {
				t.Errorf("SignersFromBlock() gotAuthor = %v, want %v", gotAuthor, tt.wantAuthor)
			}
			if !reflect.DeepEqual(gotSigners, tt.wantSigners) {
				t.Errorf("SignersFromBlock() gotSigners = %v, want %v", gotSigners, tt.wantSigners)
			}
		})
	}
}

func TestDisplayPeersInfo(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DisplayPeersInfo(tt.args.url)
		})
	}
}

func Test_readNodeListFromGithub(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    []*NodeInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readNodeListFromGithub(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("readNodeListFromGithub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readNodeListFromGithub() = %v, want %v", got, tt.want)
			}
		})
	}
}
