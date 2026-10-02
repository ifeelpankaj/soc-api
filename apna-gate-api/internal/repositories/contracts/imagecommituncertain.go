package contracts

type ImageCommitUncertain struct{}

func (*ImageCommitUncertain) Error() string { return "image commit outcome uncertain" }
