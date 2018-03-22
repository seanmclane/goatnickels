package block

import(
  "fmt"
  "io/ioutil"
  "os"
  "time"
  "encoding/json"
  "strconv"
  "golang.org/x/crypto/sha3"
  "crypto/elliptic"
  "crypto/ecdsa"
  "crypto/rand"
  "math/big"
  "encoding/hex"
  "net/http"
  "bytes"
)

//define config structure
//does this belong here?
type Config struct {
  Directory string `json:"directory"`
  Nodes []string `json:"nodes"`
  Account string `json:"account"`
}

//load config
func LoadConfig() (config Config) {
  c, err := os.Open("config.json")
  if err != nil {
    panic(err)
  }

  //fix this to json unmarshal
  decoder := json.NewDecoder(c)
  err = decoder.Decode(&config)
  if err != nil {
    fmt.Println("error:", err)
  }

  return config
}

type KeyStore struct {
  PublicKey string `json:"public_key"`
  PrivateKey string `json:"private_key"`
}

func LoadKeyStore() (keystore KeyStore) {
  k, err := os.Open("keystore.json")
  if err != nil {
    panic(err)
  }

  //fix this to json unmarshal
  decoder := json.NewDecoder(k)
  err = decoder.Decode(&keystore)
  if err != nil {
    fmt.Println("error:", err)
  }

  return keystore
}

//initializing blockchain objects here for now
//need to have last validated block (index at minimum)
var LastGoatBlock Block

//need to have last state
var Accounts map[string]Account

//create candidate set of transactions
var CandidateSet []Transaction

//create vote set to store votes from network
var VoteSet []Vote

// create http client for all network requests
var client = &http.Client{
  Timeout: time.Second * 10,
  }

//removing blockchain since it's not needed to store whole chain in memory
// type Blockchain []Block
// var GoatChain Blockchain

type Block struct {
  Index int `json:"index"`
  Timestamp int `json:"timestamp"`
  Data Data `json:"data"`
  LastHash []byte `json:"last_hash"`
  Hash []byte `json:"hash"`
}

type StoredBlock struct {
  Index int `json:"index"`
  Timestamp int `json:"timestamp"`
  Data Data `json:"data"`
  LastHash string `json:"last_hash"`
  Hash string `json:"hash"`
}

type Data struct {
  State map[string]Account `json:"state"`
  Transactions []Transaction `json:"transactions"`
}

type Account struct {
  Balance int `json:"balance"`
  Sequence int `json:"sequence"`
}

type Transaction struct {
  From string `json:"from"`
  To string `json:"to"`
  Amount int `json:"amount"`
  Sequence int `json:"sequence"`
  R string `json:"r"`
  S string `json:"s"`
}

type Signature struct {
  R string `json:"r"`
  S string `json:"s"`
}

type Vote struct {
  Account string `json:"account"`
  Hash string `json:"hash"`
  Signature Signature `json:"signature"`
}

func AsciiGoat() {
  a := "\x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x2d \x2d \x2e \x5f \x2c \x2d \x2d \x2e \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x27 \x20 \x20 \x2c \x27 \x20 \x20 \x20 \x2c \x2d \x60 \x2e \x0a \x28 \x60 \x2d \x2e \x5f \x5f \x20 \x20 \x20 \x20 \x2f \x20 \x20 \x2c \x27 \x20 \x20 \x20 \x2f \x0a \x20 \x60 \x2e \x20 \x20 \x20 \x60 \x2d \x2d \x27 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x5f \x5f \x2c \x2d \x2d \x27 \x2d \x2e \x0a \x20 \x20 \x20 \x60 \x2d \x2d \x2f \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x2d \x2e \x20 \x20 \x5f \x5f \x5f \x5f \x5f \x5f \x2f \x0a \x20 \x20 \x20 \x20 \x20 \x28 \x6f \x2d \x2e \x20 \x20 \x20 \x20 \x20 \x2c \x6f \x2d \x20 \x2f \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x60 \x2e \x20 \x3b \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x7c \x3a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x27 \x60 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x20 \x20 \x20 \x5c \x0a \x20 \x20 \x20 \x20 \x20 \x28 \x6f \x20 \x6f \x20 \x2c \x20 \x20 \x2d \x2d \x27 \x20 \x20 \x20 \x20 \x20 \x3a \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x2d \x2d \x27 \x2c \x27 \x2e \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x3b \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x60 \x3b \x3b \x20 \x20 \x3a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2f \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x3b \x27 \x20 \x20 \x3b \x20 \x20 \x2c \x27 \x20 \x2c \x27 \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x27 \x2c \x27 \x20 \x20 \x3a \x20 \x20 \x27 \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x20 \x5c \x20 \x20 \x20 \x3a \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x60"
  fmt.Println(a, "\n")
}

func (b *Block) HashBlock() {
  //create a hash of all values in the block
  //TODO: handle error
  hash_data, _ := json.Marshal(b.Data)
  block_string := strconv.Itoa(b.Index)+string(hash_data)+hex.EncodeToString(b.LastHash[:])
  fixed_hash := sha3.Sum512([]byte(block_string))
  b.Hash = fixed_hash[:]
}

func CreateGenesisBlock() {
  //temporary
  //manually adding accounts for now
  Accounts = make(map[string]Account)

  Accounts["goat_04dbb67ae9650ca3258071909f74be5400fe53fc2e5dcc82103020f3aeefeee5f9980c4c05bb8696215458dfa7ddaa1505d2826cab3d246b8930b0694f766a22f8bb63932368c0b12bf80cfaee8a18db1d7ce19df0a84215d20b0bbfbd30d95c25"] = Account{
    Balance: 50884323425,
    Sequence: 0,
  }
  Accounts["goat_04ab1594a3b65e440653b1a54952aee3cb7f5c41cb476f7ecd3ce58dc23cef0923beb45fc275ff4149cd9f0417f8ca885e882b3b68d00bab2988b22f2eaf7f6683ba3e672abd668e5788a8ecb4d055cd024f004ff03db06158f18e5bd02914685a"] = Account{
    Balance: 94043534214,
    Sequence: 0,
  }
  Accounts["goat_04c7cb2cef7da5cda83333f34fba7f07b3d1a7572ca909487c7ed20d147706b731e26983c18659bc1caf260a4fd4fc390d9bec208c92d123498faad57ae365ba3aebcd4a93e74802adee03cfbac8f71ed7f5d00824de59bf292c20b2b73bd3228d"] = Account{
    Balance: 38763423645,
    Sequence: 0,
  }
  Accounts["goat_045b4dfabe49048ef6fb6e47fc4e2b33dd54e46b3ed4ab008f8dce7457f588f7a6975690328db4bd48eb874ff909c579fe37ae4f39e9b9b10ac1f2f49083c7d2d8fe91ff5314b2742d58e894681d55682876417f33f851e8091f9c00045a7a9ebc"] = Account{
    Balance: 76457654265,
    Sequence: 0,
  }
  //set arbitrary data
  data := Data{
    State: Accounts,
    Transactions: nil,
  }

  //convert [64]byte to []byte
  fixed_hash := sha3.Sum512([]byte("Goatnickels baby!"))
  hash := fixed_hash[:]

  //genesis block for now
  b := Block {
    Index: 1,
    Timestamp: 0, //TODO: convert this to the birthdate of GoatNickels
    Data: data,
    LastHash: hash,
  }

  b.HashBlock()

  b.WriteBlockToLocalStorage()

}

//TODO: figure out where to keep data structures and have one way imports
type MaxBlockResponse struct {
  MaxBlock int64 `json:"max_block"`
}

func InitializeState() {

  config := LoadConfig()

  max_list := GetMaxBlockFromNetwork()

  //genesis block only created by calling function manually
  //always check the network for max block, then start

  max := FindMaxBlock()

  //check if different and get blocks from network if behind
  //TODO: loop through max_list and find highest block?? most common? or just make get max block only keep the highest returned block?
  if max < max_list[0] {
    GetBlockChainFromNetwork(max, max_list[0], config.Nodes[0])
    max = max_list[0]
  }

  b := ReadBlockFromLocalStorage(strconv.Itoa(int(max)))

  //make bytestring to Block
  err := json.Unmarshal(b, &LastGoatBlock)
  if err != nil {
    fmt.Println("error:", err)
  }

}

func GetMaxBlockFromNetwork() (max_list []int64){
  
  config := LoadConfig()

  for key, node := range config.Nodes {
    max_list = append(max_list, 0)
    r, err := client.Get("http://"+node+":3000/api/v1/maxblock")
    if err != nil {
      fmt.Println("no response from node:", node)
      fmt.Println("error:", err)
    } else {
      defer r.Body.Close()
      var res MaxBlockResponse
      err = json.NewDecoder(r.Body).Decode(&res)
      if err != nil {
        fmt.Println("error:", err)
      }
      max_list[key] = res.MaxBlock
    }
  }

  return max_list
}

func GetBlockChainFromNetwork(local_max int64, network_max int64, node string) {
  for i := local_max+1; i <= network_max; i++ {
    //TODO get block from network, then write to local storage for each block
    r, err := client.Get("http://"+node+":3000/api/v1/block/"+strconv.Itoa(int(i)))
    if err != nil {
      fmt.Println("could not get block from", node)
      fmt.Println("error:", err)
    } else {
      defer r.Body.Close()
      var b Block
      err = json.NewDecoder(r.Body).Decode(&b)
      if err != nil {
        fmt.Println("error:", err)
      }
      b.WriteBlockToLocalStorage()
    }
    //TODO validate blocks!
  }
}

func ReadBlockFromLocalStorage(index string) (b []byte) {
  config := LoadConfig()
  b, _ = ioutil.ReadFile(string(config.Directory)+index)
  return b
}

func FindMaxBlock() (max int64) {
  config := LoadConfig()
  files, err := ioutil.ReadDir(config.Directory)
  if err != nil {
    panic(err)
  }

  max = 0
  for _, file := range files {
    cur, err := strconv.ParseInt(file.Name(), 10, 0)
    if err != nil {
      fmt.Println("error:", err)
    }
    if cur > max {
      max = cur
    }
  }

  return max
}

func MakeNextBlockData() (data Data){
  
  var empty_txions []Transaction

  data = Data{
    State: LastGoatBlock.Data.State,
    Transactions: empty_txions,
  }

  data.ApplyTransactions()
  
  return data
}

func (v *Vote) AddVote() (ok bool) {
  ok = v.VerifyVote()
  if ok != true {
    return false
  }
  VoteSet = append(VoteSet, *v)
  fmt.Println(VoteSet)
  return ok
}

func (v *Vote) VerifyVote() (ok bool) {
  //TODO: check that the vote is not in the voteset already
  //TODO: check that there is not another vote from the same account in this round

  //verify signature of account sending the vote
  //TODO: abstract key recreation into a function (hash, r, s) (ok bool)
  hash := v.HashVote()
  //check that key is well formed
  if len(v.Account) < 100 {
    return false
  }
  //remove goat_ from key
  public_key := v.Account[5:]
  //recreate ecdsa.PublicKey from pub
  byte_key, err := hex.DecodeString(public_key)
  if err != nil {
    fmt.Println("error:", err)
  }
  x, y := elliptic.Unmarshal(elliptic.P384(), byte_key)
  if x == nil || y == nil {
    return false
  }
  pub := new(ecdsa.PublicKey)
  pub.Curve = elliptic.P384()
  pub.X, pub.Y = x, y
  //convert r and s back to big ints
  byte_r, err := hex.DecodeString(v.Signature.R)
  if err != nil {
    fmt.Println("error:", err)
  }
  r := new(big.Int).SetBytes(byte_r)
  byte_s, err := hex.DecodeString(v.Signature.S)
  if err != nil {
    fmt.Println("error:", err)
  }
  s := new(big.Int).SetBytes(byte_s)
  //verify signature
  sig_ok := ecdsa.Verify(pub, hash, r, s)
  
  if sig_ok {
    fmt.Println("vote verified")
    return true
  } else {
    return false
  }
}

func (v *Vote) SignVote() (r string, s string) {
  keystore := LoadKeyStore()
  private_key := keystore.PrivateKey

  hash := v.HashVote()
  //recreate ecdsa.PrivateKey from private_key
  byte_key, err := hex.DecodeString(private_key)
  if err != nil {
    fmt.Println("error:", err)
  }
  bigint_key := new(big.Int).SetBytes(byte_key)
  priv := new(ecdsa.PrivateKey)
  priv.PublicKey.Curve = elliptic.P384()
  priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(byte_key)
  priv.D = bigint_key
  r_int, s_int, err := ecdsa.Sign(rand.Reader, priv, hash)
  if err != nil {
    fmt.Println("error:", err)
  }
  r = hex.EncodeToString(r_int.Bytes())
  s = hex.EncodeToString(s_int.Bytes())
  return r, s
}

func (v *Vote) HashVote() (h []byte) {
  hash_string := v.Account+v.Hash
  fixed_hash := sha3.Sum512([]byte(hash_string))
  h = fixed_hash[:]
  return h
}

func SendVoteToNetwork() {
  config := LoadConfig()

  v := Vote {
    Account: config.Account,
    Hash: hex.EncodeToString(HashCandidateSet(&CandidateSet)),
  }

  r, s := v.SignVote()

  v.Signature = Signature{
    R: r,
    S: s,
  }

  for _, node := range config.Nodes {
    json, err := json.Marshal(v)
    if err != nil {
      fmt.Println("error:", err)
    }
    req, err := http.NewRequest("POST", "http://"+node+":3000/api/v1/vote", bytes.NewBuffer(json))
    if err != nil {
      fmt.Println("error:", err)
    }
    req.Header.Set("Content-Type", "application/json")
    r, err := client.Do(req)
    if err != nil {
      fmt.Println("error:", err)
    } else {
      fmt.Println("Status:", r.Status)
    }
  }
}

func CheckConsensus() {
  //TODO: final criteria for consensus = 2/3 of stakes sign hash of candidate set transaction

  cs_hash := HashCandidateSet(&CandidateSet)

  var match int
  var total int

  //TODO: get proportions for all, this won't work if this node is the one out of sync
  //changed the consensus to check votes in voteset
  //if agreed-upon hash is same as local hash, apply transactions
  //TODO: if not, wait and request new block

  for _, v := range VoteSet {
    v_hash, err := hex.DecodeString(v.Hash)
    if err != nil {
      fmt.Println("error:", err)
    }
    if bytes.Equal(v_hash, cs_hash) {
      match += 1
    }
    total += 1
  }

  if total < 1 {
    fmt.Println("No consensus reached due to lack of votes")
    return
  }

  if match / total >= 2/3 {
    NextBlock()
  } else {
    //TODO: tell network to restart consensus round
    fmt.Println("No consensus reached")
  }
}

func HashCandidateSet(cs *[]Transaction) (h []byte){
  var sum string
  for _, txion := range *cs {
    sum += hex.EncodeToString(txion.HashTransaction())
  }
  fixed_hash := sha3.Sum512([]byte(sum))
  h = fixed_hash[:]

  fmt.Println(hex.EncodeToString(h))
 
  return h
}

func NextBlock() {

  next_block := Block {
    Index: LastGoatBlock.Index+1,
    Timestamp: int(time.Now().UTC().Unix()),
    Data: MakeNextBlockData(),
    LastHash: LastGoatBlock.Hash,
  }
  
  next_block.HashBlock()

  DescribeBlock(next_block)

  LastGoatBlock = next_block

  next_block.WriteBlockToLocalStorage()

}

func DescribeBlock(b Block) {
  fmt.Printf("----------------------------------------------------------------------------------------------\n\n")
  fmt.Println("Block ID:", b.Index)
  fmt.Printf("\n---Block State---\n")
  for key, val := range b.Data.State {
    fmt.Printf("Account: %s\nBalance: %d\nSequence: %d\n", key, val.Balance, val.Sequence)
  }
  fmt.Printf("\n---Block Transactions---\n")
  for _, txion := range b.Data.Transactions {
    fmt.Printf("To: %s\nFrom: %s\nAmount: %d\n", txion.To, txion.From, txion.Amount)
  }
  fmt.Printf("\n---Hashes---\n")
  fmt.Println("Last Hash:", hex.EncodeToString(b.LastHash[:]))
  fmt.Println("Block Hash:", hex.EncodeToString(b.Hash[:]))
  fmt.Println("\n")
  fmt.Println("Block Time:", time.Unix(int64(b.Timestamp),0))
  fmt.Println("----------------------------------------------------------------------------------------------")
}

func (t *Transaction) AddTransaction() (ok bool) {
  ok = t.VerifyTransaction()
  if ok != true {
    return false
  }
  //check if transaction exists in candidate set
  //if so, return and don't broadcast
  for _, c := range CandidateSet {
    if *t == c {
      return
    }
  }
  CandidateSet = append(CandidateSet, *t)
  t.Broadcast()
  return ok
}

func (d *Data) ApplyTransactions() {
  //add and subtract from accounts
  for _, txion := range CandidateSet {
    //check if sequence is incremented by one
    ok := txion.VerifySequence()
    if ok == true {
      //increment account sequence and debit balance
      fnb := d.State[txion.From].Balance - txion.Amount
      d.State[txion.From] = Account{
        Balance: fnb,
        Sequence: txion.Sequence,
      }
      //credit balance
      tnb := d.State[txion.To].Balance + txion.Amount
      d.State[txion.To] = Account{
        Balance: tnb,
        Sequence: d.State[txion.To].Sequence,
      } 
      //add valid transaction to transaction list
      d.Transactions = append(d.Transactions, txion)
    } else {
      //TODO: else mark transaction as failed and do something with it? broadcast to network?

    }
  }
  //reset candidate transactions to apply
  CandidateSet = nil

}

func (b *Block) WriteBlockToLocalStorage() {
  config := LoadConfig()

  //convert data to plain json
  out, err := json.Marshal(b)
  if err != nil {
    fmt.Println("error:", err)
  }

  //write json to file at config directory
  //TODO: check if file exists and don't overwrite
  err = ioutil.WriteFile(string(config.Directory)+strconv.Itoa(b.Index), out, 0644)
  if err != nil {
      panic(err)
  }

  fmt.Println("Block "+strconv.Itoa(b.Index)+" written successfully!")

}

type AccountResponse struct {
  PrivateKey string `json:"private_key"`
  PublicKey string `json:"public_key"`
}

func GenerateAccount() {
  //create the keypair
  priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
  if err != nil {
    fmt.Println("error:", err)
  }

  //create the address from the public key variables
  pub := priv.PublicKey
  pubkey := elliptic.Marshal(elliptic.P384(), pub.X, pub.Y)

  response := AccountResponse{
    PrivateKey: hex.EncodeToString(priv.D.Bytes()),
    PublicKey: "goat_"+hex.EncodeToString(pubkey),
  }

  bytes, err := json.Marshal(response)
  if err != nil {
    fmt.Println("error:", err)
  }

  fmt.Println(string(bytes))
}

//TODO: make this real and not a test of some hardcoded values
func (t *Transaction) SignTransaction (private_key string) (r, s string) {
  hash := t.HashTransaction()
  //recreate ecdsa.PrivateKey from private_key
  byte_key, err := hex.DecodeString(private_key)
  if err != nil {
    fmt.Println("error:", err)
  }
  bigint_key := new(big.Int).SetBytes(byte_key)
  priv := new(ecdsa.PrivateKey)
  priv.PublicKey.Curve = elliptic.P384()
  priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(byte_key)
  priv.D = bigint_key
  r_int, s_int, err := ecdsa.Sign(rand.Reader, priv, hash)
  if err != nil {
    fmt.Println("error:", err)
  }
  r = hex.EncodeToString(r_int.Bytes())
  s = hex.EncodeToString(s_int.Bytes())
  return r, s
}

func (t *Transaction) HashTransaction() (h []byte) {
  hash_string := t.To+t.From+strconv.Itoa(t.Amount)+strconv.Itoa(t.Sequence)
  fixed_hash := sha3.Sum512([]byte(hash_string))
  h = fixed_hash[:]
  return h
}

func (t *Transaction) VerifyTransaction() (ok bool) {
  //check if t.R and t.S ok with public key
  //what is being signed exactly? hash of transaction sequence, to, from, and amount
  hash := t.HashTransaction()
  //public_key := "goat_04c12951412edfc215fe6d288491eb1251e2d8d99375c01049588dd228c6346f068246353d84702418f797d672af512d89742f6842b32f43541ea703f08170a67687f75fe0c6f15bd518764dee5476c86f9ba33f28036a76d018c1d7c8b14c307f"
  //check that key is well formed
  if len(t.From) < 100 {
    return false
  }
  //remove goat_ from key
  public_key := t.From[5:]
  //recreate ecdsa.PublicKey from pub
  byte_key, err := hex.DecodeString(public_key)
  if err != nil {
    fmt.Println("error:", err)
  }
  x, y := elliptic.Unmarshal(elliptic.P384(), byte_key)
  if x == nil || y == nil {
    return false
  }
  pub := new(ecdsa.PublicKey)
  pub.Curve = elliptic.P384()
  pub.X, pub.Y = x, y
  //convert r and s back to big ints
  byte_r, err := hex.DecodeString(t.R)
  if err != nil {
    fmt.Println("error:", err)
  }
  r := new(big.Int).SetBytes(byte_r)
  byte_s, err := hex.DecodeString(t.S)
  if err != nil {
    fmt.Println("error:", err)
  }
  s := new(big.Int).SetBytes(byte_s)
  //verify signature
  sig_ok := ecdsa.Verify(pub, hash, r, s)
  //verify balance is sufficient
  spend_ok := t.VerifyNegativeSpend()
  //do not check sequence here, so you can have more than one transaction per block
  if sig_ok && spend_ok {
    return true
  } else {
    return false
  }

}

func (t *Transaction) VerifyNegativeSpend() (ok bool) {
  if LastGoatBlock.Data.State[t.From].Balance < t.Amount {
    return false
  } else {
    return true
  }
}

func (t *Transaction) VerifySequence() (ok bool) {
  //sequence must be current account sequence number plus one
  if t.Sequence == LastGoatBlock.Data.State[t.From].Sequence + 1 {
    return true
  } else {
    return false
  }
}

func (t *Transaction) Broadcast() {

  config := LoadConfig()

  for _, node := range config.Nodes {
    json, err := json.Marshal(t)
    if err != nil {
      fmt.Println("error:", err)
    }
    req, err := http.NewRequest("POST", "http://"+node+":3000/api/v1/txion", bytes.NewBuffer(json))
    if err != nil {
      fmt.Println("error:", err)
    }
    req.Header.Set("Content-Type", "application/json")
    r, err := client.Do(req)
    if err != nil {
      fmt.Println("error:", err)
    } else {
      fmt.Println("Status:", r.Status)
    }
  }
}