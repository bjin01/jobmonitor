package spmigration

func (m *Target_Minions) Add_Online_Minions(list []Minion_Data) {

	unique_minions := []Minion_Data{}

	for _, minion := range list {
		exists := false
		// Check if the Minion_ID is already in the map
		for _, minion2 := range m.Minion_List {
			if minion.Minion_ID == minion2.Minion_ID {
				// If it is, don't add it again
				exists = true
			}
		}
		if !exists {
			//fmt.Printf("Adding Minion to Minion_List: %s\n", minion.Minion_Name)
			unique_minions = append(unique_minions, minion)
		}

	}
	m.Minion_List = append(m.Minion_List, unique_minions...)
}

func (m *Target_Minions) Add_Offline_Minions(list []string) {

	unique_minions := []string{}

	for _, minion := range list {
		exists := false
		// Check if the Minion_ID is already in the map
		for _, minion2 := range m.Offline_Minions {
			if minion == minion2 {
				// If it is, don't add it again
				exists = true
			}
		}
		if !exists {
			unique_minions = append(unique_minions, minion)
		}

	}
	m.Offline_Minions = append(m.Offline_Minions, unique_minions...)
}

func (m *Target_Minions) Add_No_Target_Minions(list []string) {
	//m.No_Targets_Minions = make([]string, len(list))
	unique_minions := []string{}

	for _, minion := range list {
		exists := false
		// Check if the Minion_ID is already in the map
		for _, minion2 := range m.No_Targets_Minions {
			if minion == minion2 {
				// If it is, don't add it again
				exists = true
			}
		}
		if !exists {
			unique_minions = append(unique_minions, minion)
		}

	}
	m.No_Targets_Minions = append(m.No_Targets_Minions, unique_minions...)
}
