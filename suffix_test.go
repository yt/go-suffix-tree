package suffix

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	RunAlhoc = flag.Bool("alhoc", false, "Run alhoc testing")
)

// go test -args -alhoc to enable alhoc testing
func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func getFixtures() ([]string, *Tree) {
	tree := NewTree()
	lists := []string{
		"edible", "presentable", "abominable", "credible",
		"picturesque", "statuesque", "nothing", "something", "thing", "nonsense",
		"random word", "word", "table", "unbelievable", "believable", "sense",
	}
	for _, s := range lists {
		tree.Insert([]byte(s), s)
	}
	return lists, tree
}

func TestInsertReturn(t *testing.T) {
	tree := NewTree()
	oldValue, ok := tree.Insert([]byte("sth"), "sth")
	assert.True(t, ok)
	assert.Nil(t, oldValue)

	oldValue, ok = tree.Insert([]byte("sth"), "else")
	assert.True(t, ok)
	assert.Equal(t, "sth", oldValue.(string))
}

func assertGet(t *testing.T, tree *Tree, expectedValue string, expectedFound bool) {
	value, found := tree.Get([]byte(expectedValue))
	assert.Equal(t, expectedFound, found)
	if expectedFound && value != nil {
		assert.Equal(t, expectedValue, value.(string))
	}
}

func TestGet_EmptyTree(t *testing.T) {
	tree := NewTree()
	assertGet(t, tree, "sth", false)
}

func TestGet_Base(t *testing.T) {
	tree := NewTree()
	tree.Insert([]byte("sth"), "sth")
	assertGet(t, tree, "sth", true)
	assertGet(t, tree, "else", false)
	assertGet(t, tree, "any sth", false)

	tree = NewTree()
	tree.Insert([]byte("sth else"), "sth else")
	tree.Insert([]byte("else"), "else")
	assertGet(t, tree, "else", true)
	assertGet(t, tree, "sth else", true)

	tree = NewTree()
	tree.Insert([]byte("else"), "else")
	tree.Insert([]byte("sth else"), "sth else")
	tree.Insert([]byte("any sth else"), "any sth else")
	tree.Insert([]byte("anything sth else"), "anything sth else")
	tree.Insert([]byte("at any sth else"), "at any sth else")
	assertGet(t, tree, "else", true)
	assertGet(t, tree, "sth else", true)
	assertGet(t, tree, "any sth else", true)
}

// GetPredecessor is mostly like Get, but please notice their slight differnces.
func assertGetPredecessor(t *testing.T, tree *Tree, key string, expectedFound bool) {
	matchedKey, value, found := tree.GetPredecessor([]byte(key))
	assert.Equal(t, expectedFound, found)
	if expectedFound && value != nil {
		assert.Equal(t, string(matchedKey), value.(string))
	}
}

func assertGetPredecessorCheckKey(t *testing.T, tree *Tree, key string,
	expectedKey string) {

	matchedKey, value, found := tree.GetPredecessor([]byte(key))
	assert.True(t, found, "expected %s, got nothing", key)
	if value != nil {
		assert.Equal(t, expectedKey, string(matchedKey))
		assert.Equal(t, expectedKey, value.(string))
	}
}

func TestGetPredecessor_EmptyTree(t *testing.T) {
	tree := NewTree()
	assertGetPredecessor(t, tree, "banana", false)
}

func TestGetPredecessor_Base(t *testing.T) {
	tree := NewTree()
	tree.Insert([]byte("sth"), "sth")
	assertGetPredecessor(t, tree, "th", false)
	assertGetPredecessor(t, tree, "else", false)
	assertGetPredecessorCheckKey(t, tree, "sth", "sth")
	assertGetPredecessorCheckKey(t, tree, "any sth", "sth")

	tree.Insert([]byte("else sth"), "else sth")
	assertGetPredecessor(t, tree, "sth", true)
	assertGetPredecessorCheckKey(t, tree, "lse sth", "sth")
	assertGetPredecessorCheckKey(t, tree, "any else sth", "else sth")

	tree.Insert([]byte("any sth"), "tenth")
	assertGetPredecessor(t, tree, "th", false)
	assertGetPredecessor(t, tree, "fourth", false)
}

func assertGetSuccessor(t *testing.T, tree *Tree, key string, expectedFound bool) {
	matchedKey, value, found := tree.GetSuccessor([]byte(key))
	assert.Equal(t, expectedFound, found)
	if expectedFound && value != nil {
		assert.Equal(t, string(matchedKey), value.(string))
	}
}

func assertGetSuccessorCheckKey(t *testing.T, tree *Tree, key string,
	expectedKey string) {

	matchedKey, value, found := tree.GetSuccessor([]byte(key))
	assert.True(t, found, "expected %s, got nothing", key)
	if value != nil {
		assert.Equal(t, expectedKey, string(matchedKey))
		assert.Equal(t, expectedKey, value.(string))
	}
}

func TestGetSuccessor_EmptyTree(t *testing.T) {
	tree := NewTree()
	assertGetSuccessor(t, tree, "apple", false)
}

func TestGetSuccessor_Base(t *testing.T) {
	tree := NewTree()
	tree.Insert([]byte("sth"), "sth")
	assertGetSuccessor(t, tree, "sth", true)
	assertGetSuccessorCheckKey(t, tree, "th", "sth")
	assertGetSuccessor(t, tree, "else sth", false)

	tree.Insert([]byte("any sth"), "any sth")
	tree.Insert([]byte("else any sth"), "else any sth")
	assertGetSuccessorCheckKey(t, tree, "y sth", "any sth")

	tree.Insert([]byte("goose any sth"), "goose any sth")
	assertGetSuccessor(t, tree, "geese any sth", false)
}

func TestRemove_EmptyTree(t *testing.T) {
	tree := NewTree()
	_, found := tree.Remove([]byte("anything"))
	assert.False(t, found)
}

func TestRemove_Base(t *testing.T) {
	tree := NewTree()
	tree.Insert([]byte("else"), "else")
	_, found := tree.Remove([]byte("lse"))
	assert.False(t, found)
	_, found = tree.Remove([]byte("anything"))
	assert.False(t, found)

	tree.Insert([]byte("sth else"), "sth else")
	oldValue, found := tree.Remove([]byte("sth else"))
	assert.True(t, found)
	assert.Equal(t, "sth else", oldValue.(string))
	assertGet(t, tree, "sth else", false)
	assertGet(t, tree, "else", true)

	tree.Remove([]byte("else"))
	assertGet(t, tree, "else", false)

	tree = NewTree()
	tree.Insert([]byte("sth else"), "sth else")
	tree.Insert([]byte("else"), "else")
	_, found = tree.Remove([]byte("else"))
	assert.True(t, found)
	assertGet(t, tree, "else", false)
	assertGet(t, tree, "sth else", true)
}

func TestWalk(t *testing.T) {
	lists, tree := getFixtures()

	result := map[string]string{}
	tree.Walk(func(key []byte, value interface{}) bool {
		result[string(key)] = value.(string)
		return false
	})
	for _, s := range lists {
		assert.Equal(t, s, result[s])
	}

	count := 0
	tree.Walk(func(key []byte, value interface{}) bool {
		if string(key) == "believable" {
			return true
		}
		count++
		return false
	})
	assert.Equal(t, 3, count)

	count = 0
	tree.Walk(func(key []byte, value interface{}) bool {
		if string(key) == "something" {
			return true
		}
		count++
		return false
	})
	assert.Equal(t, 15, count)
}

func TestWalkSuffix_EmptyTree(t *testing.T) {
	tree := NewTree()
	count := 0
	tree.WalkSuffix([]byte{}, func(key []byte, value interface{}) bool {
		count++
		return false
	})
	assert.Equal(t, 0, count)
}

func TestWalkSuffix_Base(t *testing.T) {
	lists, tree := getFixtures()

	count := len(lists)
	tree.WalkSuffix([]byte{}, func(key []byte, value interface{}) bool {
		count--
		return false
	})
	assert.Equal(t, 0, count)

	count = 0
	tree.WalkSuffix([]byte("nonexist"), func(key []byte, value interface{}) bool {
		count--
		return false
	})
	assert.Equal(t, 0, count)

	count = 5
	suffix := []byte("able")
	tree.WalkSuffix(suffix, func(key []byte, value interface{}) bool {
		if !bytes.HasSuffix(key, suffix) {
			assert.FailNowf(t, "The walked key %v should have given suffix %v",
				string(key), string(suffix))
		}
		count--
		return false
	})
	assert.Equal(t, 0, count)

	count = 2
	tree.WalkSuffix([]byte("word"), func(key []byte, value interface{}) bool {
		count--
		return false
	})
	assert.Equal(t, 0, count)

	count = 1
	suffix = []byte("redible")
	tree.WalkSuffix(suffix, func(key []byte, value interface{}) bool {
		count--
		return false
	})
	assert.Equal(t, 0, count)

	count = 0
	tree.WalkSuffix([]byte("anything"), func(key []byte, value interface{}) bool {
		count--
		return false
	})
	assert.Equal(t, 0, count)
}

func dumpTestData(wordRef map[string]bool, tree *Tree, ops []string, errMsg string) {
	opDumpFile, _ := ioutil.TempFile("", "suffix_test_op_dump_")
	defer opDumpFile.Close()
	for _, op := range ops {
		println(op)
		opDumpFile.Write(append([]byte(op), []byte("\n")...))
	}
	println("\nWord status:")
	words := []string{}
	for word, _ := range wordRef {
		words = append(words, word)
	}
	sort.Sort(sort.StringSlice(words))
	for _, word := range words {
		if wordRef[word] {
			println(word, "existed")
		} else {
			println(word, "removed")
		}
	}
	println("\nTree nodes:")
	nodeDumpFile, _ := ioutil.TempFile("", "suffix_test_node_dump_")
	defer nodeDumpFile.Close()
	tree.walkNode(func(labels [][]byte, value interface{}) {
		if labels[0] == nil {
			return
		}
		suffixes := []string{}
		for _, l := range labels {
			suffixes = append(suffixes, string(l))
		}
		println(strings.Join(suffixes, ":"))

		nodeDumpFile.Write(append(bytes.Join(labels, []byte(":")), []byte("\n")...))
	})
	println(errMsg)
	println("Also dump operation records to", opDumpFile.Name())
	println("Also dump final node status to", nodeDumpFile.Name())
}

func checkLabelOrder(tree *Tree) (string, bool) {
	var preLabelLen int
	var msg string
	inOrder := true
	tree.walkNode(func(labels [][]byte, value interface{}) {
		if labels[0] != nil {
			if len(labels[0]) < preLabelLen {
				msg = fmt.Sprintf("expect label len not shorter than %d, actual len(%s) is %d",
					preLabelLen, string(labels[0]), len(labels[0]))
				inOrder = false
			}
			preLabelLen = len(labels[0])
		} else {
			preLabelLen = 0
		}
	})
	return msg, inOrder
}

const banner = `
Start alhoc test.
Repeat below steps in 30 seconds.
1. Generate 256 random words, and insert them into a new Tree.
2. Perform 2048 random operations with pre-generated 256 words.
3. Dump the generated test data once failed.
`

func TestAlhoc(t *testing.T) {
	if !*RunAlhoc {
		t.SkipNow()
	}
	println(banner)
	OpNum := 2048
	WordNum := 256
	// Put some variable definitions outside for loop,
	// so that we could refer it in test dump.
	var wordRef map[string]bool
	var randomWords []string
	var ops []string
	rand.Seed(time.Now().UnixNano())
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	mismatchLetters := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	testTurns := 0
	testEnd := time.NewTimer(30 * time.Second)
	var errMsg string
	var tree *Tree
	defer func() {
		if r := recover(); r != nil {
			errMsg = fmt.Sprintf("Panic happened: %v", r)
			dumpTestData(wordRef, tree, ops, errMsg)
			t.FailNow()
		}
	}()

	for {
		select {
		case <-testEnd.C:
			println(testTurns, "turns of alhoc tests pass.")
			println("Alhoc test is finished.")
			return
		default:
		}
		tree = NewTree()
		wordRef = map[string]bool{}
		randomWords = []string{}
		ops = []string{}
		for wordCount := 0; wordCount < WordNum; {
			b := make([]byte, rand.Intn(12))
			for i := range b {
				b[i] = letters[rand.Intn(len(letters))]
			}
			bs := string(b)
			if _, ok := wordRef[bs]; !ok {
				wordRef[bs] = true
				wordCount += 1
			}
			ops = append(ops, "Insert\t"+bs)
			tree.Insert(b, bs)
			value, found := tree.Get(b)
			if !found {
				errMsg = fmt.Sprintf("expect get %v after insertion, actual not found", bs)
			} else {
				if value.(string) != bs {
					errMsg = fmt.Sprintf("expect insert %v, actual %v", bs, value.(string))
					goto failed
				}
			}
			msg, inOrder := checkLabelOrder(tree)
			if !inOrder {
				errMsg = msg
				goto failed
			}
		}
		for s := range wordRef {
			randomWords = append(randomWords, s)
		}
		for i := 0; i < OpNum; i++ {
			if OpNum%256 == 0 {
				existedNum := 0
				for _, existed := range wordRef {
					if existed {
						existedNum += 1
					}
				}
				if tree.Len() != existedNum {
					errMsg = fmt.Sprintf("expect len %v, actual %v", existedNum, tree.Len())
					goto failed
				}
			}
			word := randomWords[rand.Intn(WordNum)]
			existed := wordRef[word]
			switch rand.Intn(6) {
			case 0:
				ops = append(ops, "Get\t"+word)
				value, found := tree.Get([]byte(word))
				if found {
					if !existed {
						errMsg = fmt.Sprintf("expect not found %v, actual found", word)
						goto failed
					}
					if value.(string) != word {
						errMsg = fmt.Sprintf("expect get %v, actual %v", word, value.(string))
						goto failed
					}
				} else if existed {
					errMsg = fmt.Sprintf("expect found %v, actual not found", word)
					goto failed
				}
			case 1:
				ops = append(ops, "Insert\t"+word)
				value, _ := tree.Insert([]byte(word), word)
				if existed {
					if value.(string) != word {
						errMsg = fmt.Sprintf("expect insert %v, actual %v", word, value.(string))
						goto failed
					}
				}
				wordRef[word] = true
				_, found := tree.Get([]byte(word))
				if !found {
					errMsg = fmt.Sprintf("expect get %v after insertion, actual not found", word)
				}
				msg, inOrder := checkLabelOrder(tree)
				if !inOrder {
					errMsg = msg
					goto failed
				}
			case 2:
				ops = append(ops, "Remove\t"+word)
				value, found := tree.Remove([]byte(word))
				wordRef[word] = false
				if found {
					if !existed {
						errMsg = fmt.Sprintf("expect not found %v in removal, actual found", word)
						goto failed
					}
					if value.(string) != word {
						errMsg = fmt.Sprintf("expect remove %v, actual %v", word, value.(string))
						goto failed
					}
					_, found = tree.Get([]byte(word))
					if found {
						errMsg = fmt.Sprintf("expect %v not found after removal, actual found", word)
					}
					msg, inOrder := checkLabelOrder(tree)
					if !inOrder {
						errMsg = msg
						goto failed
					}
				} else if existed {
					errMsg = fmt.Sprintf("expect found %v in removal, actual not found", word)
					goto failed
				}
			case 3:
				mismatchLabel := make([]byte, rand.Intn(3))
				for i := range mismatchLabel {
					mismatchLabel[i] = mismatchLetters[rand.Intn(len(mismatchLetters))]
				}
				suffix := string(mismatchLabel) + word
				ops = append(ops, "GetPredecessor\t"+suffix)
				key, value, found := tree.GetPredecessor([]byte(suffix))
				if existed {
					if !found {
						errMsg = fmt.Sprintf("expect getPredecessor found %v with %v, actual not found",
							word, suffix)
						goto failed
					}
					if value.(string) != string(key) {
						errMsg = fmt.Sprintf(
							"expect getPredecessor %v, actual %v", string(key), value.(string))
						goto failed
					}
				} else {
					if found {
						if !strings.HasSuffix(suffix, string(key)) {
							errMsg = fmt.Sprintf(
								"expect getPredecessor found suffix matched %v, actual found %v",
								suffix, string(key))
							goto failed
						}
					}
				}
			case 4:
				suffix := word[:len(word)]
				bsuffix := []byte(suffix)
				ops = append(ops, "GetSuccessor\t"+suffix)
				matchedKey, matchedValue, found := tree.GetSuccessor(bsuffix)
				if !found {
					if existed {
						errMsg = fmt.Sprintf("expect getSuccessor found %v with %v, actual not found",
							word, suffix)
						goto failed
					}
				} else {
					matched := false
					// Walk method walks through the tree as the same way as GetSuccessor
					tree.Walk(func(key []byte, value interface{}) bool {
						if bytes.HasSuffix(key, bsuffix) {
							word = string(key)
							matched = bytes.Equal(matchedKey, key) &&
								matchedValue.(string) == matchedValue.(string)
							return true
						}
						return false
					})
					if !matched {
						errMsg = fmt.Sprintf("expect getSuccessor found %v with %v, actual found %v",
							word, suffix, string(matchedKey))
						goto failed
					}
				}
			case 5:
				var suffix string
				if len(word) > 0 {
					suffix = word[rand.Intn(len(word)):]
				}
				bsuffix := []byte(suffix)
				ops = append(ops, "WalkSuffix\t"+suffix)
				shouldMatchKeys := []string{}
				matchedKeys := map[string]bool{}
				tree.Walk(func(key []byte, value interface{}) bool {
					if bytes.HasSuffix(key, bsuffix) {
						shouldMatchKeys = append(shouldMatchKeys, string(key))
					}
					return false
				})
				tree.WalkSuffix(bsuffix, func(key []byte, value interface{}) bool {
					matchedKeys[string(key)] = true
					return false
				})
				if len(shouldMatchKeys) != len(matchedKeys) {
					errMsg = fmt.Sprintf("expect walkSuffix with %v matches %v keys, actual %v",
						suffix, len(shouldMatchKeys), len(matchedKeys))
					goto failed
				}
				for _, s := range shouldMatchKeys {
					if _, ok := matchedKeys[s]; !ok {
						errMsg = fmt.Sprintf("expect walkSuffix with %v travels %v",
							suffix, s)
						goto failed
					}
				}
			}
		}
		testTurns++
	}
failed:
	dumpTestData(wordRef, tree, ops, errMsg)
	t.FailNow()
}
