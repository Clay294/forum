package threads

import "time"

func NewReqCreateThread() *ReqCreateThreads {
	return &ReqCreateThreads{
		Tags:                 make([]string, 0, 128),
		Status:               StatusDraft,
		ReqCreateThreadsMeta: &ReqCreateThreadsMeta{},
	}
}
func NewThread(rct *ReqCreateThreads) *Thread {
	return &Thread{
		ThreadBase: (*ThreadBase)(rct),
		ThreadMeta: NewThreadMeta(),
	}
}
func NewThreadMeta() *ThreadMeta {
	now := time.Now().UnixMilli()
	return &ThreadMeta{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewReqSearchByMainHome() *ReqSearchByMainHome {
	return &ReqSearchByMainHome{
		ReqSearchMeta: NewReqSearchMeta(),
	}
}
