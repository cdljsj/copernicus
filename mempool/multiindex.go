package mempool

import (
	"sort"

	"github.com/btcboost/copernicus/utils"
)

const (
	DescendantScore = iota
	MiningScore
	AncestorScore
	TimeSort
)

// MultiIndex the struct for support memPool store node, to implement MultiIndex sort
type MultiIndex struct {
	PoolNode              map[utils.Hash]*TxMempoolEntry // unique
	nodeKey               []*TxMempoolEntry
	byDescendantScoreSort []*TxMempoolEntry // ordered_non_unique;
	byEntryTimeSort       []*TxMempoolEntry // ordered_non_unique;
	byScoreSort           []*TxMempoolEntry // ordered_unique;
	byAncestorFeeSort     []*TxMempoolEntry // ordered_non_unique;
}

func NewMultiIndex() *MultiIndex {
	multi := MultiIndex{}
	multi.PoolNode = make(map[utils.Hash]*TxMempoolEntry)
	multi.byDescendantScoreSort = make([]*TxMempoolEntry, 0)
	multi.byEntryTimeSort = make([]*TxMempoolEntry, 0)
	multi.byScoreSort = make([]*TxMempoolEntry, 0)
	multi.byAncestorFeeSort = make([]*TxMempoolEntry, 0)
	multi.nodeKey = make([]*TxMempoolEntry, 0)

	return &multi
}

// AddElement add the element to the multiIndex; the element must meet multiIndex's keys various criterions;
func (multiIndex *MultiIndex) AddElement(hash utils.Hash, txEntry *TxMempoolEntry) {
	if _, has := multiIndex.PoolNode[hash]; has {
		return
	}
	multiIndex.PoolNode[hash] = txEntry
	multiIndex.nodeKey = append(multiIndex.nodeKey, txEntry)
}

// DelEntryByHash : delete the key correspond value In multiIndex;
func (multiIndex *MultiIndex) DelEntryByHash(hash utils.Hash) {
	if _, ok := multiIndex.PoolNode[hash]; ok {
		delete(multiIndex.PoolNode, hash)
		for i, v := range multiIndex.nodeKey {
			oriHash := v.TxRef.Hash
			if (&oriHash).IsEqual(&hash) {
				multiIndex.nodeKey = append(multiIndex.nodeKey[:i], multiIndex.nodeKey[i+1:]...)
				break
			}
		}
	}
}

// GetEntryByHash : return the key correspond value In multiIndex;
// And modify The return value will be Influence the multiIndex;
func (multiIndex *MultiIndex) GetEntryByHash(hash utils.Hash) *TxMempoolEntry {
	if v, ok := multiIndex.PoolNode[hash]; ok {
		return v
	}
	return nil
}

func (multiIndex *MultiIndex) Size() int {
	return len(multiIndex.PoolNode)
}

// GetByDescendantScoreSort : return the sort slice by descendantScore
func (multiIndex *MultiIndex) GetByDescendantScoreSort() []*TxMempoolEntry {
	multiIndex.updateSort(DescendantScore)
	return multiIndex.byDescendantScoreSort
}

func (multiIndex *MultiIndex) GetByDescendantScoreSortBegin() interface{} {
	multiIndex.updateSort(DescendantScore)
	if len(multiIndex.byDescendantScoreSort) > 0 {
		return multiIndex.byDescendantScoreSort[0]
	}
	return nil
}

func (multiIndex *MultiIndex) updateSort(flag int) {
	switch flag {
	case DescendantScore:
		multiIndex.byDescendantScoreSort = make([]*TxMempoolEntry, len(multiIndex.nodeKey))
		copy(multiIndex.byDescendantScoreSort, multiIndex.nodeKey)
		sort.SliceStable(multiIndex.byDescendantScoreSort, func(i, j int) bool {
			return CompareTxMemPoolEntryByDescendantScore(multiIndex.byDescendantScoreSort[i], multiIndex.byDescendantScoreSort[j])
		})
	case AncestorScore:
		multiIndex.byAncestorFeeSort = make([]*TxMempoolEntry, len(multiIndex.nodeKey))
		copy(multiIndex.byAncestorFeeSort, multiIndex.nodeKey)
		sort.SliceStable(multiIndex.byAncestorFeeSort, func(i, j int) bool {
			return CompareTxMemPoolEntryByAncestorFee(multiIndex.byAncestorFeeSort[i], multiIndex.byAncestorFeeSort[j])
		})
	case MiningScore:
		multiIndex.byScoreSort = make([]*TxMempoolEntry, len(multiIndex.nodeKey))
		copy(multiIndex.byScoreSort, multiIndex.nodeKey)
		sort.SliceStable(multiIndex.byScoreSort, func(i, j int) bool {
			return CompareTxMempoolEntryByScore(multiIndex.byScoreSort[i], multiIndex.byScoreSort[j])
		})
	case TimeSort:
		multiIndex.byEntryTimeSort = make([]*TxMempoolEntry, len(multiIndex.nodeKey))
		copy(multiIndex.byEntryTimeSort, multiIndex.nodeKey)
		sort.SliceStable(multiIndex.byEntryTimeSort, func(i, j int) bool {
			return CompareTxMemPoolEntryByEntryTime(multiIndex.byEntryTimeSort[i], multiIndex.byEntryTimeSort[j])
		})

	}
}

func (multiIndex *MultiIndex) GetbyEntryTimeSort() []*TxMempoolEntry {
	multiIndex.updateSort(TimeSort)
	return multiIndex.byEntryTimeSort
}

func (multiIndex *MultiIndex) GetbyEntryTimeSortBegin() interface{} {
	multiIndex.updateSort(TimeSort)
	if len(multiIndex.byEntryTimeSort) > 0 {
		return multiIndex.byEntryTimeSort[len(multiIndex.byEntryTimeSort)-1]
	}
	return nil
}

func (multiIndex *MultiIndex) GetbyAncestorFeeSort() []*TxMempoolEntry {
	multiIndex.updateSort(AncestorScore)
	return multiIndex.byAncestorFeeSort
}

func (multiIndex *MultiIndex) GetbyAncestorFeeSortBegin() interface{} {
	multiIndex.updateSort(AncestorScore)
	if len(multiIndex.byAncestorFeeSort) > 0 {
		return multiIndex.byAncestorFeeSort[len(multiIndex.byAncestorFeeSort)-1]
	}
	return nil
}

func (multiIndex *MultiIndex) GetbyScoreSort() []*TxMempoolEntry {
	multiIndex.updateSort(MiningScore)
	return multiIndex.byScoreSort
}

func (multiIndex *MultiIndex) GetbyScoreSortBegin() interface{} {
	multiIndex.updateSort(MiningScore)
	if len(multiIndex.byScoreSort) > 0 {
		return multiIndex.byScoreSort[len(multiIndex.byScoreSort)-1]
	}
	return nil
}
