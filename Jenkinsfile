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
  stage('SonarQube Analysis') {
    steps {
      script {
        scannerHome = tool 'SonarScanner';
      }
      withSonarQubeEnv('vezgammon') {
        sh "${scannerHome}/bin/sonar-scanner"
      }
    }
  }
}

