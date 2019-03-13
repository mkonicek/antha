package liquidhandling

type RobotInstructionVisitor interface {
	Transfer(*TransferInstruction)
	TransferBlock(*TransferBlockInstruction)
	ChannelBlock(*ChannelBlockInstruction)
	ChannelTransfer(*ChannelTransferInstruction)
	StateChange(*StateChangeInstruction)
	LoadTipsMove(*LoadTipsMoveInstruction)
	UnloadTipsMove(*UnloadTipsMoveInstruction)
	Reset(*ResetInstruction)
	ChangeAdaptor(*ChangeAdaptorInstruction)
	Aspirate(*AspirateInstruction)
	Dispense(*DispenseInstruction)
	Blowout(*BlowoutInstruction)
	PTZ(*PTZInstruction)
	Move(*MoveInstruction)
	MoveRaw(*MoveRawInstruction)
	LoadTips(*LoadTipsInstruction)
	UnloadTips(*UnloadTipsInstruction)
	Suck(*SuckInstruction)
	Blow(*BlowInstruction)
	SetPipetteSpeed(*SetPipetteSpeedInstruction)
	SetDriveSpeed(*SetDriveSpeedInstruction)
	Initialize(*InitializeInstruction)
	Finalize(*FinalizeInstruction)
	Wait(*WaitInstruction)
	LightsOn(*LightsOnInstruction)
	LightsOff(*LightsOffInstruction)
	Open(*OpenInstruction)
	Close(*CloseInstruction)
	LoadAdaptor(*LoadAdaptorInstruction)
	UnloadAdaptor(*UnloadAdaptorInstruction)
	MoveMix(*MoveMixInstruction)
	Mix(*MixInstruction)
	Message(*MessageInstruction)
	MovAsp(*MovAsp)
	MovDsp(*MovDsp)
	MovMix(*MovMix)
	MovBlo(*MovBlo)
	RemoveAllPlates(*RemoveAllPlatesInstruction)
	AddPlateTo(*AddPlateToInstruction)
	SplitBlock(*SplitBlockInstruction)
}

type RobotInstructionBaseVisitor struct {
	HandleTransfer        func(*TransferInstruction)
	HandleTransferBlock   func(*TransferBlockInstruction)
	HandleChannelBlock    func(*ChannelBlockInstruction)
	HandleChannelTransfer func(*ChannelTransferInstruction)
	HandleStateChange     func(*StateChangeInstruction)
	HandleLoadTipsMove    func(*LoadTipsMoveInstruction)
	HandleUnloadTipsMove  func(*UnloadTipsMoveInstruction)
	HandleReset           func(*ResetInstruction)
	HandleChangeAdaptor   func(*ChangeAdaptorInstruction)
	HandleAspirate        func(*AspirateInstruction)
	HandleDispense        func(*DispenseInstruction)
	HandleBlowout         func(*BlowoutInstruction)
	HandlePTZ             func(*PTZInstruction)
	HandleMove            func(*MoveInstruction)
	HandleMoveRaw         func(*MoveRawInstruction)
	HandleLoadTips        func(*LoadTipsInstruction)
	HandleUnloadTips      func(*UnloadTipsInstruction)
	HandleSuck            func(*SuckInstruction)
	HandleBlow            func(*BlowInstruction)
	HandleSetPipetteSpeed func(*SetPipetteSpeedInstruction)
	HandleSetDriveSpeed   func(*SetDriveSpeedInstruction)
	HandleInitialize      func(*InitializeInstruction)
	HandleFinalize        func(*FinalizeInstruction)
	HandleWait            func(*WaitInstruction)
	HandleLightsOn        func(*LightsOnInstruction)
	HandleLightsOff       func(*LightsOffInstruction)
	HandleOpen            func(*OpenInstruction)
	HandleClose           func(*CloseInstruction)
	HandleLoadAdaptor     func(*LoadAdaptorInstruction)
	HandleUnloadAdaptor   func(*UnloadAdaptorInstruction)
	HandleMoveMix         func(*MoveMixInstruction)
	HandleMix             func(*MixInstruction)
	HandleMessage         func(*MessageInstruction)
	HandleMovAsp          func(*MovAsp)
	HandleMovDsp          func(*MovDsp)
	HandleMovMix          func(*MovMix)
	HandleMovBlo          func(*MovBlo)
	HandleRemoveAllPlates func(*RemoveAllPlatesInstruction)
	HandleAddPlateTo      func(*AddPlateToInstruction)
	HandleSplitBlock      func(*SplitBlockInstruction)
}

func (self RobotInstructionBaseVisitor) Transfer(ins *TransferInstruction) {
	if self.HandleTransfer != nil {
		self.HandleTransfer(ins)
	}
}
func (self RobotInstructionBaseVisitor) TransferBlock(ins *TransferBlockInstruction) {
	if self.HandleTransferBlock != nil {
		self.HandleTransferBlock(ins)
	}
}
func (self RobotInstructionBaseVisitor) ChannelBlock(ins *ChannelBlockInstruction) {
	if self.HandleChannelBlock != nil {
		self.HandleChannelBlock(ins)
	}
}
func (self RobotInstructionBaseVisitor) ChannelTransfer(ins *ChannelTransferInstruction) {
	if self.HandleChannelTransfer != nil {
		self.HandleChannelTransfer(ins)
	}
}
func (self RobotInstructionBaseVisitor) StateChange(ins *StateChangeInstruction) {
	if self.HandleStateChange != nil {
		self.HandleStateChange(ins)
	}
}
func (self RobotInstructionBaseVisitor) LoadTipsMove(ins *LoadTipsMoveInstruction) {
	if self.HandleLoadTipsMove != nil {
		self.HandleLoadTipsMove(ins)
	}
}
func (self RobotInstructionBaseVisitor) UnloadTipsMove(ins *UnloadTipsMoveInstruction) {
	if self.HandleUnloadTipsMove != nil {
		self.HandleUnloadTipsMove(ins)
	}
}
func (self RobotInstructionBaseVisitor) Reset(ins *ResetInstruction) {
	if self.HandleReset != nil {
		self.HandleReset(ins)
	}
}
func (self RobotInstructionBaseVisitor) ChangeAdaptor(ins *ChangeAdaptorInstruction) {
	if self.HandleChangeAdaptor != nil {
		self.HandleChangeAdaptor(ins)
	}
}
func (self RobotInstructionBaseVisitor) Aspirate(ins *AspirateInstruction) {
	if self.HandleAspirate != nil {
		self.HandleAspirate(ins)
	}
}
func (self RobotInstructionBaseVisitor) Dispense(ins *DispenseInstruction) {
	if self.HandleDispense != nil {
		self.HandleDispense(ins)
	}
}
func (self RobotInstructionBaseVisitor) Blowout(ins *BlowoutInstruction) {
	if self.HandleBlowout != nil {
		self.HandleBlowout(ins)
	}
}
func (self RobotInstructionBaseVisitor) PTZ(ins *PTZInstruction) {
	if self.HandlePTZ != nil {
		self.HandlePTZ(ins)
	}
}
func (self RobotInstructionBaseVisitor) Move(ins *MoveInstruction) {
	if self.HandleMove != nil {
		self.HandleMove(ins)
	}
}
func (self RobotInstructionBaseVisitor) MoveRaw(ins *MoveRawInstruction) {
	if self.HandleMoveRaw != nil {
		self.HandleMoveRaw(ins)
	}
}
func (self RobotInstructionBaseVisitor) LoadTips(ins *LoadTipsInstruction) {
	if self.HandleLoadTips != nil {
		self.HandleLoadTips(ins)
	}
}
func (self RobotInstructionBaseVisitor) UnloadTips(ins *UnloadTipsInstruction) {
	if self.HandleUnloadTips != nil {
		self.HandleUnloadTips(ins)
	}
}
func (self RobotInstructionBaseVisitor) Suck(ins *SuckInstruction) {
	if self.HandleSuck != nil {
		self.HandleSuck(ins)
	}
}
func (self RobotInstructionBaseVisitor) Blow(ins *BlowInstruction) {
	if self.HandleBlow != nil {
		self.HandleBlow(ins)
	}
}
func (self RobotInstructionBaseVisitor) SetPipetteSpeed(ins *SetPipetteSpeedInstruction) {
	if self.HandleSetPipetteSpeed != nil {
		self.HandleSetPipetteSpeed(ins)
	}
}
func (self RobotInstructionBaseVisitor) SetDriveSpeed(ins *SetDriveSpeedInstruction) {
	if self.HandleSetDriveSpeed != nil {
		self.HandleSetDriveSpeed(ins)
	}
}
func (self RobotInstructionBaseVisitor) Initialize(ins *InitializeInstruction) {
	if self.HandleInitialize != nil {
		self.HandleInitialize(ins)
	}
}
func (self RobotInstructionBaseVisitor) Finalize(ins *FinalizeInstruction) {
	if self.HandleFinalize != nil {
		self.HandleFinalize(ins)
	}
}
func (self RobotInstructionBaseVisitor) Wait(ins *WaitInstruction) {
	if self.HandleWait != nil {
		self.HandleWait(ins)
	}
}
func (self RobotInstructionBaseVisitor) LightsOn(ins *LightsOnInstruction) {
	if self.HandleLightsOn != nil {
		self.HandleLightsOn(ins)
	}
}
func (self RobotInstructionBaseVisitor) LightsOff(ins *LightsOffInstruction) {
	if self.HandleLightsOff != nil {
		self.HandleLightsOff(ins)
	}
}
func (self RobotInstructionBaseVisitor) Open(ins *OpenInstruction) {
	if self.HandleOpen != nil {
		self.HandleOpen(ins)
	}
}
func (self RobotInstructionBaseVisitor) Close(ins *CloseInstruction) {
	if self.HandleClose != nil {
		self.HandleClose(ins)
	}
}
func (self RobotInstructionBaseVisitor) LoadAdaptor(ins *LoadAdaptorInstruction) {
	if self.HandleLoadAdaptor != nil {
		self.HandleLoadAdaptor(ins)
	}
}
func (self RobotInstructionBaseVisitor) UnloadAdaptor(ins *UnloadAdaptorInstruction) {
	if self.HandleUnloadAdaptor != nil {
		self.HandleUnloadAdaptor(ins)
	}
}
func (self RobotInstructionBaseVisitor) MoveMix(ins *MoveMixInstruction) {
	if self.HandleMoveMix != nil {
		self.HandleMoveMix(ins)
	}
}
func (self RobotInstructionBaseVisitor) Mix(ins *MixInstruction) {
	if self.HandleMix != nil {
		self.HandleMix(ins)
	}
}
func (self RobotInstructionBaseVisitor) Message(ins *MessageInstruction) {
	if self.HandleMessage != nil {
		self.HandleMessage(ins)
	}
}
func (self RobotInstructionBaseVisitor) MovAsp(ins *MovAsp) {
	if self.HandleMovAsp != nil {
		self.HandleMovAsp(ins)
	}
}
func (self RobotInstructionBaseVisitor) MovDsp(ins *MovDsp) {
	if self.HandleMovDsp != nil {
		self.HandleMovDsp(ins)
	}
}
func (self RobotInstructionBaseVisitor) MovMix(ins *MovMix) {
	if self.HandleMovMix != nil {
		self.HandleMovMix(ins)
	}
}
func (self RobotInstructionBaseVisitor) MovBlo(ins *MovBlo) {
	if self.HandleMovBlo != nil {
		self.HandleMovBlo(ins)
	}
}
func (self RobotInstructionBaseVisitor) RemoveAllPlates(ins *RemoveAllPlatesInstruction) {
	if self.HandleRemoveAllPlates != nil {
		self.HandleRemoveAllPlates(ins)
	}
}
func (self RobotInstructionBaseVisitor) AddPlateTo(ins *AddPlateToInstruction) {
	if self.HandleAddPlateTo != nil {
		self.HandleAddPlateTo(ins)
	}
}
func (self RobotInstructionBaseVisitor) SplitBlock(ins *SplitBlockInstruction) {
	if self.HandleSplitBlock != nil {
		self.HandleSplitBlock(ins)
	}
}
