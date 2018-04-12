package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/seanmclane/goatnickels/rpc"
	"math/big"
	"strconv"
	"time"
)

type Vote struct {
	Account   string    `json:"account"`
	Hash      string    `json:"hash"`
	Signature Signature `json:"signature"`
}

//create vote set to store votes from network
var VoteSet []Vote

//create bool to capture whether voting
var Voting bool

func CheckConsensus() {
	//TODO: final criteria for consensus = 2/3 of stakes sign hash of candidate set transaction

	voteCount := make(map[string]int)
	total := 0

	//TODO: get proportions for all, this won't work if this node is the one out of sync
	//changed the consensus to check votes in voteset
	//if agreed-upon hash is same as local hash, apply transactions
	//TODO: if not, wait and request new block

	for _, v := range VoteSet {
		voteCount[v.Hash] += 1
		total += 1
	}

	winningVote := 0
	var winningHash string
	for key, vote := range voteCount {
		if vote > winningVote {
			winningVote = vote
			winningHash = key
		}
	}

	if total < 1 {
		//catch up to network
		time.Sleep(3 * time.Second)
		maxBlock, nodes := GetMaxBlockNumberFromNetwork()
		localMax := FindMaxBlock()
		if localMax < maxBlock {
			GetBlockChainFromNetwork(localMax, maxBlock, nodes[0])
			LastGoatBlock = ReadBlockFromLocalStorage(strconv.Itoa(maxBlock))
			//TODO: loop through to get real max block
		} else {
			fmt.Println("No consensus reached due to lack of votes")
		}
		return
	}

	if winningVote/total >= 2/3 {
		csHash := HashCandidateSet(&CandidateSet)
		vHash, err := hex.DecodeString(winningHash)
		if err != nil {
			fmt.Println("error:", err)
		}

		if bytes.Equal(vHash, csHash) {
			//reward voters on winning hash (add transactions to candidate set)
			RewardVoters(winningHash)
			//build next block
			//TODO: consider whether something needs to be refactored out to make this simpler
			NextBlock()
		} else {
			time.Sleep(3 * time.Second)
			localMax := FindMaxBlock()
			maxBlock, nodes := GetMaxBlockNumberFromNetwork()
			GetBlockChainFromNetwork(localMax, maxBlock, nodes[0])
			LastGoatBlock = ReadBlockFromLocalStorage(strconv.Itoa(maxBlock))
		}
		//TODO: loop through to get real max block
	} else {
		//TODO: tell network to restart consensus round?
		fmt.Println("No consensus reached")
	}
}

func ResetVoteSet() {
	VoteSet = []Vote{}
}

func RewardVoters(hash string) {
	for _, v := range VoteSet {
		if v.Hash == hash {
			CandidateSet = append(CandidateSet, Transaction{
				From:   "mine",
				To:     v.Account,
				Amount: 100000000, //TODO: scale reward to stake
			})
		} else {
			CandidateSet = append(CandidateSet, Transaction{
				From:   "mine",
				To:     v.Account,
				Amount: -100000000, //TODO: scale reward to stake
			})
		}
	}
}

func (v *Vote) AddVote() (ok bool) {
	ok = v.VerifyVote()
	if ok != true {
		return false
	}

	VoteSet = append(VoteSet, *v)
	return ok
}

func (v *Vote) VerifyVote() (ok bool) {

	//check that the vote is not in the voteset already
	for _, vote := range VoteSet {
		if v.Account == vote.Account && v.Hash == vote.Hash {
			return false
		}
		//allow votes with the same account but different hashes to be added
	}

	//verify signature of account sending the vote
	//TODO: abstract key recreation into a function (hash, r, s) (ok bool)
	hash := v.HashVote()
	//check that key is well formed
	if len(v.Account) < 100 {
		return false
	}
	//remove goat_ from key
	publicKey := v.Account[5:]
	//recreate ecdsa.PublicKey from pub
	byteKey, err := hex.DecodeString(publicKey)
	if err != nil {
		fmt.Println("error:", err)
	}
	x, y := elliptic.Unmarshal(elliptic.P384(), byteKey)
	if x == nil || y == nil {
		return false
	}
	pub := new(ecdsa.PublicKey)
	pub.Curve = elliptic.P384()
	pub.X, pub.Y = x, y
	//convert r and s back to big ints
	byteR, err := hex.DecodeString(v.Signature.R)
	if err != nil {
		fmt.Println("error:", err)
	}
	r := new(big.Int).SetBytes(byteR)
	byteS, err := hex.DecodeString(v.Signature.S)
	if err != nil {
		fmt.Println("error:", err)
	}
	s := new(big.Int).SetBytes(byteS)
	//verify signature
	sigOk := ecdsa.Verify(pub, hash, r, s)

	if sigOk {
		fmt.Println("vote verified")
		return true
	} else {
		return false
	}
}

func (v *Vote) SignVote() (r string, s string) {
	keystore := LoadKeyStore()
	privateKey := keystore.PrivateKey

	hash := v.HashVote()
	//recreate ecdsa.PrivateKey from private_key
	byteKey, err := hex.DecodeString(privateKey)
	if err != nil {
		fmt.Println("error:", err)
	}
	bigintKey := new(big.Int).SetBytes(byteKey)
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P384()
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(byteKey)
	priv.D = bigintKey
	rInt, sInt, err := ecdsa.Sign(rand.Reader, priv, hash)
	if err != nil {
		fmt.Println("error:", err)
	}
	r = hex.EncodeToString(rInt.Bytes())
	s = hex.EncodeToString(sInt.Bytes())
	return r, s
}

func (v *Vote) HashVote() (h []byte) {
	hashString := v.Account + v.Hash
	fixedHash := sha512.Sum512([]byte(hashString))
	h = fixedHash[:]
	return h
}

func SendVoteToNetwork() {

	config := LoadConfig()

	v := Vote{
		Account: config.Account,
		Hash:    hex.EncodeToString(HashCandidateSet(&CandidateSet)),
	}

	r, s := v.SignVote()

	v.Signature = Signature{
		R: r,
		S: s,
	}

	//ensure the vote exists locally in the voteset, then broadcast
	v.AddVote()

	out, err := json.Marshal(v)
	if err != nil {
		fmt.Println("error:", err)
	}

	//broadcast vote to network through websocket
	hub.Broadcast <- rpc.BuildNotification("vote", out)
}
