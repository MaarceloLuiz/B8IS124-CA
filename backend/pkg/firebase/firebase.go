package firebase

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FirebaseConfig struct {
	ProjectID string
	Bucket    string
}

var firebaseApp *firebase.App

func GetImageURL(fileName string) (string, error) {
	config, err := loadFirebaseConfig()
	if err != nil {
		logrus.Errorf("Failed to load Firebase config: %v", err)
		return "", err
	}

	url := fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/silhouettes%%2F%s.png?alt=media",
		config.Bucket,
		fileName,
	)

	// logrus.Infof("Generated image URL: %s", url) // DEBUG LOG
	return url, nil
}

func GetAllTerritoriesFromStorage() ([]string, error) {
	if err := initFirebaseApp(); err != nil {
		logrus.Errorf("Failed to initialize Firebase app: %v", err)
		return nil, fmt.Errorf("failed to initialize Firebase app: %v", err)
	}

	ctx := context.Background()

	logrus.Info("Getting Firebase storage client...")
	firebaseClient, err := firebaseApp.Storage(ctx)
	if err != nil {
		logrus.Errorf("Error getting storage client: %v", err)
		return nil, fmt.Errorf("error getting storage client: %v", err)
	}

	logrus.Info("Getting default bucket...")
	bucket, err := firebaseClient.DefaultBucket()
	if err != nil {
		logrus.Errorf("Error getting default bucket: %v", err)
		return nil, fmt.Errorf("error getting default bucket: %v", err)
	}

	var countries []string

	logrus.Info("Listing objects with prefix 'silhouettes/'...")
	it := bucket.Objects(ctx, &storage.Query{
		Prefix: "silhouettes/",
	})

	fileCount := 0
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			logrus.Infof("Finished iterating. Found %d files total", fileCount)
			break
		}
		if err != nil {
			logrus.Errorf("Error iterating objects: %v", err)
			return nil, fmt.Errorf("error iterating objects: %v", err)
		}

		fileCount++
		logrus.Infof("Found file #%d: %s", fileCount, attrs.Name)

		if strings.HasSuffix(attrs.Name, ".png") {
			fileName := strings.TrimPrefix(attrs.Name, "silhouettes/")
			countryName := strings.TrimSuffix(fileName, ".png")
			countries = append(countries, countryName)
			logrus.Infof("Added country: %s", countryName)
		}
	}

	if len(countries) == 0 {
		logrus.Error("No countries found in Firebase Storage")
		return nil, fmt.Errorf("no countries found in Firebase Storage")
	}

	logrus.Infof("Successfully loaded %d countries", len(countries))
	return countries, nil
}

func initFirebaseApp() error {
	config, err := loadFirebaseConfig()
	if err != nil {
		return fmt.Errorf("failed to load firebase config: %v", err)
	}

	if firebaseApp != nil {
		logrus.Info("Firebase app already initialized")
		return nil
	}

	firebaseConfig := &firebase.Config{
		StorageBucket: config.Bucket,
	}

	ctx := context.Background()
	serviceAccountPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	// Check if service account file exists
	if serviceAccountPath != "" {
		if _, err := os.Stat(serviceAccountPath); os.IsNotExist(err) {
			logrus.Errorf("Service account file does not exist at path: %s", serviceAccountPath)
			return fmt.Errorf("service account file not found: %s", serviceAccountPath)
		}
		logrus.Info("Service account file exists")
	}

	var app *firebase.App

	if serviceAccountPath != "" {
		logrus.Info("Using service account credentials")
		app, err = firebase.NewApp(ctx, firebaseConfig, option.WithCredentialsFile(serviceAccountPath))
	} else {
		logrus.Info("Using default credentials (Cloud Run)")
		app, err = firebase.NewApp(ctx, firebaseConfig)
	}

	if err != nil {
		logrus.Errorf("Error initializing Firebase app: %v", err)
		return fmt.Errorf("error initializing firebase app: %v", err)
	}

	firebaseApp = app
	logrus.Info("Firebase app initialized successfully")
	return nil
}

func loadFirebaseConfig() (*FirebaseConfig, error) {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		logrus.Error("FIREBASE_PROJECT_ID environment variable is missing")
		return nil, fmt.Errorf("missing FIREBASE_PROJECT_ID environment variable")
	}

	bucket := fmt.Sprintf("%s.firebasestorage.app", projectID)
	// logrus.Infof("Loaded Firebase config - Project: %s, Bucket: %s", projectID, bucket) // DEBUG LOG
	logrus.Info("Firebase config loaded successfully")

	return &FirebaseConfig{
		ProjectID: projectID,
		Bucket:    bucket,
	}, nil
}
