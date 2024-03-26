package dataverse

type Select struct {
	Query SelectQuery `json:"query"`
}

type SelectQuery struct {
	Prefixes []Prefix         `json:"prefixes"`
	Select   []SelectItem     `json:"select"`
	Where    []WhereCondition `json:"where"`
	Limit    uint64           `json:"limit"`
}

type SelectResponse struct {
	Head struct {
		Vars []string `json:"vars"`
	} `json:"head"`
	Results struct {
		Bindings []map[string]struct {
			Type     string      `json:"type"`
			Value    interface{} `json:"value"`
			Lang     *string     `json:"xml:lang,omitempty"`
			Datatype *IRI        `json:"datatype,omitempty"`
		} `json:"bindings"`
	} `json:"results"`
}

type Prefix struct {
	Prefix    string `json:"prefix"`
	Namespace string `json:"namespace"`
}

type SelectItem struct {
	Variable string `json:"variable"`
}

type WhereCondition struct {
	Simple SimpleWhereCondition `json:"simple"`
}

type SimpleWhereCondition struct {
	TriplePattern TriplePattern `json:"triple_pattern"`
}

type TriplePattern struct {
	Subject   VarOrNode          `json:"subject"`
	Predicate VarOrNamedNode     `json:"predicate"`
	Object    VarOrNodeOrLiteral `json:"object"`
}

type VarOrNode struct {
	Variable string `json:"variable,omitempty"`
	Node     *Node  `json:"node,omitempty"`
}

type VarOrNamedNode struct {
	Variable  string `json:"variable,omitempty"`
	NamedNode *IRI   `json:"named_node,omitempty"`
}

type VarOrNodeOrLiteral struct {
	Variable string   `json:"variable,omitempty"`
	Node     *Node    `json:"node,omitempty"`
	Literal  *Literal `json:"literal,omitempty"`
}

type Node struct {
	NamedNode *IRI   `json:"named_node,omitempty"`
	BlankNode string `json:"blank_node,omitempty"`
}

type IRI struct {
	Prefixed string `json:"prefixed,omitempty"`
	Full     string `json:"full,omitempty"`
}

type Literal struct {
	Simple               string `json:"simple,omitempty"`
	LanguageTaggedString *struct {
		Value    string `json:"value"`
		Language string `json:"language"`
	} `json:"language_tagged_string,omitempty"`
	TypedValue *struct {
		Value    string `json:"value"`
		Datatype IRI    `json:"datatype"`
	} `json:"typed_value,omitempty"`
}