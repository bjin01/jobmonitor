package pkg_updates

func Find_Next_Stage(wf []Workflow_Step, minion Minion_Data) string {
	for _, step := range wf {
		/* fmt.Printf("step: %+v, current_stage: %s\n", step.Name, minion.Migration_Stage)
		fmt.Printf("step: %+v, current_stage status: %s\n", step.Name, minion.Migration_Stage_Status) */
		if minion.Migration_Stage == "" && minion.JobID == 0 {
			for _, w := range wf {
				if w.Order == 1 {
					logger.Infof("Set minion %s to next stage: %s\n", minion.Minion_Name, w.Name)
					return w.Name // return the first stage
				}
			}
		}

		if step.Name == minion.Migration_Stage && minion.Migration_Stage_Status == "completed" {
			if step.Order < len(wf)-1 {
				step.Order = step.Order + 1
				for _, w := range wf {
					if step.Order == w.Order {
						logger.Infof("Set minion %s to next stage: %s\n", minion.Minion_Name, w.Name)
						return w.Name // return the next stage
					}
				}
			}

			if step.Order == len(wf)-1 {
				for _, w := range wf {
					if step.Order == w.Order {
						logger.Infof("Set minion %s to next stage: %s\n", minion.Minion_Name, w.Name)
						return w.Name // return the last stage
					}
				}
			}
		}
	}
	return ""
}
