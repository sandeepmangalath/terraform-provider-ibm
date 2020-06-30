module github.com/IBM-Cloud/terraform-provider-ibm

go 1.12

require (
	github.com/Bowery/prompt v0.0.0-20190916142128-fa8279994f75 // indirect
	github.com/IBM-Cloud/bluemix-go v0.0.0-20200618045128-dcf70b676caf
	github.com/IBM-Cloud/power-go-client v1.0.24
	github.com/IBM/apigateway-go-sdk v0.0.0-20200414212859-416e5948678a
	github.com/IBM/cloudpakfordata-go-sdk v0.0.0-20191003093330-fead1522985f // indirect
	github.com/IBM/dns-svcs-go-sdk v0.0.3
<<<<<<< HEAD
	github.com/IBM/go-sdk-core v1.1.0
	github.com/IBM/go-sdk-core/v3 v3.3.1
	github.com/IBM/ibm-cos-sdk-go v1.2.0
=======
	github.com/IBM/go-sdk-core/v3 v3.2.4
	github.com/IBM/ibm-cos-sdk-go v1.3.1
	github.com/IBM/ibm-cos-sdk-go-config v1.0.0
>>>>>>> 4a6664ff625d188991437b2326c2e11c545f6938
	github.com/IBM/keyprotect-go-client v0.3.5-0.20200325142150-b63163832e26
	github.com/apache/openwhisk-client-go v0.0.0-20200201143223-a804fb82d105
	github.com/apparentlymart/go-cidr v1.0.1
	github.com/cloudfoundry/jibber_jabber v0.0.0-20151120183258-bcc4c8345a21 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/structs v1.1.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/runtime v0.19.15 // indirect
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/go-openapi/validate v0.19.8 // indirect
	github.com/go-test/deep v1.0.4 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/hashicorp/go-getter v1.4.2-0.20200106182914-9813cbd4eb02 // indirect
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/hcl/v2 v2.3.0 // indirect
	github.com/hashicorp/terraform-config-inspect v0.0.0-20191212124732-c6ae6269b9d7 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.6.0
	github.com/hokaccha/go-prettyjson v0.0.0-20170213120834-e6b9231a2b1c // indirect
	github.com/minsikl/netscaler-nitro-go v0.0.0-20170827154432-5b14ce3643e3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/renier/xmlrpc v0.0.0-20170708154548-ce4a1a486c03 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/softlayer/softlayer-go v0.0.0-20190814165317-b9062a914a22
	github.com/zclconf/go-cty v1.2.1 // indirect
	github.ibm.com/Bluemix/riaas-go-client v0.0.0-20191018070922-afd27ac04d4f
	github.ibm.com/ibmcloud/namespace-go-sdk v0.0.0-00010101000000-000000000000
	github.ibm.com/ibmcloud/vpc-go-sdk v0.0.0-00010101000000-000000000000
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	github.ibm.com/ibmcloud/networking-go-sdk v0.0.0-00010101000000-000000000000
	github.ibm.com/ibmcloud/vpc-go-sdk v0.0.0-00010101000000-000000000000
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
)

replace github.ibm.com/Bluemix/riaas-go-client v0.0.0-20191018070922-afd27ac04d4f => ./common/github.ibm.com/Bluemix/riaas-go-client

replace github.ibm.com/ibmcloud/vpc-go-sdk => ./common/github.ibm.com/ibmcloud/vpc-go-sdk

replace github.ibm.com/ibmcloud/namespace-go-sdk => ./common/github.ibm.com/ibmcloud/namespace-go-sdk

replace github.com/softlayer/softlayer-go v0.0.0-20190814165317-b9062a914a22 => ./common/github.com/softlayer/softlayer-go

replace github.ibm.com/ibmcloud/networking-go-sdk => ./common/github.ibm.com/ibmcloud/networking-go-sdk
