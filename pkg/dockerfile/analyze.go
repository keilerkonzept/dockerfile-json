package dockerfile

func (d *Dockerfile) analyzeStages() {
	seenStageNames := make(map[string]int)
	for i, stage := range d.Stages {
		stageIndex, stageIndexOK := seenStageNames[stage.BaseName]
		switch {
		case stageIndexOK:
			stage.From.Stage = &FromStage{
				Named: &stage.BaseName,
				Index: stageIndex,
			}
		case stage.BaseName == "scratch":
			stage.From.Scratch = true
			stage.From.Image = nil
		default:
			stage.From.Image = &stage.BaseName
		}
		if stage.Stage.Name != "" {
			stage.Name = &stage.Stage.Name
			seenStageNames[stage.Stage.Name] = i
		}
		for i, command := range stage.Commands {
			stage.Commands[i].Name = command.Command.Name()
		}
	}
	return
}
