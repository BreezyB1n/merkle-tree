package merkleTree

import (
	"crypto/sha256"
	"encoding/hex"
)

type MerkleTree struct {
	Root *MerkleTreeNode
}

type MerkleTreeNode struct {
	Left  *MerkleTreeNode
	Right *MerkleTreeNode
	Hash  string
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []*MerkleTreeNode

	// 创建叶子节点
	for _, d := range data {
		node := NewMerkleNode(nil, nil, d)
		nodes = append(nodes, node)
	}

	// 构建 Merkle 树
	for len(nodes) > 1 {
		var newLevel []*MerkleTreeNode

		// 每两个节点合并为一个节点
		for i := 0; i < len(nodes); i += 2 {
			left := nodes[i]
			var right *MerkleTreeNode
			if i+1 < len(nodes) {
				right = nodes[i+1]
			} else {
				// 节点个数为奇数，复制最后一个节点
				right = nodes[i]
			}

			parent := NewMerkleNode(left, right, nil)
			newLevel = append(newLevel, parent)
		}

		nodes = newLevel
	}
	return &MerkleTree{Root: nodes[0]}
}

func NewMerkleNode(left, right *MerkleTreeNode, data []byte) *MerkleTreeNode {
	node := &MerkleTreeNode{
		Left:  left,
		Right: right,
	}

	// 判断是否为叶子节点
	if data != nil {
		node.Hash = hash(data)
	} else {
		// 如果是父节点，对左右子节点值进行XOR运算
		node.Hash = hex.EncodeToString(xor(hexToBytes(node.Left.Hash), hexToBytes(node.Right.Hash)))
	}

	return node
}

func xor(a, b []byte) []byte {
	result := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}
	return result
}

func hexToBytes(hexStr string) []byte {
	bytes, _ := hex.DecodeString(hexStr)
	return bytes
}

func hash(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// GetMerkleRoot 返回 Merkle Tree 的根哈希
func (mt *MerkleTree) GetMerkleRoot() string {
	return mt.Root.Hash
}

func (mt *MerkleTree) GetMerkleProof(data []byte) []string {
	var proof []string

	targetHash := hash(data)
	currentNode := mt.Root

	// 从根节点开始遍历
	for currentNode.Left != nil || currentNode.Right != nil {
		if contains(currentNode.Left, targetHash) {
			proof = append(proof, currentNode.Right.Hash)
			currentNode = currentNode.Left
		} else {
			proof = append(proof, currentNode.Left.Hash)
			currentNode = currentNode.Right
		}
	}
	return proof
}

func contains(node *MerkleTreeNode, targetHash string) bool {
	if node == nil {
		return false
	}
	return node.Hash == targetHash
}
