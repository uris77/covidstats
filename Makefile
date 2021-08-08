buildGCR: buildLinux
	gcloud builds submit --tag gcr.io/epi-belize/covidstats
deployCloudRun:
	gcloud run deploy --image gcr.io/epi-belize/covidstats --platform managed --service-account covidstats@epi-belize.iam.gserviceaccount.com --set-env-vars GCP_PROJECT_ID=epi-belize --memory 1024
buildLocal:
	go build -mod=readonly -o bin/server cmd/main.go
buildLinux:
	GOOS=linux GOARCH=amd64 go build -mod=readonly -o bin/server cmd/main.go
dockerBuild: buildLinux
	docker build -t gcr.io/epi-belize/covidstats .
dockerRun:
	docker run --env GCP_PROJECT_ID=epi-belize --rm -p 8080:8080 gcr.io/epi-belize/covidstats
