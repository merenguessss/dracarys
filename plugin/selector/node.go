package selector

// Node 服务结点.
type Node struct {
	Key   string
	Value string
	// Meta 服务结点的内容,包括黑白名单,权重等.
	Meta map[string]string
}
