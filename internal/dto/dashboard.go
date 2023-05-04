package dto

type PanelGroupDataOutput struct {
	ServiceNum int64 `json:"serviceNum"`
	AppNum     int64 `json:"appNum"`
}

type DashServiceStatItemOutput struct {
	Name     string `json:"name"`
	LoadType int    `json:"load_type"`
	Value    int64  `json:"value"`
}

type DashServiceStatOutput struct {
	Legend []string                    `json:"legend"`
	Data   []DashServiceStatItemOutput `json:"data"`
}
