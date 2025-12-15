# worldle-firebase

## Gcloud auth
```bash
gcloud projects create YOUR-PROJECT-ID --name="YOUR PROJECT NAME"
gcloud config set project YOUR-PROJECT-ID
```

### APIs
```bash
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable firebase.googleapis.com
```
## Firebase auth:
```bash
firebase login
firebase projects:addfirebase YOUR-PROJECT-ID
```
## Firebase Storage Deploy Rules
```bash
firebase deploy --only storage
```
- After the storage was created, the images were manually uploaded using the UI

## Get Google Application Credentials (Local Development)
1. Make sure you're in the right project
```bash
gcloud config set project YOUR_PROJECT_ID
```
2. Create a service account (if you haven't already)
```bash
gcloud iam service-accounts create firebase-admin --display-name="Firebase Admin SDK Service Account"
```

3. Grant the service account Storage Admin role
```bash
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID --member="serviceAccount:firebase-admin@YOUR_PROJECT_ID.iam.gserviceaccount.com" --role="roles/storage.admin"
```

4. Generate and download the key file
```bash
gcloud iam service-accounts keys create serviceAccountKey.json --iam-account=firebase-admin@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

## Deploy Backend to Cloud RUN
1. Check current project
```bash
gcloud config get-value project
```
2. If not set, list your projects
```bash
gcloud projects list
```
3. Set your project
```bash
gcloud config set project YOUR_PROJECT_ID
```
4. Build and Push Dockerimage to gcloud
```bash
docker build -f YOUR_DOCKERFILE -t gcr.io/YOUR_ARTFACTORY_REPO/YOUR_IMAGE_NAME .
docker push gcr.io/YOUR_ARTFACTORY_REPO/YOUR_IMAGE_NAME
```

5. Deploy to Cloud Run
```bash
gcloud run deploy YOUR_PROJECT_NAME \
  --image gcr.io/YOUR_PROJECT_ID/YOUR_PROJECT_NAME \
  --platform managed \
  --region europe-west2 \
  --allow-unauthenticated \
  --set-env-vars YOUR_ENV=YOUR_ENV,ALLOWED_ORIGIN=http://localhost:3000
```

6. Check the deployment status
```bash
gcloud run services describe YOUR_PROJECT_NAME --region europe-west2
```

## Deploy Frontend to Firebase Hosting
1. Login to Firebase (if not already)
```bash
firebase login
```
2. Initialize Firebase
```bash
firebase init hosting
```
3. Build the React App for Prod
```bash
npm install
npm run build
```
4. Deploy to Firebase Hosting
```bash
firebase deploy --only hosting
```
## Update Backend CORS
```bash
gcloud run services update worldle-backend --region europe-west2 --update-env-vars ALLOWED_ORIGIN=https://YOUR_FRONTEND_URL
```

## Automated Deployment Pipeline

This project uses **Google Cloud Build** for continuous deployment. When code is pushed to the `main` branch, it automatically:
1. Builds and deploys the **Backend** to Cloud Run
2. Builds and deploys the **Frontend** to Firebase Hosting

### Pipeline Configuration

The deployment pipeline is defined in `cloudbuild.yaml` and consists of 6 steps:
- **Backend:** Build Docker image → Push to GCR → Deploy to Cloud Run
- **Frontend:** Install dependencies → Build production bundle → Deploy to Firebase Hosting

### Setting Up the Trigger (via Google Cloud Console)

1. **Navigate to Cloud Build Triggers:**
   - Go to: [Cloud Build Triggers](https://console.cloud.google.com/cloud-build/triggers)
   - Select your project

2. **Create a New Trigger:**
   - Click **"Create Trigger"**
   - **Name:** `worldle-deploy`
   - **Event:** Push to a branch
   - **Source:** Connect your GitHub repository (`MaarceloLuiz/B8IS124-CA`)
   - **Branch:** `^main$`
   - **Configuration:** Cloud Build configuration file (`cloudbuild.yaml`)

3. **Add Substitution Variables:**
   Under **Advanced > Substitution variables**, add:
   - `_MAPS_API_KEY`: Your Google Maps API Key
   - `_FIREBASE_TOKEN`: Generate with `firebase login:ci`
   - `_ALLOWED_ORIGIN`: Your Firebase Hosting URL
   - `_BACKEND_URL`: Your Cloud Run backend URL

4. **Create and Test:**
   - Click **"Create"**
   - Push to `main` branch to trigger automatic deployment

### Required Permissions

The Cloud Build service account needs these IAM roles:
```bash
# Grant Cloud Run Admin
gcloud projects add-iam-policy-binding PROJECT_ID \
  --member="serviceAccount:PROJECT_NUMBER@cloudbuild.gserviceaccount.com" \
  --role="roles/run.admin"

# Grant Service Account User
gcloud projects add-iam-policy-binding PROJECT_ID \
  --member="serviceAccount:PROJECT_NUMBER@cloudbuild.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"
```

5. Check triggers list
```bash
gcloud builds triggers list
```