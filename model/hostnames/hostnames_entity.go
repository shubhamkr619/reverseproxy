package hostnames

import (
	"fmt"
	"sync"

	"github.com/revproxy/src/gorm-library/database"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	mutex       = &sync.Mutex{}
	HostNameMap *HostNames
)

type HostNames struct {
	names map[string]string
}

func getHostnameInstance() *HostNames {
	if HostNameMap == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if HostNameMap == nil {
			HostNameMap = &HostNames{
				names: make(map[string]string),
			}
		} else {
			fmt.Printf("Struct already exist...+%v\n\n", HostNameMap)
		}
	} else {
		fmt.Printf("Struct already exist...+%v\n\n", HostNameMap)
	}
	return HostNameMap
}

type HostName struct {
	gorm.Model
	Domain string `json:"domain" gorm:"index:,unique"`
	Bucket string `json:"bucket"`
}

func init() {
	db := database.UseDB()
	db.AutoMigrate(&HostName{})
}

const (
	tableName = "hostnames"
	pkName    = "domain"
)

func (h HostName) TableName() string {
	return tableName
}

func (h HostName) String() string {
	return fmt.Sprintf("{domain: %s, bucket: %s}", h.Domain, h.Bucket)
}

func CreateHostName(db *gorm.DB, hostname *HostName) (err error) {
	err = db.Create(hostname).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAllHostNames(db *gorm.DB) ([]HostName, error) {
	fmt.Println("called GetAllHostNames")
	var hostnames []HostName
	err := db.Find(&hostnames).Error
	if err != nil {
		return []HostName{}, err
	}
	getHostnameInstance()
	for _, row := range hostnames {
		HostNameMap.names[row.Domain] = row.Bucket
	}
	return hostnames, nil
}

func GetBucketByHostName(db *gorm.DB, hostName string) string {
	getHostnameInstance()
	if len(HostNameMap.names) == 0 {
		GetAllHostNames(db)
	}
	bucket, ok := findInMap(hostName, HostNameMap.names)
	if !ok {
		fmt.Printf("No value found for the hostname in map %v \n", HostNameMap)
		// Need to execute a refresh for given hostname. Could cause DOS attach in some cases,
		//however the front lb should drop this traffic.
		GetAllHostNames(db)
		bucket, ok = findInMap(hostName, HostNameMap.names)
		if !ok {
			fmt.Printf("retried: No value found for the hostname in map %v \n", HostNameMap)
			// Need to execute a refresh for given hostname. Could cause DOS attach in some cases,
			//however the front lb should drop this traffic.
			//panic(fmt.Errorf("error occurred as the host is not available on the proxy %v \n", hostName))
			fmt.Errorf("error occurred as the host is not available on the proxy %v \n", hostName)
		}

	}
	return bucket
}

func findInMap(r string, m map[string]string) (string, bool) {
	val, ok := m[r]
	return val, ok
}
