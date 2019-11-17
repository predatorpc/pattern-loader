package models

type ConfigApi struct {
	Url string
	Method string
}

type ConfigSettings struct {
	AddHeader string
	TimeOutSec int
	SendPostBack bool
}

type ConfigParams struct {
	Click_Hash bool
	Flow_Hash bool
	Sub1 bool
	Sub2 bool
	Sub3 bool
	Sub4 bool
	Sub5 bool
	Sub6 bool
	Sub7 bool
	Sub8 bool
	Sub9 bool
	Sub10 bool
	Pl_ID bool
	Lp_ID bool
	Webmaster_currency_ID bool
	Webmaster_ID bool
	Geo bool
	External_Api_ID bool
}

type GeneralConfigFile struct {
	IntegrationName string
	Api ConfigApi
	Settings ConfigSettings
	Params ConfigParams
}
