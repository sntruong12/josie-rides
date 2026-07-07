# Josie Rides

This repository follows the book Let's Go by Alex Edwards. It houses golang code to be deployed on AWS. I'm using this as practice for modern cloud development. 

## Golang Notes

| CLI Command | Description |
| ----------- | ----------- |
| go mod init nameOfPackage | creates your golang project as a module, creates a document with all the project dependencies, your module's path, and module's go version |
| go run pathToPackage | a shortcut that compiles your code, creates an executable binary in your /tmp directory, and then runs this binary in one step |

## Shell Commands

| Command | Description |
| ------- | ----------- |
| curl -i -X POST url | sends and http request with the POST method and include the header in the output |
| curl | in general wrap the url parameter in single quotes so that you won't have to escape some special characters |

## Project Structure

| directory/file | description |
| -------------- | ----------- |
| cmd | The cmd directory will contain the application-specific code for the executable applications in the project. For now we’ll have just one executable application — the web application — which will live under the cmd/web directory. |
| internal | The internal directory will contain the ancillary non-application-specific code used in the project. We’ll use it to hold potentially reusable code like validation helpers and the SQL database models for the project. |
| ui | The ui directory will contain the user-interface assets used by the web application. Specifically, the ui/html directory will contain HTML templates, and the ui/static directory will contain static files (like CSS and images).|

Two benefits for using this structure
1. Clean separation of Go and non-Go assets. Go code lives in cmd and internal directories, leaving the root directory free to hold non-Go assets like UI files, makefiles, module definitions like our go.mod file. Easier to manage when it comes to building and deploying the app.
2. Scales like if you want to add more executable application to the project. Adding a CLI to automate some admin tasks, would be able to reuse internal code keep your code DRY.

Some tradeoffs here
1. Blast Radius, if multiple services share the internal code, a bug in internal could affect multiple services like the web server and a cloud function all at once.
2. Slower CI/CD pipelines, new updates to code would have to be handled specifically. CI flows would have to figure what changed and which services need to be built. Tests and builds can take longer, thus impacting speed of delivery.

