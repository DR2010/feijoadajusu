package commonstruct

// DatabaseX is a struct
type DatabaseX struct {
	Location   string // location of the database localhost, something.com, etc
	Database   string // database name
	Collection string // collection name
	APIServer  string // apiserver name
}

// Resultado is a struct
type Resultado struct {
	ErrorCode        string // error code
	ErrorDescription string // description
	IsSuccessful     string // Y or N
	ReturnedValue    string
}

// RestEnvVariables = restaurante environment variables
//
type RestEnvVariables struct {
	Organisation            string // Name of the organisation
	APIMongoDBLocation      string // location of the database localhost, something.com, etc
	APIMongoDBDatabase      string // database name
	APIAPIServerPort        string // collection name
	APIAPIServerIPAddress   string // apiserver name
	WEBDebug                string // debug
	RecordCurrencyTick      string // debug
	RunningFromServer       string // debug
	WEBServerPort           string // collection name
	ConfigFileFound         string // collection name
	ApplicationID           string // collection name
	UserID                  string // collection name
	AppFestaJuninaEnabled   string
	AppBelnorthEnabled      string
	AppBitcoinEnabled       string
	MSAPIdishesPort         string // Microservices Port Dishes
	MSAPIordersPort         string // Microservices Port Orders
	MSAPIdishesIPAddress    string // Microservices IP Address
	MSAPIordersIPAddress    string // Microservices IP Address
	MSAPImainPort           string // Microservices IP Address
	MSAPImainIPAddress      string // Microservices IP Address
	SecurityMicroservice    string // Microservices IP Address
	SecurityMicroserviceURL string // Microservices IP Address
	SYSID                   string //
}
