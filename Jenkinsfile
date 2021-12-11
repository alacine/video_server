pipeline {
    agent {
        docker {
            image 'golang:1.16-alpine'
            args '-u root:root'
        }
    }
    stages {
        stage('build') {
            steps {
                sh 'apk add --no-cache make'
                sh 'make'
            }
        }
    }
}
