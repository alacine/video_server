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
                [
                     key: 'appname',
                     value: '$.app.name',
                     expressionType: 'JSONPath', 
                     defaultValue: 'default app name'
                ]
            ],
            token: 'abc123',
            causeString: 'Triggered on $appname',
            printContributedVariables: true,
            printPostContent: true
        )
    }
    stages {
        stage('env') {
            steps {
                sh 'echo $appname'
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
