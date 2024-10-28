pipeline {
  agent any

  tools {
    go '1.23.2'
  }

  stages {
    stage('Test go') {
      steps {
        sh 'go test -v ./...'
      }
    }
  }
}

