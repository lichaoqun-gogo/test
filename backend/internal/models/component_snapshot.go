package models

// ComponentSnapshot 开源依赖包表
type ComponentSnapshot struct {
	ID                int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CurrentID         string     `gorm:"column:current_id" json:"current_id"`
	ParentComponentID *int64     `gorm:"column:PARENT_COMPONENT_ID" json:"PARENT_COMPONENT_ID"`
	Author            string     `gorm:"column:AUTHOR" json:"AUTHOR"`
	Publisher         string     `gorm:"column:PUBLISHER" json:"PUBLISHER"`
	Group             string     `gorm:"column:_group" json:"_group"`
	Name              string     `gorm:"column:_name" json:"_name"`
	Version           string     `gorm:"column:_version" json:"_version"`
	Description       string     `gorm:"column:DESCRIPTION" json:"DESCRIPTION"`
	Classifier        string     `gorm:"column:CLASSIFIER" json:"CLASSIFIER"`
	Filename          string     `gorm:"column:FILENAME" json:"FILENAME"`
	Extension         string     `gorm:"column:EXTENSION" json:"EXTENSION"`
	Internal          *int       `gorm:"column:INTERNAL" json:"INTERNAL"`
	MD5               string     `gorm:"column:MD5" json:"MD5"`
	SHA1              string     `gorm:"column:SHA1" json:"SHA1"`
	SHA256            string     `gorm:"column:SHA_256" json:"SHA_256"`
	SHA384            string     `gorm:"column:SHA_384" json:"SHA_384"`
	SHA512            string     `gorm:"column:SHA_512" json:"SHA_512"`
	SHA3_256          string     `gorm:"column:SHA3_256" json:"SHA3_256"`
	SHA3_384          string     `gorm:"column:SHA3_384" json:"SHA3_384"`
	SHA3_512          string     `gorm:"column:SHA3_512" json:"SHA3_512"`
	BLAKE2B_256       string     `gorm:"column:BLAKE2B_256" json:"BLAKE2B_256"`
	BLAKE2B_384       string     `gorm:"column:BLAKE2B_384" json:"BLAKE2B_384"`
	BLAKE2B_512       string     `gorm:"column:BLAKE2B_512" json:"BLAKE2B_512"`
	BLAKE3            string     `gorm:"column:BLAKE3" json:"BLAKE3"`
	Notes             string     `gorm:"column:NOTES;type:mediumtext" json:"NOTES"`
	PURL              string     `gorm:"column:PURL" json:"PURL"`
	PURLCoordinates   string     `gorm:"column:PURLCOORDINATES" json:"PURLCOORDINATES"`
	CPE               string     `gorm:"column:CPE" json:"CPE"`
	SWIDTagID         string     `gorm:"column:SWIDTAGID" json:"SWIDTAGID"`
	DirectDependencies string    `gorm:"column:DIRECT_DEPENDENCIES;type:mediumtext" json:"DIRECT_DEPENDENCIES"`
	ExternalReferences []byte    `gorm:"column:EXTERNAL_REFERENCES;type:mediumblob" json:"EXTERNAL_REFERENCES"`
	Copyright         string     `gorm:"column:COPYRIGHT" json:"COPYRIGHT"`
	License           string     `gorm:"column:LICENSE" json:"LICENSE"`
	LicenseID         *int64     `gorm:"column:LICENSE_ID" json:"LICENSE_ID"`
	LastRiskscore     float64    `gorm:"column:LAST_RISKSCORE" json:"LAST_RISKSCORE"`
	ProjectID         *int64     `gorm:"column:PROJECT_ID" json:"PROJECT_ID"`
}

func (ComponentSnapshot) TableName() string {
	return "component_snapshot"
}
