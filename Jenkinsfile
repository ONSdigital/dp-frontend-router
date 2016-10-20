#!groovy

node {
    def gopath = pwd()

    ws("${gopath}/src/github.com/ONSdigital/dp-frontend-router") {
        stage('Checkout') {
            checkout scm
            sh 'git clean -dfx'
        }

        stage('Build') {
            sh "GOPATH=${gopath} go build"
        }

        stage('Test') {
            sh "GOPATH=${gopath} go test ./..."
        }
    }
}
