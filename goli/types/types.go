package types

type GoliRequestBody struct {
	Name      string `json:"name"`
	Image     string `json:"image"`
	Network   string `json:"network"`
	Port_Ex   string `json:"port_ex"`
	Port_In   string `json:"port_in"`
	V_Map     bool   `json:"v_map"`
	Volume_Ex string `json:"volume_ex"`
	Volume_In string `json:"volume_in"`
	Opts      string `json:"opts"`
}
