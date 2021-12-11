pipeline {
    agent { docker 'golang:1.16-alpine' }
    stages {
        stage('build') {
            steps {
                sh 'apk add make'
                sh 'make'
            }
        }
    }
}
