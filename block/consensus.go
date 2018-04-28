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
	Index     int       `json:"index"`
	Hash      string    `json:"hash"`
	Signature Signature `json:"signature"`
}

//create channel to avoid write conflicts to voteset
var VoteChannel chan rpc.JsonRpcMessage

//create vote set to store votes from network
var VoteSet []Vote

//create bool to capture whether voting
var Voting bool

func runVoting() {
	//initialize voting to whether a voting round on the candidate set is happening
	Voting = false

	//TODO: wait for state to sync but not in a way that hangs...
	//after ten failures request maxblock through api?

	/* for state.maxBlock < 1 {
		fmt.Println("waiting for maxblock to sync")
		time.Sleep(time.Second * 1)
	} */

	for {
		voteStartList := [6]int{8, 18, 28, 38, 48, 58}
		voteEndList := [6]int{0, 10, 20, 30, 40, 50}

		for _, sec := range voteStartList {
			if time.Now().Second() == sec {
				//TODO: Move any transactions not reaching 80 percent to the staging set
				Voting = true
				fmt.Println("---------- Start voting round ---------")
				SendVoteToNetwork()
			}
		}
		for _, sec := range voteEndList {
			if time.Now().Second() == sec {
				CheckConsensus()
				Voting = false
				//add staging transactions back into candidate set
				//for now I will overwrite the candidate set with the staging set
				//TODO: ensure transactions that were not in the applied candidate set stay in the new candidate set with all the staging transactions
				ResetCandidateSet()
				ResetVoteSet()
				fmt.Println("---------- End voting round ---------")

				//temporarily generate transactions for testing
				txion := Transaction{
					To:       "goat_04ab1594a3b65e440653b1a54952aee3cb7f5c41cb476f7ecd3ce58dc23cef0923beb45fc275ff4149cd9f0417f8ca885e882b3b68d00bab2988b22f2eaf7f6683ba3e672abd668e5788a8ecb4d055cd024f004ff03db06158f18e5bd02914685a",
					From:     "goat_04dbb67ae9650ca3258071909f74be5400fe53fc2e5dcc82103020f3aeefeee5f9980c4c05bb8696215458dfa7ddaa1505d2826cab3d246b8930b0694f766a22f8bb63932368c0b12bf80cfaee8a18db1d7ce19df0a84215d20b0bbfbd30d95c25",
					Amount:   LastGoatBlock.Index * 1000,
					Sequence: LastGoatBlock.Index - 1,
				}
				r, s := txion.SignTransaction("131952e67d981417fe5acf0eaca5d576e701fe1f523c647bbb3c387d8234cd94c55f3f457de0dff2d19adb66b654fbc7")
				txion.R = r
				txion.S = s
				out, _ := json.Marshal(txion)
				hub.Broadcast <- rpc.BuildNotification("transaction", out)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func CheckConsensus() {
	//TODO: final criteria for consensus = 2/3 of stakes sign hash of candidate set transaction

	voteCount := make(map[string]int)
	total := 0

	//TODO: get proportions for all, this won't work if this node is the one out of sync
	//changed the consensus to check votes in voteset
	//if agreed-upon hash is same as local hash, apply transactions
	//TODO: if not, wait and request new block

	fmt.Println("VoteSet")

	for _, v := range VoteSet {
		voteCount[v.Hash] += 1
		total += 1
		fmt.Println("Account:", v.Account, "Index:", v.Index)
	}

	winningVote := 0
	var winningHash string
	for key, vote := range voteCount {
		if vote > winningVote {
			winningVote = vote
			winningHash = key
		}
	}

	//TODO: replace maxBlock with state.maxBlock but wait until the state is accurate

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

func handleVoteChannel() {
	VoteChannel = make(chan rpc.JsonRpcMessage)
	for {
		msg := <-VoteChannel
		var vote Vote
		err := json.Unmarshal(msg.Params, &vote)
		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}
		if vote.AddVote() {
			hub.Broadcast <- msg
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
	indexOk := v.Index == LastGoatBlock.Index+1
	if sigOk && indexOk {
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
	hashString := v.Account + v.Hash + strconv.Itoa(v.Index)
	fixedHash := sha512.Sum512([]byte(hashString))
	h = fixedHash[:]
	return h
}

func SendVoteToNetwork() {

	config := LoadConfig()

	v := Vote{
		Account: config.Account,
		Index:   LastGoatBlock.Index + 1,
		Hash:    hex.EncodeToString(HashCandidateSet(&CandidateSet)),
	}

	r, s := v.SignVote()

	v.Signature = Signature{
		R: r,
		S: s,
	}

	out, err := json.Marshal(v)
	if err != nil {
		fmt.Println("error:", err)
	}

	//broadcast vote to network through websocket, do not vote locally to ensure vote propagation
	hub.Broadcast <- rpc.BuildNotification("vote", out)
}
