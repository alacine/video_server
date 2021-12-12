pipeline {
    agent {
        docker {
            image 'golang:1.16-alpine'
            args '-u root:root'
        }
    }
    triggers {
        GenericTrigger(
            genericVariables: [
              [key: 'appn', value: '$.app.name']
            ],
            token: 'abc123',
            causeString: 'Triggered on $ref',
            printContributedVariables: true,
            printPostContent: true
        )
    }
    stages {
        stage('env') {
            steps {
                sh 'echo $appn'
                sh 'printenv'
            }
        }
        stage('build') {
            steps {
                sh 'apk add --no-cache make'
                sh 'make'
            }
        }
    }
}
