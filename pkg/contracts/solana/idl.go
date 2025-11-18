package solana

type IDL struct {
	Address      string        `json:"address"`
	Metadata     Metadata      `json:"metadata"`
	Instructions []Instruction `json:"instructions"`
	Accounts     []Account     `json:"accounts"`
	Errors       []Error       `json:"errors"`
	Types        []Type        `json:"types"`
}

type Metadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Spec        string `json:"spec"`
	Description string `json:"description"`
}

type Instruction struct {
	Name          string    `json:"name"`
	Discriminator []byte    `json:"discriminator"`
	Accounts      []Account `json:"accounts"`
	Args          []Arg     `json:"args"`
}

type Account struct {
	Name     string `json:"name"`
	Writable bool   `json:"writable,omitempty"`
	Signer   bool   `json:"signer,omitempty"`
	Address  string `json:"address,omitempty"`
	PDA      *PDA   `json:"pda,omitempty"`
}

type PDA struct {
	Seeds []Seed `json:"seeds"`
}

type Seed struct {
	Kind  string `json:"kind"`
	Value []byte `json:"value,omitempty"`
}

type Arg struct {
	Name string `json:"name"`
	Type any    `json:"type"`
}

type Error struct {
	Code int    `json:"code"`
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

type Type struct {
	Name string    `json:"name"`
	Type TypeField `json:"type"`
}

type TypeField struct {
	Kind   string  `json:"kind"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name string `json:"name"`
	Type any    `json:"type"`
}
