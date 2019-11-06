package consulStruct

type ConsulInfo struct {
	LockIndex   int
	Key         string
	Flags       int
	Value       string
	CreateIndex int
	ModifyIndex int
}
