package replies

import "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/threads"

type target struct {
	reply  Reply
	thread threads.Thread
}

func (target *target) IsReply() bool {
	return target.reply != nil
}

func (target *target) Reply() Reply {
	return target.reply
}

func (target *target) IsThread() bool {
	return target.thread != nil
}

func (target *target) Thread() threads.Thread {
	return target.thread
}
