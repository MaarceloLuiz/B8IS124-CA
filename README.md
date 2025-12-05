# worldle-firebase

## Gcloud auth
```
gcloud projects create YOUR-PROJECT-ID --name="YOUR PROJECT NAME"
gcloud config set project YOUR-PROJECT-ID
```

### APIs
```
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable firebase.googleapis.com
```
## Firebase auth:
```
firebase login
firebase projects:addfirebase YOUR-PROJECT-ID
```
## Firebase Storage Deploy Rules
```
firebase deploy --only storage
```
- After the storage was created, the images were manually uploaded using the UI

## Get Google Application Credentials (Local Development)
1. Make sure you're in the right project
```
gcloud config set project YOUR_PROJECT_ID
```
2. Create a service account (if you haven't already)
```
gcloud iam service-accounts create firebase-admin --display-name="Firebase Admin SDK Service Account"
```

3. Grant the service account Storage Admin role
```
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID --member="serviceAccount:firebase-admin@YOUR_PROJECT_ID.iam.gserviceaccount.com" --role="roles/storage.admin"
```

4. Generate and download the key file
```
gcloud iam service-accounts keys create serviceAccountKey.json --iam-account=firebase-admin@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

## Deploy Backend to Cloud RUN
1. Check current project
```
gcloud config get-value project
```
2. If not set, list your projects
```
gcloud projects list
```
3. Set your project
```
gcloud config set project YOUR_PROJECT_ID
```
4. Build and Push Dockerimage to gcloud
```
docker build -f YOUR_DOCKERFILE -t gcr.io/YOUR_ARTFACTORY_REPO/YOUR_IMAGE_NAME .
docker push gcr.io/YOUR_ARTFACTORY_REPO/YOUR_IMAGE_NAME
```

5. Deploy to Cloud Run
```
gcloud run deploy YOUR_PROJECT_NAME \
  --image gcr.io/YOUR_PROJECT_ID/YOUR_PROJECT_NAME \
  --platform managed \
  --region europe-west2 \
  --allow-unauthenticated \
  --set-env-vars YOUR_ENV=YOUR_ENV,ALLOWED_ORIGIN=http://localhost:3000
```

6. Check the deployment status
```
gcloud run services describe YOUR_PROJECT_NAME --region europe-west2
```

## Deploy Frontend to Firebase Hosting
1. Login to Firebase (if not already)
```
firebase login
```
2. Initialize Firebase
```
firebase init hosting
```
3. Build the React App for Prod
```
npm install
npm run build
```
4. Deploy to Firebase Hosting
```
firebase deploy --only hosting
```
## Update Backend CORS
```
gcloud run services update worldle-backend --region europe-west2 --update-env-vars ALLOWED_ORIGIN=https://YOUR_FRONTEND_URL
```

## Automated Deployment Pipeline (NOT FINISHED)
- WORKING ON IT