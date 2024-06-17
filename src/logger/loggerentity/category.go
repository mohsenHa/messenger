package loggerentity

type Category string
type SubCategory string
type ExtraKey string

const (
	CategoryNotDefined      Category = "NotDefined"
	CategoryRequestResponse Category = "RequestResponse"
	CategoryGeneral         Category = "General"
	CategoryInternal        Category = "Internal"
	CategoryMysql           Category = "Mysql"
	CategoryRabbitMQ        Category = "RabbitMQ"
	CategoryValidation      Category = "Validation"
)

const (
	SubCategoryNotDefined     SubCategory = "Not defined"
	SubCategoryGeneralStartup SubCategory = "Startup"

	SubCategoryRabbitMQConnection SubCategory = "Connection"
	SubCategoryRabbitMQChannel    SubCategory = "Channel"
	SubCategoryRabbitMQAck        SubCategory = "Ack"

	SubCategoryMysqlMigration SubCategory = "Migration"
	SubCategoryMysqlSelect    SubCategory = "Select"
	SubCategoryMysqlRollback  SubCategory = "Rollback"
	SubCategoryMysqlCommit    SubCategory = "Commit"
	SubCategoryMysqlUpdate    SubCategory = "Update"
	SubCategoryMysqlDelete    SubCategory = "Delete"
	SubCategoryMysqlInsert    SubCategory = "Insert"

	SubCategoryInternalRequest  SubCategory = "Request"
	SubCategoryInternalResponse SubCategory = "Response"

	SubCategoryValidationMobileValidation   SubCategory = "MobileValidation"
	SubCategoryValidationPasswordValidation SubCategory = "PasswordValidation"
)

const (
	ExtraKeyAppName       ExtraKey = "AppName"
	ExtraKeyLoggerName    ExtraKey = "Logger"
	ExtraKeyClientIP      ExtraKey = "ClientIP"
	ExtraKeyRemoteIP      ExtraKey = "RemoteIP"
	ExtraKeyMethod        ExtraKey = "Method"
	ExtraKeyStatusCode    ExtraKey = "StatusCode"
	ExtraKeyResponseSize  ExtraKey = "ResponseSize"
	ExtraKeyURI           ExtraKey = "KeyURI"
	ExtraKeyURIPath       ExtraKey = "KeyURI"
	ExtraKeyLatency       ExtraKey = "Latency"
	ExtraKeyRequestBody   ExtraKey = "RequestBody"
	ExtraKeyRequestID     ExtraKey = "RequestId"
	ExtraKeyHost          ExtraKey = "Host"
	ExtraKeyContentLength ExtraKey = "ContentLength"
	ExtraKeyProtocol      ExtraKey = "Protocol"
	ExtraKeyResponseBody  ExtraKey = "ResponseBody"
	ExtraKeyErrorMessage  ExtraKey = "ErrorMessage"
	ExtraKeyParentTrace   ExtraKey = "ParentTrace"
)
