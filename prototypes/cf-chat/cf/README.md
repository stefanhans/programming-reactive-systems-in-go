### Cloud Functions As List Of Members Service

Using the Google Cloud Functions ([Go 1.11.5 runtime](https://cloud.google.com/functions/docs/concepts/go-runtime)), 
I have implemented a simple list of members service with [Google Firestore](https://cloud.google.com/firestore/docs/).

After preparing the environment, we start setting some environment variables:

* GCP_PROJECT=\<project id>
* GOOGLE_APPLICATION_CREDENTIALS=\<json file>
* GCP_REGION=\<region>

We use Go modules for our dependencies, therefore, 
having the source files in place, we execute the following:

```
export GO111MODULE=on
go mod init && go mod vendor
```

Now, we can deploy our functions:

```
gcloud alpha functions deploy subscribe --region $GCP_REGION --entry-point Subscribe --runtime go111 --trigger-http
gcloud alpha functions deploy list --region $GCP_REGION --entry-point List --runtime go111 --trigger-http
gcloud alpha functions deploy unsubscribe --region $GCP_REGION --entry-point Unsubscribe --runtime go111 --trigger-http
gcloud alpha functions deploy reset --region $GCP_REGION --entry-point Reset --runtime go111 --trigger-http
```

Using the URL from the deployment, we test the functions:


```
curl <url>/subscribe <name> <ip> <port>
curl <url>/list
curl <url>/unsubscribe <id>
curl <url>/reset
```
