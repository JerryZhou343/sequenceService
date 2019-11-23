package repository

import "github.com/mfslog/sequenceService/domain/segment/entity"

type SegmentSeqRepository interface {
	Get(pid, bid int32) (*entity.SegmentSeq, error)
	Save(data *entity.SegmentSeq) error
	List(page, pageSize int32) ([]*entity.SegmentSeq, int, error)
}
