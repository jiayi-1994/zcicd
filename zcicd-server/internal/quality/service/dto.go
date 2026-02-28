package service

type CreateTestConfigReq struct {
	Name      string `json:"name" binding:"required"`
	TestType  string `json:"test_type"`
	Framework string `json:"framework"`
	Command   string `json:"command"`
	Timeout   int    `json:"timeout"`
}

type UpdateTestConfigReq struct {
	Name      string `json:"name"`
	TestType  string `json:"test_type"`
	Framework string `json:"framework"`
	Command   string `json:"command"`
	Timeout   int    `json:"timeout"`
	Enabled   *bool  `json:"enabled"`
}

type CreateScanConfigReq struct {
	Name            string `json:"name" binding:"required"`
	ScanType        string `json:"scan_type"`
	SonarProjectKey string `json:"sonar_project_key"`
}

type UpdateScanConfigReq struct {
	Name            string `json:"name"`
	ScanType        string `json:"scan_type"`
	SonarProjectKey string `json:"sonar_project_key"`
	Enabled         *bool  `json:"enabled"`
}

type QualityGateReq struct {
	MinCoverage        *float64 `json:"min_coverage"`
	MaxBugs            *int     `json:"max_bugs"`
	MaxVulnerabilities *int     `json:"max_vulnerabilities"`
	MaxCodeSmells      *int     `json:"max_code_smells"`
	MaxDuplications    *float64 `json:"max_duplications"`
	BlockDeploy        *bool    `json:"block_deploy"`
}
