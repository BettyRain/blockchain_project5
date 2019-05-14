package p1

import (
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
	"strings"
)

type Flag_value struct {
	encoded_prefix []uint8
	value          string //here can be hashed_value of node (if branch or ext)
}

type Node struct {
	node_type    int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value   Flag_value
}

type MerklePatriciaTrie struct {
	db map[string]Node
	kv map[string]string //key value pair
	ks map[string][]byte
	//string is a hashed_value of Node
	root string
}

func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {
	key_hex := keyToHex(key)
	root_node := mpt.db[mpt.root]
	n := Node{}
	switch root_node.node_type {
	case 1:
		n = mpt.GetByNode(root_node, []byte(key_hex))
	case 2:
		root_key := root_node.flag_value.encoded_prefix
		root_key_decoded := compact_decode(root_key)
		if intersectionCount(string(root_key_decoded), string(key_hex)) > 0 {
			n = mpt.GetByNode(root_node, []byte(key_hex))
		} else {
			return "", errors.New("path_not_found")
		}

	}

	switch n.node_type {
	case 1:
		if n.branch_value[16] != "" {
			return n.branch_value[16], errors.New("")
		} else {
			return "", errors.New("path_not_found")
		}

	case 2:
		if n.flag_value.value != "" {
			return n.flag_value.value, errors.New("")
		} else {
			return "", errors.New("path_not_found")
		}
	}

	return "", errors.New("path_not_found")
}

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string, sign []byte) {
	n := Node{}
	//add key-value pair
	if mpt.kv == nil {
		mpt.kv = make(map[string]string)
	}
	mpt.kv[key] = new_value

	//add signatures
	if mpt.ks == nil {
		mpt.ks = make(map[string][]byte)
	}
	mpt.ks[key] = sign

	encoded_key := keyToHex(key)
	root_node := mpt.db[mpt.root]

	switch root_node.node_type {
	case 1:
		n = mpt.GetByNode(root_node, []byte(encoded_key))
	case 2:
		root_key := root_node.flag_value.encoded_prefix
		root_key_decoded := compact_decode(root_key)
		if intersectionCount(string(root_key_decoded), string(encoded_key)) > 0 {
			n = mpt.GetByNode(root_node, []byte(encoded_key))
		}
	}
	if n.node_type == 2 {
		n_hash := mpt.KeyByValue(n)
		n.flag_value.value = new_value
		mpt.db[n_hash] = n
	} else if n.node_type == 1 {
		n_hash := mpt.KeyByValue(n)
		n.branch_value[16] = new_value
		mpt.db[n_hash] = n
	} else {
		//if mpt is empty
		if mpt.db == nil {
			encoded_key = append(encoded_key, 16)
			encoded_key = compact_encode(encoded_key)
			n.node_type = 2
			n.branch_value = [17]string{}
			n.flag_value.encoded_prefix = encoded_key
			n.flag_value.value = new_value
			root := n.hash_node()
			mpt.root = root
			mpt.db = map[string]Node{
				root: n,
			}
		} else {
			encoded_key = append(encoded_key)
			encoded_key = compact_encode(encoded_key)
			mpt.InsertByNode(mpt.db[mpt.root], []byte(encoded_key), new_value, true)
		}
	}

	//	previousHash := root_node.flag_value.value
	//	rootHash := root_node.hash_node()
	//	root_node.flag_value.value = rootHash
	//	mpt.Delete(previousHash)

	//	mpt.db[rootHash] = root_node
	//	mpt.root = rootHash

}

func (mpt *MerklePatriciaTrie) InsertByNode(node Node, key []byte, new_value string, isRoot bool) string {

	if ((len(key) == 0) && node.node_type != 1) || (node.node_type == 0) {
		return ""
	}

	key_node := node.flag_value.encoded_prefix
	hash_node := node.flag_value.value //hash_value of Node
	hashed_old_leaf := mpt.KeyByValue(node)

	switch node.node_type {
	//branch node
	case 1:
		decoded_key := []uint8{}
		branch_value := -1
		rest_path := []uint8{}
		if key != nil {
			decoded_key = compact_decode(key)
			branch_value = int(decoded_key[0])
			rest_path = append(rest_path, decoded_key[1:]...)
		} else {
			branch_value = 16
		}

		rest_path = append(rest_path, 16)
		rest_path = compact_encode(rest_path)

		if node.branch_value[branch_value] == "" {
			//cell is empty
			//create new leaf
			f_leaf := Flag_value{encoded_prefix: rest_path, value: new_value}
			n_leaf := Node{node_type: 2, branch_value: [17]string{}, flag_value: f_leaf}
			n_leaf_hash := n_leaf.hash_node()

			// add leaf to cell
			node.branch_value[branch_value] = n_leaf.hash_node()
			mpt.db[hashed_old_leaf] = node
			mpt.db[n_leaf_hash] = n_leaf

		} else {
			//cell is not empty
			//move to next node
			hash_next_node := node.branch_value[branch_value]
			next_node := mpt.db[hash_next_node]
			hash_next_node = mpt.InsertByNode(next_node, rest_path, new_value, false)

			if hash_next_node != "" {
				node.branch_value[branch_value] = hash_next_node
			}
			mpt.db[hashed_old_leaf] = node
		}

		//NEWMEOW
		old_hash := mpt.KeyByValue(node)
		if old_hash != node.hash_node() {
			delete(mpt.db, old_hash)
			new_hash := node.hash_node()
			mpt.db[new_hash] = node
			if mpt.root == old_hash {
				mpt.root = new_hash
			} else {
				mpt.ChangeHash(old_hash, new_hash)
			}
		}

	//ext or leaf
	case 2:
		node_prefix := key_node[0] / 16

		//check for leaf and ext
		switch node_prefix {
		//extension node
		case 0, 1:
			decoded_key_node := compact_decode(key_node) //nibbles of ext decoded
			decoded_key := compact_decode(key)
			//new leaf with empty prefix
			f_leaf := Flag_value{encoded_prefix: []uint8{}, value: new_value}
			n_leaf := Node{node_type: 2, branch_value: [17]string{}, flag_value: f_leaf}
			n_leaf_hash := n_leaf.hash_node()
			index_ext := intersectionCount(string(decoded_key_node), string(decoded_key))

			if index_ext == len(decoded_key_node) {
				//all same nibbles
				branch_node := mpt.db[hash_node]
				//here we go to insert branch node
				key_next := []uint8{}
				if index_ext == len(decoded_key) {
					key_next = nil
				} else {
					key_next = decoded_key[index_ext:]
					key_next = compact_encode(key_next)
				}
				mpt.InsertByNode(branch_node, key_next, new_value, false)

			} else if (index_ext <= (len(decoded_key_node) - 1)) && (index_ext > 0) {
				//some same nibbles - more than one left (<)
				//take same for new ext + take one for branch + leave other for old ext
				//some same nibbles - one nibble left (=)
				//take same for new ext + take one for branch + delete old ext + add old branch to index[] of new branch

				//have to change the root
				common_path := decoded_key_node[:index_ext]
				rest := decoded_key[index_ext:]
				branch_nib := decoded_key_node[index_ext:]

				//create branch
				branch := [17]string{}

				if (len(rest)) > 0 && (len(branch_nib)) > 1 {
					branch_path := rest[0]
					branch_nibble := branch_nib[0]
					rest_path := decoded_key[(index_ext + 1):]
					rest_nibble := decoded_key_node[(index_ext + 1):]

					branch[branch_path] = n_leaf_hash
					branch[branch_nibble] = hashed_old_leaf

					rest_path = append(rest_path, 16)
					rest_path = compact_encode(rest_path)
					rest_nibble = compact_encode(rest_nibble)

					//change prefixes
					node.flag_value.encoded_prefix = rest_nibble
					n_leaf.flag_value.encoded_prefix = rest_path
					mpt.db[hashed_old_leaf] = node
					mpt.db[n_leaf_hash] = n_leaf

				} else if len(rest) == 0 && len(branch_nib) > 1 {
					branch_nibble := branch_nib[0]
					rest_nibble := decoded_key_node[(index_ext + 1):]
					branch[branch_nibble] = hashed_old_leaf
					branch[16] = new_value

					rest_nibble = compact_encode(rest_nibble)
					node.flag_value.encoded_prefix = rest_nibble
					mpt.db[hashed_old_leaf] = node

				} else if len(rest) == 0 && len(branch_nib) == 1 {
					//delete ext
					//add old branch to new
					branch_nibble := branch_nib[0]
					branch[branch_nibble] = hash_node
					delete(mpt.db, hashed_old_leaf)
					branch[16] = new_value

				} else if len(branch_nib) == 1 {
					//delete ext
					//add old branch to new
					branch_nibble := branch_nib[0]
					branch[branch_nibble] = hash_node
					delete(mpt.db, hashed_old_leaf)

					branch_path := rest[0]
					rest_path := decoded_key[(index_ext + 1):]
					branch[branch_path] = n_leaf_hash
					rest_path = append(rest_path, 16)
					rest_path = compact_encode(rest_path)

					n_leaf.flag_value.encoded_prefix = rest_path
					mpt.db[hashed_old_leaf] = node
					mpt.db[n_leaf_hash] = n_leaf
				}
				f_br := Flag_value{encoded_prefix: []uint8{}, value: ""}
				n_br := Node{node_type: 1, branch_value: branch, flag_value: f_br}
				n_br_hash := n_br.hash_node()
				mpt.db[n_br_hash] = n_br

				//create ext
				common_path = compact_encode(common_path)
				f_ext := Flag_value{encoded_prefix: common_path, value: n_br_hash}
				n_ext := Node{node_type: 2, branch_value: [17]string{}, flag_value: f_ext}
				n_ext_hash := n_ext.hash_node()
				mpt.db[n_ext_hash] = n_ext
				if isRoot == true {
					mpt.root = n_ext_hash
				}
				return n_ext_hash

			} else {
				//no same nibbles - more than one in ext
				//create a branch + take one for branch + add ext to branch []
				//no same - one in ext
				//create a branch + take nibble for branch + delete ext + add old branch to [] in new branch
				//have to change the root -> to branch

				//create branch
				branch := [17]string{}
				first_key := decoded_key[0]
				first_node_key := decoded_key_node[0]

				if len(decoded_key_node) > 1 {
					branch[first_key] = n_leaf_hash
					branch[first_node_key] = hashed_old_leaf

					rest_nibble := []uint8{}
					rest_nibble = append(rest_nibble, decoded_key_node[1:]...)
					rest_nibble = compact_encode(rest_nibble)

					rest_path := []uint8{}
					rest_path = append(rest_path, decoded_key[1:]...)
					rest_path = append(rest_path, 16)
					rest_path = compact_encode(rest_path)

					//change prefixes
					node.flag_value.encoded_prefix = rest_nibble
					n_leaf.flag_value.encoded_prefix = rest_path

					mpt.db[hashed_old_leaf] = node
					mpt.db[n_leaf_hash] = n_leaf

				} else {
					branch[first_key] = n_leaf_hash
					branch[first_node_key] = hash_node
					delete(mpt.db, hashed_old_leaf)

					rest_path := []uint8{}
					rest_path = append(rest_path, decoded_key[1:]...)
					rest_path = append(rest_path, 16)
					rest_path = compact_encode(rest_path)

					//change prefixes
					n_leaf.flag_value.encoded_prefix = rest_path
					mpt.db[n_leaf_hash] = n_leaf
				}
				f_br := Flag_value{encoded_prefix: []uint8{}, value: ""}
				n_br := Node{node_type: 1, branch_value: branch, flag_value: f_br}
				n_br_hash := n_br.hash_node()
				mpt.db[n_br_hash] = n_br
				if isRoot == true {
					mpt.root = n_br.hash_node()
				}
				return n_br_hash
			}

		//leaf node
		case 2, 3:
			//check the key
			//keys are equal
			decoded_key_node := compact_decode(key_node)
			decoded_key := compact_decode(key)

			if len(decoded_key_node) == 0 {
				parent_node_hash, parent_index := mpt.FindParentNode(node, hashed_old_leaf)
				parent_node := mpt.db[parent_node_hash]

				hashed_old_leaf := mpt.KeyByValue(node)
				old_leaf := mpt.db[hashed_old_leaf]
				delete(mpt.db, hashed_old_leaf)

				branch := [17]string{}
				rest_path := []uint8{}
				rest_path = append(rest_path, decoded_key[1:]...)
				rest_path = append(rest_path, 16)

				f_leaf := Flag_value{encoded_prefix: compact_encode(rest_path), value: new_value}
				n_leaf := Node{node_type: 2, branch_value: [17]string{}, flag_value: f_leaf}
				n_leaf_hash := n_leaf.hash_node()

				branch[decoded_key[0]] = n_leaf_hash
				branch[16] = old_leaf.flag_value.value

				f_br := Flag_value{encoded_prefix: []uint8{}, value: ""}
				n_br := Node{node_type: 1, branch_value: branch, flag_value: f_br}
				n_br_hash := n_br.hash_node()

				if parent_node.node_type == 1 {
					parent_node.branch_value[parent_index] = n_br_hash
				} else {
					parent_node.flag_value.value = n_br_hash
				}

				mpt.db[n_br_hash] = n_br
				mpt.db[n_leaf_hash] = n_leaf
				mpt.db[parent_node_hash] = parent_node

				parent_node = mpt.db[parent_node_hash]
				return n_br_hash

			} else if len(decoded_key) == 0 {

				parent_node_hash, parent_index := mpt.FindParentNode(node, hashed_old_leaf)
				parent_node := mpt.db[parent_node_hash]

				hashed_old_leaf := mpt.KeyByValue(node)
				old_leaf := mpt.db[hashed_old_leaf]
				prefix := compact_decode(old_leaf.flag_value.encoded_prefix)

				branch := [17]string{}
				branch[prefix[0]] = hashed_old_leaf
				branch[16] = new_value

				prefix = prefix[1:]
				prefix = append(prefix, 16)
				old_leaf.flag_value.encoded_prefix = compact_encode(prefix)

				f_br := Flag_value{encoded_prefix: []uint8{}, value: ""}
				n_br := Node{node_type: 1, branch_value: branch, flag_value: f_br}
				n_br_hash := n_br.hash_node()

				if parent_node.node_type == 1 {
					parent_node.branch_value[parent_index] = n_br_hash
				} else {
					parent_node.flag_value.value = n_br_hash
				}
				mpt.db[n_br_hash] = n_br
				mpt.db[hashed_old_leaf] = old_leaf
				mpt.db[parent_node_hash] = parent_node
				return n_br_hash

			} else if decoded_key_node[0] == decoded_key[0] {
				hashed_old_leaf := mpt.KeyByValue(node)
				old_leaf := mpt.db[hashed_old_leaf]
				delete(mpt.db, hashed_old_leaf)

				index_ext := intersectionCount(string(decoded_key_node), string(decoded_key))
				common_path := decoded_key_node[:index_ext]
				rest := decoded_key[index_ext:]
				branch_nib := decoded_key_node[index_ext:]

				//create new leaf
				f_leaf := Flag_value{encoded_prefix: []uint8{}, value: new_value}
				n_leaf := Node{node_type: 2, branch_value: [17]string{}, flag_value: f_leaf}
				n_leaf_hash := n_leaf.hash_node()

				//create branch
				//if common = key -> add value
				//decrease on one value
				branch := [17]string{}

				if (len(rest)) > 0 && (len(branch_nib)) > 0 {
					branch_path := rest[0]
					branch_nibble := branch_nib[0]
					rest_path := decoded_key[(index_ext + 1):]
					rest_nibble := decoded_key_node[(index_ext + 1):]

					branch[branch_path] = n_leaf_hash
					branch[branch_nibble] = hashed_old_leaf

					rest_path = append(rest_path, 16)
					rest_path = compact_encode(rest_path)
					rest_nibble = append(rest_nibble, 16)
					rest_nibble = compact_encode(rest_nibble)

					//change paths in leafs
					node.flag_value.encoded_prefix = rest_nibble
					n_leaf.flag_value.encoded_prefix = rest_path

					mpt.db[hashed_old_leaf] = node
					mpt.db[n_leaf_hash] = n_leaf

				} else if len(rest) == 0 {
					branch_nibble := branch_nib[0]
					rest_nibble := decoded_key_node[(index_ext + 1):]
					branch[branch_nibble] = hashed_old_leaf
					branch[16] = new_value

					rest_nibble = append(rest_nibble, 16)
					rest_nibble = compact_encode(rest_nibble)

					node.flag_value.encoded_prefix = rest_nibble
					mpt.db[hashed_old_leaf] = node

				} else if len(branch_nib) == 0 {
					branch_path := rest[0]
					rest_path := decoded_key[(index_ext + 1):]
					branch[branch_path] = n_leaf_hash

					rest_path = append(rest_path, 16)
					rest_path = compact_encode(rest_path)
					//change old leaf
					branch[16] = old_leaf.flag_value.value
					n_leaf.flag_value.encoded_prefix = rest_path
					mpt.db[n_leaf_hash] = n_leaf
				}

				f_br := Flag_value{encoded_prefix: []uint8{}, value: ""}
				n_br := Node{node_type: 1, branch_value: branch, flag_value: f_br}
				n_br_hash := n_br.hash_node()
				mpt.db[n_br_hash] = n_br

				//create ext
				common_path = compact_encode(common_path)
				f_ext := Flag_value{encoded_prefix: common_path, value: n_br_hash}
				n_ext := Node{node_type: 2, branch_value: [17]string{}, flag_value: f_ext}
				n_ext_hash := n_ext.hash_node()
				mpt.db[n_ext_hash] = n_ext
				if isRoot == true {
					mpt.root = n_ext_hash
				}
				return n_ext_hash
			} else {
				//not equal keys
				hashed_old_leaf := mpt.KeyByValue(node)
				delete(mpt.db, hashed_old_leaf)

				//create new leaf
				f_leaf := Flag_value{encoded_prefix: decoded_key, value: new_value}
				n_leaf := Node{node_type: 2, branch_value: [17]string{}, flag_value: f_leaf}
				n_leaf_hash := n_leaf.hash_node()

				//create branch
				//decrease on one value
				branch := [17]string{}

				if (len(decoded_key)) > 0 && (len(decoded_key_node)) > 0 {
					br_path_value := decoded_key[0]
					path := []uint8{}
					decoded_key = append(path, decoded_key[1:]...)

					br_nibble_value := decoded_key_node[0]
					nibble := []uint8{}
					decoded_key_node = append(nibble, decoded_key_node[1:]...)

					branch[br_path_value] = n_leaf_hash
					branch[br_nibble_value] = hashed_old_leaf

					decoded_key = append(decoded_key, 16)
					decoded_key = compact_encode(decoded_key)
					decoded_key_node = append(decoded_key_node, 16)
					decoded_key_node = compact_encode(decoded_key_node)

					//change paths in leafs
					node.flag_value.encoded_prefix = decoded_key_node
					n_leaf.flag_value.encoded_prefix = decoded_key

					mpt.db[hashed_old_leaf] = node
					mpt.db[n_leaf_hash] = n_leaf

				} else if len(decoded_key) == 0 {
					br_nibble_value := decoded_key_node[0]
					nibble := []uint8{}
					decoded_key_node = append(nibble, decoded_key_node[1:]...)

					branch[br_nibble_value] = hashed_old_leaf
					branch[16] = new_value

					decoded_key_node = append(decoded_key_node, 16)
					decoded_key_node = compact_encode(decoded_key_node)

					node.flag_value.encoded_prefix = decoded_key_node
					mpt.db[hashed_old_leaf] = node

				} else if len(decoded_key_node) == 0 {
					br_path_value := decoded_key[0]
					path := []uint8{}
					decoded_key = append(path, decoded_key[1:]...)

					branch[br_path_value] = n_leaf_hash

					decoded_key = append(decoded_key, 16)
					decoded_key = compact_encode(decoded_key)
					//change old leaf

					branch[16] = hash_node
					n_leaf.flag_value.encoded_prefix = decoded_key
					mpt.db[n_leaf_hash] = n_leaf

				}

				f_br := Flag_value{encoded_prefix: []uint8{}, value: ""}
				n_br := Node{node_type: 1, branch_value: branch, flag_value: f_br}
				n_br_hash := n_br.hash_node()
				mpt.db[n_br_hash] = n_br
				if isRoot == true {
					mpt.root = n_br_hash
				}
				return n_br_hash
			}
		}
		//NEWMEOW
		/*	old_hash := mpt.KeyByValue(node)
			if old_hash != node.hash_node(){
				delete(mpt.db, old_hash)
				new_hash := node.hash_node()
				mpt.db[new_hash] = node
				if mpt.root == old_hash {
					mpt.root = new_hash
				}
				mpt.ChangeHash(old_hash, new_hash)
			}
		*/
	}

	return ""
}

func (mpt *MerklePatriciaTrie) GetByNode(node Node, key []byte) Node {
	n := Node{}
	key_node := node.flag_value.encoded_prefix
	isExt := true
	if len(key_node) > 0 {
		node_prefix := key_node[0] / 16
		if node_prefix == 2 || node_prefix == 3 {
			isExt = false
		}
	}

	if len(key) == 0 && isExt {
		return node
	}

	if node.node_type == 0 {
		return n
	}

	switch node.node_type {
	case 2:
		key_decoded := compact_decode(key_node)
		hash_node := node.flag_value.value //hash_value of child Node
		n = mpt.db[hash_node]              //child node (next node)
		node_prefix := key_node[0] / 16

		//check for leaf and ext
		switch node_prefix {
		//extension node
		case 0, 1:
			//decerase nibbles
			index := intersectionCount(string(key_decoded), string(key))
			if len(key) > len(key_decoded) && index > 0 {
				//decrease key from beginning
				return mpt.GetByNode(n, key[index:])
				//same key all - go to branch node - return value
			} else if len(key) == len(key_decoded) && string(key_decoded) == string(key) {
				return mpt.db[node.flag_value.value]
			} else if (len(key) <= len(key_decoded)) && string(key_decoded) != string(key) {
				return Node{}
			} else {
				return Node{}
			}

		//leaf node
		case 2, 3:
			//same key - return node
			//not same - return nil
			if reflect.DeepEqual(key_decoded, key) {
				return node
			} else {
				return Node{}
			}

		default:
			return node
		}
	//branch node
	case 1:
		br_path_value := key[0]
		//no hash in branch
		if len(key) == 0 {
			return node
		} else if node.branch_value[br_path_value] != "" {
			path := []uint8{}
			path = append(path, key[1:]...)
			new_hash_node := node.branch_value[br_path_value]
			n = mpt.db[new_hash_node]
			return mpt.GetByNode(n, path)
		} else {
			return Node{}
		}

	default:
		return Node{}
	}
	return Node{}
}

func (mpt *MerklePatriciaTrie) Delete(key string) (string, error) {
	key_hex := keyToHex(key)
	root_node := mpt.db[mpt.root]
	n := Node{}
	_, ok := mpt.kv[key]
	if ok {
		delete(mpt.kv, key)
	}

	//find leaf node with that key
	switch root_node.node_type {
	case 1:
		n = mpt.GetByNode(root_node, []byte(key_hex))
	case 2:
		root_key := root_node.flag_value.encoded_prefix
		root_key_decoded := compact_decode(root_key)
		if intersectionCount(string(root_key_decoded), string(key_hex)) > 0 {
			n = mpt.GetByNode(root_node, []byte(key_hex))
		} else {
			return "", errors.New("path_not_found")
		}
	}
	if n.node_type == 0 {
		return "", errors.New("path_not_found")
	}

	mpt.DeleteByNode(n, n.flag_value.value)
	return "", errors.New("")

}

func (mpt *MerklePatriciaTrie) DeleteByNode(node Node, search_value string) {

	key := mpt.KeyByValue(node) //hash value of node with value we need to delete
	parent_node_hash, parent_index := mpt.FindParentNode(node, key)
	parent_node := mpt.db[parent_node_hash]

	switch node.node_type {
	//node is branch
	case 1:
		//count number of values in branch
		count := 0
		index_deletion := 0

		for _, branch_values := range node.branch_value {
			if branch_values != "" {
				count += 1
			}
			if branch_values == search_value {
				index_deletion = 16
			}
		}
		search_value = node.branch_value[16]
		if index_deletion == 16 {
			node.branch_value[16] = ""
		}
		//find hash value of left leaf
		hash_last_leaf := ""
		index_count := -1
		index_lasf_leaf := 0
		for _, branch_values := range node.branch_value {
			index_count += 1
			if branch_values != "" {
				hash_last_leaf = branch_values
				index_lasf_leaf = index_count
			}
		}

		//if branch is root
		if (key == mpt.root) && (count-1 <= 1) {
			last_leaf := mpt.db[hash_last_leaf]
			prefix := last_leaf.flag_value.encoded_prefix
			prefix = compact_decode(prefix)
			prefix_new := []uint8{}
			prefix_new = append(prefix_new, uint8(index_lasf_leaf))
			prefix_new = append(prefix_new, prefix...)
			prefix_new = append(prefix_new, 16)
			last_leaf.flag_value.encoded_prefix = compact_encode(prefix_new)
			delete(mpt.db, key)
			mpt.db[hash_last_leaf] = last_leaf
			mpt.root = hash_last_leaf
		} else if count-1 > 1 {
			//branch has > 2 values
			//delete hash from branch
			node.branch_value[16] = ""
			mpt.db[key] = node
		} else {
			//branch left with 1 value
			child_node := mpt.db[hash_last_leaf]
			if (index_lasf_leaf == 16) && (parent_node.node_type == 2) {
				//value from branch and parent is ext
				prefix := parent_node.flag_value.encoded_prefix
				prefix = compact_decode(prefix)
				prefix = append(prefix, 16)
				parent_node.flag_value.encoded_prefix = compact_encode(prefix)
				parent_node.flag_value.value = search_value
				delete(mpt.db, key)
				mpt.db[parent_node_hash] = parent_node
			} else if (index_lasf_leaf == 16) && (parent_node.node_type == 1) {
				//value from branch and parent is branch
				delete(mpt.db, key)
				en_prefix := []uint8{}
				en_prefix = append(en_prefix, 16)
				f_leaf := Flag_value{encoded_prefix: compact_encode(en_prefix), value: search_value}
				new_leaf := Node{node_type: 2, branch_value: [17]string{}, flag_value: f_leaf}
				new_leaf_hash := new_leaf.hash_node()
				parent_node.branch_value[parent_index] = new_leaf_hash

				mpt.db[new_leaf_hash] = new_leaf
				mpt.db[parent_node_hash] = parent_node

			} else if (child_node.node_type == 1) && (parent_node.node_type == 1) {
				//parent node -> branch
				//child node (hash_last_leaf) -> branch
				node.node_type = 2
				node.branch_value = [17]string{}
				node.flag_value.value = hash_last_leaf
				node.flag_value.encoded_prefix = compact_encode([]uint8{uint8(index_lasf_leaf)})
				mpt.db[key] = node
			} else if (child_node.node_type == 1) && (parent_node.node_type == 2) {
				//parent node -> ext
				//child node (hash_last_leaf) -> branch
				ext_nibbles := compact_decode(parent_node.flag_value.encoded_prefix)
				ext_nibbles = append(ext_nibbles, uint8(index_lasf_leaf))
				parent_node.flag_value.value = hash_last_leaf
				parent_node.flag_value.encoded_prefix = compact_encode(ext_nibbles)
				delete(mpt.db, key)
				mpt.db[parent_node_hash] = parent_node

			} else if (child_node.node_type == 2) && (parent_node.node_type == 1) {
				//parent node -> branch
				//child node (hash_last_leaf) -> ext
				ext_nibbles := []uint8{}
				ext_nibbles = append(ext_nibbles, uint8(index_lasf_leaf))
				ext_nibbles = append(ext_nibbles, compact_decode(child_node.flag_value.encoded_prefix)...)
				prefix := child_node.flag_value.encoded_prefix
				if prefix[0]/16 > 1 {
					ext_nibbles = append(ext_nibbles, 16)
				}
				//parent_node.flag_value.value = hash_last_leaf
				parent_node.branch_value[parent_index] = hash_last_leaf
				child_node.flag_value.encoded_prefix = compact_encode(ext_nibbles)
				delete(mpt.db, key)
				mpt.db[hash_last_leaf] = child_node
				mpt.db[parent_node_hash] = parent_node

			} else {
				//parent node -> ext
				//child node (hash_last_leaf) -> ext
				ext_par_nibbles := compact_decode(parent_node.flag_value.encoded_prefix)
				ext_par_nibbles = append(ext_par_nibbles, uint8(index_lasf_leaf))
				ext_par_nibbles = append(ext_par_nibbles, compact_decode(child_node.flag_value.encoded_prefix)...)
				prefix := child_node.flag_value.encoded_prefix
				if prefix[0]/16 > 1 {
					ext_par_nibbles = append(ext_par_nibbles, 16)
				}

				parent_node.flag_value.encoded_prefix = compact_encode(ext_par_nibbles)
				parent_node.flag_value.value = child_node.flag_value.value
				delete(mpt.db, key)
				delete(mpt.db, hash_last_leaf)
				mpt.db[parent_node_hash] = parent_node
			}

		}

	//node is leaf
	case 2:
		//count number of values in branch
		count := 0
		for _, branch_values := range parent_node.branch_value {
			if branch_values != "" {
				count += 1
			}
		}
		//if leaf is root
		if key == mpt.root {
			delete(mpt.db, key)
			mpt.root = ""
		} else if count-1 > 1 {
			//more than 1 value left in branch
			//delete hash from branch
			parent_node.branch_value[parent_index] = ""
			mpt.db[parent_node_hash] = parent_node
			delete(mpt.db, key)
		} else {
			//only 1 value leaft in branch
			parent_node.branch_value[parent_index] = ""
			mpt.db[parent_node_hash] = parent_node
			delete(mpt.db, key)
			mpt.DeleteByNode(parent_node, search_value)
		}
	}
}

func compact_decode(hex_array []uint8) []uint8 {
	decoded_arr := []uint8{}
	for i := 0; i < len(hex_array); i += 1 {
		firstPart := hex_array[i] / 16
		secondPart := hex_array[i] % 16
		decoded_arr = append(decoded_arr, firstPart)
		decoded_arr = append(decoded_arr, secondPart)
	}

	if decoded_arr[0] == 0 || decoded_arr[0] == 2 {
		decoded_arr = append(decoded_arr[:0], decoded_arr[1:]...)
		decoded_arr = append(decoded_arr[:0], decoded_arr[1:]...)
	} else if decoded_arr[0] == 1 || decoded_arr[0] == 3 {
		decoded_arr = append(decoded_arr[:0], decoded_arr[1:]...)
	}
	return decoded_arr
}

// If Leaf, ignore 16 at the end
func compact_encode(encoded_arr []uint8) []uint8 {
	//encoded_arr = [] {1, 6, 1}
	term := 0
	if len(encoded_arr) == 0 {
		term = 0
	} else if encoded_arr[len(encoded_arr)-1] == 16 {
		term = 1
		encoded_arr = encoded_arr[:len(encoded_arr)-1]
	} else {
		term = 0
	}

	oddlen := len(encoded_arr) % 2
	flags := 2*term + oddlen

	if oddlen == 1 {
		encoded_arr = append([]uint8{uint8(flags)}, encoded_arr...)
	} else {
		encoded_arr = append([]uint8{0}, encoded_arr...)
		encoded_arr = append([]uint8{uint8(flags)}, encoded_arr...)
	}

	result := []uint8{}
	for i := 0; i < len(encoded_arr); i += 2 {
		result = append(result, (16*encoded_arr[i] + 1*encoded_arr[i+1]))
	}
	return result
}

func test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

const hextable = "0123456789abcdef"

//convert key to hex array
func keyToHex(key string) []uint8 {
	result := []uint8{}
	encodedByteString := hex.EncodeToString([]byte(key))
	for _, encodedByte := range encodedByteString {
		result = append(result, uint8(strings.IndexByte(hextable, uint8(encodedByte))))
	}
	return result
}

func intersectionCount(key string, key_node string) int {
	count := 0
	min := 0
	if len(key) > len(key_node) {
		min = len(key_node)
	} else {
		min = len(key)
	}
	for i := 0; i < min; i += 1 {
		if key[i] == key_node[i] {
			count += 1
		} else {
			break
		}
	}
	return count
}

func (mpt *MerklePatriciaTrie) KeyByValue(node Node) string {
	for key, value := range mpt.db {
		if reflect.DeepEqual(node, value) {
			return key
		}
	}
	return ""
}

func (mpt *MerklePatriciaTrie) FindParentNode(node Node, key string) (string, int) {
	index := 0
	for key_node, node_db := range mpt.db {
		//go through branch values
		index = -1
		for _, branch_values := range node_db.branch_value {
			index += 1
			if branch_values == key {
				return key_node, index
			}
		}
		//go through ext values
		if node_db.flag_value.value == key {
			return key_node, index
		}
	}
	return "", 0
}

//-----------------------
func (node *Node) hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = node.flag_value.value
	}

	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

func (node *Node) String() string {
	str := "empty string"
	switch node.node_type {
	case 0:
		str = "[Null Node]"
	case 1:
		str = "Branch["
		for i, v := range node.branch_value[:16] {
			str += fmt.Sprintf("%d=\"%s\", ", i, v)
		}
		str += fmt.Sprintf("value=%s]", node.branch_value[16])
	case 2:
		encoded_prefix := node.flag_value.encoded_prefix
		node_name := "Leaf"
		if is_ext_node(encoded_prefix) {
			node_name = "Ext"
		}
		ori_prefix := strings.Replace(fmt.Sprint(compact_decode(encoded_prefix)), " ", ", ", -1)
		str = fmt.Sprintf("%s<%v, value=\"%s\">", node_name, ori_prefix, node.flag_value.value)
	}
	return str
}

func node_to_string(node Node) string {
	return node.String()
}

func (mpt *MerklePatriciaTrie) GetKeyValue() map[string]string {
	return mpt.kv
}

func (mpt *MerklePatriciaTrie) GetSignature() map[string][]byte {
	return mpt.ks
}

func (mpt *MerklePatriciaTrie) GetRoot() string {
	return mpt.root
}

func (mpt *MerklePatriciaTrie) GetDB() map[string]Node {
	return mpt.db
}

func (mpt *MerklePatriciaTrie) Initial() {
	mpt.db = make(map[string]Node)
}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0]/16 < 2
}

func TestCompact() {
	test_compact_encode()
}

func (mpt *MerklePatriciaTrie) String() string {
	content := fmt.Sprintf("ROOT=%s\n", mpt.root)
	for hash := range mpt.db {
		content += fmt.Sprintf("%s: %s\n", hash, node_to_string(mpt.db[hash]))
	}
	return content
}

func (mpt *MerklePatriciaTrie) Order_nodes() string {
	raw_content := mpt.String()
	content := strings.Split(raw_content, "\n")
	root_hash := strings.Split(strings.Split(content[0], "HashStart")[1], "HashEnd")[0]
	queue := []string{root_hash}
	i := -1
	rs := ""
	cur_hash := ""
	for len(queue) != 0 {
		last_index := len(queue) - 1
		cur_hash, queue = queue[last_index], queue[:last_index]
		i += 1
		line := ""
		for _, each := range content {
			if strings.HasPrefix(each, "HashStart"+cur_hash+"HashEnd") {
				line = strings.Split(each, "HashEnd: ")[1]
				rs += each + "\n"
				rs = strings.Replace(rs, "HashStart"+cur_hash+"HashEnd", fmt.Sprintf("Hash%v", i), -1)
			}
		}
		temp2 := strings.Split(line, "HashStart")
		flag := true
		for _, each := range temp2 {
			if flag {
				flag = false
				continue
			}
			queue = append(queue, strings.Split(each, "HashEnd")[0])
		}
	}
	return rs
}

func (mpt *MerklePatriciaTrie) ChangeHash(old_hash string, new_hash string) {
	for key, value := range mpt.db {
		switch value.node_type {
		case 1:
			for i := 0; i < 17; i++ {
				if value.branch_value[i] == old_hash {
					value.branch_value[i] = new_hash
					delete(mpt.db, key)
					new_hash_node := value.hash_node()
					mpt.db[new_hash_node] = value

					if mpt.root == key {
						mpt.root = new_hash_node
					}
					mpt.ChangeHash(key, new_hash_node)
				}
			}
		case 2:
			if value.flag_value.value == old_hash {
				value.flag_value.value = new_hash
				delete(mpt.db, key)
				new_hash_node := value.hash_node()
				mpt.db[new_hash_node] = value
				if mpt.root == key {
					mpt.root = new_hash_node
				}
				mpt.ChangeHash(key, new_hash_node)
			}
		}
	}
}
