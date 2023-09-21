package spmigration

func (t *Target_Minions) Update_Target_Minion_Status(analyase_minions *Target_Minions) {
	for _, minion := range t.Minion_List {
		for _, minion2 := range analyase_minions.Minion_List {
			if minion.Minion_ID == minion2.Minion_ID {
				minion.Host_Job_Info = minion2.Host_Job_Info
				minion.Migration_Stage = minion2.Migration_Stage
				minion.Migration_Stage_Status = minion2.Migration_Stage_Status
			}
		}
	}
}
