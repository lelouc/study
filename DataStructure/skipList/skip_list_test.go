package skipList

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const (
	MAX_LENGTH = 20
	MAX_SCORE  = 100
)

type Value struct {
	Key   uint64
	Score uint64
}

type Cmp struct {
}

func (this *Cmp) CmpScore(v1, v2 interface{}) int {
	s1 := v1.(*Value).Score
	s2 := v2.(*Value).Score
	switch {
	case s1 < s2:
		return -1
	case s1 > s2:
		return 1
	default:
		return 0
	}
}

func (this *Cmp) CmpKey(v1, v2 interface{}) int {
	k1 := v1.(*Value).Key
	k2 := v2.(*Value).Key

	switch {
	case k1 < k2:
		return -1
	case k1 > k2:
		return 1
	default:
		return 0
	}
}

func PrintRankList(sk *SkipList) {
	fmt.Println("//////////ranklist/////////////")
	fmt.Println("rank\tscore\tkey\t")
	rank := 0
	head := sk.Head()
	node := head
	for node != nil {
		if node != head {
			rank++
			val, ok := node.Value().(*Value)
			if !ok {
				fmt.Println("sssss")
			}

			fmt.Printf("%d\t%d\t%d\t", rank, val.Score, val.Key)
		}
		node = node.Next()
	}
}

func RandomDelete(sk *SkipList) bool {
	delRank := uint32(rand.Intn(MAX_LENGTH) + 1)
	delValue := sk.GetNodeByRank(delRank).Value()
	delKey := delValue.(*Value).Key
	delScore := delValue.(*Value).Score
	fmt.Println("////////////delete///////////")
	fmt.Println("rank\tscore\tkey\t")
	fmt.Printf("%d\t%d\t%d\t", delRank, delScore, delKey)
	return sk.Delete(delValue)
}

func PrintRankNode(rank uint32, sk *SkipList) {
	value := sk.GetNodeByRank(rank).Value().(*Value)
	fmt.Println("///////////get value/////////")
	fmt.Println("rank\tscore\tkey\t")
	fmt.Printf("%d\t%d\t%d\t", rank, value.Score, value.Key)
}

func TestExample(t *testing.T) {
	rand.Seed(int64(time.Now().Unix()))
	sk := NewSkipList(&Cmp{})
	for i := 0; i < MAX_LENGTH; i++ {
		sk.Insert(&Value{
			Key:   uint64(i + 1),
			Score: uint64(rand.Intn(MAX_SCORE) + 1),
		})
	}
	PrintRankList(sk)
	if RandomDelete(sk) {
		PrintRankList(sk)
	}
	PrintRankNode(6, sk)
}
