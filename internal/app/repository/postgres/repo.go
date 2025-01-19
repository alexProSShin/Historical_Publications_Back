package postgresrepo

import (
	"backend/internal/app/repository"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db          *gorm.DB
	minioClient *minio.Client
	bucketName  string
}

func NewPostgresRepository(db *gorm.DB, minioClient *minio.Client, bucketName string) repository.PostgresRepo {
	return &PostgresRepository{
		db:          db,
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}
