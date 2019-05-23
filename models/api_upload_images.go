package models

type UploadImagesBody struct {
	UploadImagesRequest UploadImagesRequest `xml:"q0:UploadImagesRequest"`
}

type UploadImagesRequest struct {
	Request
	Images []Image `xml:"q0:Images"`
}

type UploadImagesResponseEnvelope struct {
	Reply UploadImagesReply `xml:"Body>UploadImagesReply"`
}

func (u *UploadImagesResponseEnvelope) Error() error {
	return u.Reply.Error()
}

// UploadImagesReply : UploadImages reply root (`xml:"Body>UploadImagesReply"`)
type UploadImagesReply struct {
	Reply
}
