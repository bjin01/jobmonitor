package email

type Templates_Dir struct {
	Dir string
}

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

type SPMigration_Email_Body struct {
	Host                      string
	Port                      int
	T7user                    string
	Template_dir              string
	SPmigration_Tracking_File string
	Recipients                []string
}
